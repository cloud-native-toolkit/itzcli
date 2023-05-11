package api

import (
	"bytes"
	"fmt"
	"html/template"
	"net/http"

	"github.com/cloud-native-toolkit/itzcli/pkg"
	"github.com/cloud-native-toolkit/itzcli/pkg/auth"
	"github.com/gin-gonic/gin"
	logger "github.com/sirupsen/logrus"
)

func StartServer() {
	r := SetUpRouter()
	r.Run("localhost:8795")
}

func SetUpRouter() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.GET("/login", GetTechZoneToken)
	return r
}


func GetTechZoneToken(c *gin.Context) {
	// Grab the access token
	access_token := c.Query("token")
	if access_token == "" {
		c.Data(http.StatusUnauthorized, "text/html; charset=utf-8", []byte(errorHTML))
		logger.Debug("Missing required access token...")
		auth.ErrorGettingToken()
		return
	}
	// Okay, we have the access token so let's make our API call to TechZone now
	accessTokenurl := fmt.Sprintf("https://auth.techzone.ibm.com/user?access_token=%s", access_token)
  
	techZoneData, err := pkg.ReadHttpGetTWithFunc(accessTokenurl, "", func(code int) error {
		logger.Debugf("Handling HTTP return code %d...", code)
		return nil
	})
	if err != nil {
		c.Data(http.StatusBadRequest, "text/html; charset=utf-8", []byte(errorHTML))
		auth.ErrorGettingToken()
		return
	}
	// parse the body and grab requestJson
	jsoner := auth.NewJsonReader()
	techZoneDataR := bytes.NewReader(techZoneData)
	requestJson, err := jsoner.Read(techZoneDataR)
	var buf bytes.Buffer
	if err != nil || requestJson.Token == "" {
		c.Data(http.StatusBadRequest, "text/html; charset=utf-8", []byte(errorHTML))
		auth.ErrorGettingToken()
		return
	}
	htmlParser, _ := template.New("itzcliapi").Parse(succesHTML)
	htmlParser.Execute(&buf, requestJson)
	htmlString := buf.String()
	// Serve the HTML page 
	c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(htmlString))
	// Write the token to the config file
	auth.SaveTokenToConfig(requestJson.Token)
}

