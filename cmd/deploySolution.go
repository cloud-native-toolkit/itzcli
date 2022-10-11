/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"bytes"
	b64 "encoding/base64"
	"encoding/json"
	"fmt"
	logger "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.ibm.com/skol/atkcli/internal/prompt"
	"github.ibm.com/skol/atkcli/pkg"
	"github.ibm.com/skol/atkmod"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"
)

var fn string
var sol string
var cluster string
var rez string
var useCached bool

// deploySolutionCmd represents the deployProject command
var deploySolutionCmd = &cobra.Command{
	Use:   "deploy",
	Short: "Deploys the specified solution.",
	Long: `Use this command to deploy the specified solution
locally in your own environment. You can specify the environment by using
either --cluster-name or --reservation as a target.

    --cluster-name requires the name of a cluster that has been deployed
using ocpnow. To see the clusters that are configured, use the "atk configure 
list" command to list the available clusters. If you have none, you may need to
import the ocpnow configuration using the "atk configure import" command. See
the help for those commands for more information.

    --reservation requires the id of a reservation in the TechZone system. Use
the "atk reservation list" command to list the available reservations.`,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		SetLoggingLevel(cmd, args)
		if len(fn) == 0 && len(sol) == 0 {
			return fmt.Errorf("either \"--solution\" or \"--file\" must be specified.")
		}
		if len(cluster) == 0 && len(rez) == 0 {
			return fmt.Errorf("either \"--cluster-name\" or \"--reservation\" must be specified.")
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		logger.Infof("Deploying solution \"%s\"...", sol)
		return DeploySolution(cmd, args)
	},
}

var jenkinsBuildProjectCmd *pkg.ServiceClient
var jenkinsGetJobParamsCmd *pkg.ServiceClient

func init() {
	solutionCmd.AddCommand(deploySolutionCmd)
	deploySolutionCmd.Flags().StringVarP(&fn, "file", "f", "", "The full path to the solution file to be deployed.")
	deploySolutionCmd.Flags().StringVarP(&sol, "solution", "s", "", "The name of the solution to be deployed.")
	deploySolutionCmd.Flags().StringVarP(&cluster, "cluster-name", "c", "", "The name of the cluster created by ocpnow to target.")
	deploySolutionCmd.Flags().StringVarP(&rez, "reservation", "r", "", "The id of the reservation to target.")
	// TODO: Change this from true to false by default
	deploySolutionCmd.Flags().BoolVarP(&useCached, "use-cache", "u", true, "If true, uses a cached solution file instead of downloading from target.")
	deploySolutionCmd.MarkFlagsMutuallyExclusive("file", "solution")
	deploySolutionCmd.MarkFlagsMutuallyExclusive("reservation", "cluster-name")

	jenkinsGetJobParamsCmd = &pkg.ServiceClient{
		Method: "GET",
	}
}

