package auth

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"time"

	"github.com/AlecAivazis/survey/v2"
	logger "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type TechZoneUser struct {
    Token string `json:"token"`
    Preferredfirstname string `json:"preferredfirstname"`
    Preferredlastname string `json:"preferredlastname"`
}

const tokenName = "reservations.api.token"

func GetToken() error {
	token := loadTokenFromConfig()
	// Check to see if the api token is set
	if token != "" {
		reValidate := true
		prompt := &survey.Confirm{
    		Message: "You're already logged into TechZone. Do you want re-authenticate?",
		}
		survey.AskOne(prompt, &reValidate)
		if !reValidate {
			return nil
		}
		// null out the existing token and let the user reauthenticate
		err := SaveTokenToConfig("")
  		if err != nil {
    		return err
  		}
		token = loadTokenFromConfig()
	}
	// Okay, open the browser so the user can authenticate
	openBrowser()
	// loop while the token is is empty
	// The API endpoint will handle updating the token so lets wait for the user
	// We also don't want this to run forever, so timeout after 5 mins
	timeOut := time.Now().Add(5 * time.Minute)
  	for {
    	token = loadTokenFromConfig()
		// If the token isn't empty anymore, then break out of the loop
		if token != "" {
			break
		}
		// Check the timeout
		if time.Now().After(timeOut) {
			return errors.New("Session timeout. Please try running the auth login command again.")
		}
  	}
	if token == "error" {
		SaveTokenToConfig("")
		return errors.New("There was an error trying to login to TechZone. Please try running the auth login command again.")
	}
	// We want to give about 3 seconds for the gin server to response with the HTML before we exit the code
	time.Sleep(3 * time.Second)
	fmt.Println("Authentication successful.")
	return nil
}

func loadTokenFromConfig() string {
	viper.ReadInConfig()
	return viper.GetString(tokenName)
}

func SaveTokenToConfig(value string) error {
	viper.Set(tokenName, value)
	return viper.WriteConfig()
}

func ErrorGettingToken() {
	SaveTokenToConfig("error")
}

func openBrowser() {
	techZoneLoginURL := "https://auth.techzone.ibm.com/login?callbackUrl=http://localhost:8080/login"
	fmt.Printf("Press enter to open %s in your browser... ", techZoneLoginURL)
	// Listen for the user to hit enter
	input := bufio.NewScanner(os.Stdin); input.Scan()
	// Okay, open the url
	exec.Command("open", techZoneLoginURL).Run()
	logger.Debugf("Waiting for user response...")
}


type JsonReader struct{}

func (j *JsonReader) Read(reader io.Reader) (TechZoneUser, error) {
	var res TechZoneUser
	err := json.NewDecoder(reader).Decode(&res)
	return res, err
}

func NewJsonReader() *JsonReader {
	return &JsonReader{}
}