// DeploySolution deploys the solution by handing it off to the bifrost
// API
func DeploySolution(cmd *cobra.Command, args []string) error {
	homedir, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	// Load up the reader based on the URI provided for the solution
	bifrostURL, err := url.Parse(viper.GetString("bifrost.api.url"))
	if err != nil {
		return fmt.Errorf("error trying to parse \"bifrost.api.url\", looks like a bad URL (value was: %s): %v", err, viper.GetString("bifrost.api.url"))
	}
	builderURL, err := url.Parse(viper.GetString("ci.api.url"))
	if err != nil {
		return fmt.Errorf("error trying to parse \"ci.api.url\", looks like a bad URL (value was: %s): %v", err, viper.GetString("ci.api.url"))
	}

	services := []pkg.Service{
		{
			DisplayName: "builder",
			ImgName:     viper.GetString("ci.api.image"),
			IsLocal:     viper.GetBool("ci.api.local"),
			URL:         builderURL,
			PreStart:    pkg.StatusHandler,
			Start:       pkg.StartHandler,
			PostStart:   initTokenAndSave,
			Volumes: map[string]string{
				viper.GetString("ci.localdir"): "/var/jenkins_home",
			},
			Envvars: map[string]string{
				"JENKINS_ADMIN_ID":       viper.GetString("ci.api.user"),
				"JENKINS_ADMIN_PASSWORD": viper.GetString("ci.api.password"),
			},
			Flags: []string{"--rm", "-d", "--privileged"},
		},
		{
			DisplayName: "integration",
			ImgName:     viper.GetString("bifrost.api.image"),
			IsLocal:     viper.GetBool("bifrost.api.local"),
			URL:         bifrostURL,
			PreStart:    pkg.StatusHandler,
			Start:       withEnvUpdates,
			Flags:       []string{"--rm", "-d"},
			Envvars: map[string]string{
				"JENKINS_API_USER": viper.GetString("ci.api.user"),
				"JENKINS_API_URL":  viper.GetString("ci.api.url"),
			},
		},
	}

	out := new(bytes.Buffer)
	ctx := &atkmod.RunContext{
		Out: out,
		Log: *logger.StandardLogger(),
	}

	err = pkg.StartupServices(ctx, services, pkg.Sequential)

	if err != nil {
		return err
	}

	// TODO: Now the services are started, we can use them like we would...
	// By starting with getting the ZIP file (and saving it in /tmp)
	var archiveFile string
	if len(sol) > 0 {
		if !useCached {
			uri := fmt.Sprintf("%s/solutions/%s/automation", viper.GetString("builder.api.url"), sol)
			logger.Debugf("Downloading solution file from URL <%s>...", uri)
			data, err := pkg.ReadHttpGetT(uri, viper.GetString("builder.api.token"))
			if err != nil {
				return err
			}
			dir, err := os.MkdirTemp(os.TempDir(), "atk-")
			if err != nil {
				return err
			}
			logger.Debugf("Writing solution file to directory <%s>", dir)
			archiveFile = filepath.Join(dir, fmt.Sprintf("%s.zip", sol))
			err = pkg.WriteFile(archiveFile, data)
			logger.Trace("Finished writing solution file")
		} else {
			logger.Infof("Using cached solution file for solution %s...", sol)
			archiveFile = filepath.Join(homedir, ".atk", "cache", fmt.Sprintf("%s.zip", sol))
		}

		// Now, post the ZIP file to the bifrost endpoint...
		err = pkg.PostFileToURL(archiveFile, fmt.Sprintf("%s/api/upload/builderPackage/%s", viper.GetString("bifrost.api.url"), sol))
		if err != nil {
			return err
		}

		err = pkg.Unzip(archiveFile, filepath.Join(viper.GetString("ci.localdir"), "workspace"))
		if err != nil {
			return err
		}

		logger.Infof("Finished creating pipeline for solution %s; starting deployment now...", sol)
		vars := make([]pkg.JobParam, 0)
		jenkinsGetJobParamsCmd.BaseURL = fmt.Sprintf("%s/api/jobs/%s/parameters", viper.GetString("bifrost.api.url"), sol)
		jenkinsGetJobParamsCmd.ResponseHandler = func(reader io.ReadCloser) error {
			defer reader.Close()
			data, err := io.ReadAll(reader)
			if err != nil {
				return err
			}
			err = json.Unmarshal(data, &vars)
			if err != nil {
				return err
			}
			return nil
		}

		err = pkg.Exec(jenkinsGetJobParamsCmd)
		logger.Debugf("Got vars: %v", vars)
		if err != nil {
			return err
		}

		ocpCfg := viper.GetStringSlice("ocpnow.configFiles")
		project, err := pkg.LoadProject(ocpCfg[0])
		if err != nil {
			return err
		}
		cInfo := project.Clusters[cluster]
		clusterVars, _ := pkg.ResolveVars(&cInfo, nil)
		logger.Debugf("Got cluster vars: %v", clusterVars)
		logger.Debugf("Using region: %s", clusterVars["TF_VAR_region"])
		credInfo := project.Credentials[cInfo.CredId]
		credVars, _ := pkg.ResolveVars(&credInfo, nil)
		logger.Debugf("Got cred vars: %v", credVars)

		// Now we have a list of the required parameters (vars), and we need
		// to look at the ones that we have and that have values (clusterVars and
		// credVars), and also look into the os.Environment. We'll build a list
		// of the required ones that we don't have values for so that we can
		// prompt the user.

		askVars := make([]string, 0)
		for _, p := range vars {
			_, foundInUser := pkg.Lookup(p, clusterVars)
			_, foundInCred := pkg.Lookup(p, credVars)
			_, foundInEnv := os.LookupEnv(p.Name)
			if !foundInUser && !foundInCred && !foundInEnv && len(p.Value) == 0 {
				logger.Debugf("Found no existance of <%s>, adding to list of required vars", p.Name)
				askVars = append(askVars, p.Name)
			}
		}

		// Okay, so now I have the required vars and I can now build up the prompts
		// to ask my user for the values.
		builder := prompt.NewPromptBuilder()

		rootQuestion, err := builder.Path("proceed").
			Text(fmt.Sprintf("This will deploy the solution %s to cluster %s; continue?", sol, cluster)).
			WithOptions(prompt.YesNo()).
			Build()

		for _, v := range askVars {
			logger.Tracef("Building prompt for <%s>", v)
			subP, _ := prompt.NewPromptBuilder().
				Path(v).
				Textf("What value would you like to use for '%s'?", v).
				Build()
			rootQuestion.AddSubPrompt(subP)
		}

		nextPrompter := rootQuestion.Itr()

		for p := nextPrompter(); p != nil; p = nextPrompter() {
			logger.Tracef("Asking <%s>", p.String())
			err = prompt.Ask(p, os.Stdout, os.Stdin)
			if err != nil {
				return err
			}
		}

		jenkinsBuildProjectCmd = &pkg.ServiceClient{
			Method:             "POST",
			ContentType:        "application/x-www-form-urlencoded",
			BaseURL:            fmt.Sprintf("%s/job/%s/buildWithParameters", viper.GetString("ci.api.url"), sol),
			AuthHandler:        pkg.BasicAuthHandler(viper.GetString("ci.api.user"), viper.GetString("ci.api.token")),
			ExpectedStatusCode: 201,
			QParams: func() map[string]string {
				m := make(map[string]string)
				m["token"] = viper.GetString("ci.buildtoken")
				return m
			},
			FParams: func() map[string]string {
				m := make(map[string]string)
				for _, p := range vars {
					envVal, ok := os.LookupEnv(p.Name)
					if ok {
						logger.Tracef("Using build parameter <%s> from environment with value <%s>.", p.Name, envVal)
						m[p.Name] = fmt.Sprintf("%v", envVal)
					}
				}
				for k, v := range rootQuestion.VarMap() {
					logger.Tracef("Adding build parameter <%s> with value <%s>.", k, v)
					m[k] = v
				}
				for k, v := range clusterVars {
					logger.Tracef("Adding build parameter <%s> with value <%s>.", k, v)
					m[k] = v
				}
				for k, v := range credVars {
					logger.Tracef("Adding build parameter <%s> with value <%s>.", k, v)
					m[k] = v
				}
				return m
			},
		}

		err = pkg.Exec(jenkinsBuildProjectCmd)
		if err != nil {
			return err
		}

		logger.Infof("Started deployment pipeline for solution %s...", sol)
	}

	return nil
}

func withEnvUpdates(svc *pkg.Service, ctx *atkmod.RunContext, runner *atkmod.CliModuleRunner) bool {
	// Update the service with the API key
	runner.WithEnvvar("JENKINS_API_TOKEN", viper.GetString("ci.api.token"))
	return pkg.StartHandler(svc, ctx, runner)
}

type crumbIssuerResponse struct {
	Crumb  string `json:"crumb"`
	cookie string `json:"-"`
}

type tokenData struct {
	TokenName  string `json:"tokenName"`
	TokenUuid  string `json:"tokenUuid"`
	TokenValue string `json:"tokenValue"`
}

type generateNewTokenResponse struct {
	Status string
	Token  tokenData `json:"data"`
}

// initTokenAndSave uses the builder (Jenkins) API to create an API key for the
// configured user, which is a bit inconvenient but is required for local
// execution.
func initTokenAndSave(svc *pkg.Service, ctx *atkmod.RunContext, runner *atkmod.CliModuleRunner) bool {
	for i := 1; i < 5; i++ {
		ctx.Log.Trace("Waiting for Jenkins to become available...")
		time.Sleep(time.Second * 30)
		resp, err := http.Get(svc.URL.String())
		if err != nil {
			ctx.AddError(err)
			return false
		}
		status := resp.StatusCode
		if status != 503 {
			break
		}
	}

	// TODO: this is going to get a little hacky, but that's OK for now...
	user := viper.GetString("ci.api.user")
	password := viper.GetString("ci.api.password")
	crumbInfo, err := getJenkinsCrumb(svc.URL, user, password, ctx)
	if err != nil {
		ctx.AddError(err)
		return false
	}
	ctx.Log.Tracef("Using crumb data: %v", crumbInfo)
	apiKey, err := createApiKey(svc.URL, user, password, crumbInfo, ctx)
	if err != nil {
		ctx.AddError(err)
		return false
	}
	ctx.Log.Infof("Succesfully created API token <%s> for user <%s>", apiKey.Token.TokenValue, user)
	viper.Set("ci.api.token", apiKey.Token.TokenValue)
	err = viper.WriteConfig()
	if err != nil {
		ctx.AddError(err)
		return false
	}
	ctx.Log.Infof("Succesfully wrote API token to configuration file.")

	return true
}

// getJenkinsCrumb gets the crumb information from the crumbIssuer endpoint, which
// can then be used to create an API key for the configured bifrost user. This is
// a little convoluted and it would have been nice to re-use one of the other existing
// functions, but this needed some special header handling stuff. Maybe a refactor
// that allows us to inject the handling of the response... ?
func getJenkinsCrumb(url *url.URL, user string, password string, ctx *atkmod.RunContext) (*crumbIssuerResponse, error) {
	ctx.Log.Trace("Calling crumbIssuer to get crumb data from Jenkins...")
	client := &http.Client{}
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/crumbIssuer/api/json", url), nil)
	if err != nil {
		return nil, err
	}
	authS := fmt.Sprintf("%s:%s", user, password)
	sEnc := b64.StdEncoding.EncodeToString([]byte(authS))
	req.Header.Set("Authorization", "Basic "+sEnc)
	resp, err := client.Do(req)
	ctx.Log.Tracef("Response received; got %d", resp.StatusCode)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("error while trying to generate API token for user %s: %v", user, resp.Status)
	}
	var issuerResp crumbIssuerResponse
	json.NewDecoder(resp.Body).Decode(&issuerResp)
	issuerResp.cookie = resp.Header.Get("Set-Cookie")

	return &issuerResp, nil
}

func createApiKey(url *url.URL, user string, password string, info *crumbIssuerResponse, ctx *atkmod.RunContext) (*generateNewTokenResponse, error) {
	ctx.Log.Trace("Calling generateNewToken to generate API token in Jenkins...")
	client := &http.Client{}
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/user/%s/descriptorByName/jenkins.security.ApiTokenProperty/generateNewToken", url, user), strings.NewReader("newTokenName=bifrost-generated-token"))
	if err != nil {
		return nil, err
	}
	req.SetBasicAuth(user, password)
	req.Header.Set("Cookie", info.cookie)
	req.Header.Set("Jenkins-Crumb", info.Crumb)
	ctx.Log.Tracef("Using url to generate token: %s", req.URL)
	resp, err := client.Do(req)
	ctx.Log.Tracef("Response received; got %d", resp.StatusCode)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("error while trying to generate API token for user %s: %v", user, resp.Status)
	}

	var tokenResponse generateNewTokenResponse
	err = json.NewDecoder(resp.Body).Decode(&tokenResponse)
	return &tokenResponse, err
}
