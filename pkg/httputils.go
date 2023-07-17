package pkg

import (
	b64 "encoding/base64"
	"fmt"
	logger "github.com/sirupsen/logrus"
	"io"
	"net/http"
	"strings"
)

type ServiceClientAuthType string

const (
	Bearer ServiceClientAuthType = "Bearer"
	Basic  ServiceClientAuthType = "Basic"
)

const (
	ContentTypeMultiPart string = "multipart/form-data; boundary=\"========\""
	Boundary             string = "========"
)

type ParamBuilderFunc func() map[string]string
type ResponseHandlerFunc func(reader io.ReadCloser) error
type AuthHandlerFunc func(req *http.Request) error

// ServiceClient is a client operation that provides more structure-driven
// interaction with the backend APIs so there don't have to be so many variations
// of HTTP methods.
type ServiceClient struct {
	Method             string
	BaseURL            string
	QParams            ParamBuilderFunc
	FParams            ParamBuilderFunc
	ResponseHandler    ResponseHandlerFunc
	AuthHandler        AuthHandlerFunc
	Body               io.Reader
	ExpectedStatusCode int
	ContentType        string
}

type ReturnCodeHandlerFunc func(code int) error

func ReadHttpGetT(url string, token string) ([]byte, error) {
	return ReadHttpGetTWithFunc(url, token, nil)
}

func ReadHttpGetTWithFunc(url string, token string, handler ReturnCodeHandlerFunc) ([]byte, error) {
	return readHttpGet(url, "Bearer "+strings.TrimSpace(token), handler)
}

func ReadHttpGetB(url string, user string, password string, handler ReturnCodeHandlerFunc) ([]byte, error) {
	return ReadHttpGetBWithFunc(url, user, password, nil)
}

func ReadHttpGetBWithFunc(url string, user string, password string, handler ReturnCodeHandlerFunc) ([]byte, error) {
	data := fmt.Sprintf("%s:%s", user, password)
	sEnc := b64.StdEncoding.EncodeToString([]byte(data))
	return readHttpGet(url, "Basic "+sEnc, handler)
}

func readHttpGet(url string, auth string, handler ReturnCodeHandlerFunc) ([]byte, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", fmt.Sprintf("%s", url), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", auth)
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	//goland:noinspection GoUnhandledErrorResult
	defer resp.Body.Close()

	logger.Tracef("Got response: %d", resp.StatusCode)
	if resp.StatusCode != 200 {
		logger.Trace("Preparing to call handler...")
		if handler != nil {
			logger.Trace("Calling handler...")
			return nil, handler(resp.StatusCode)
		}
		return nil, fmt.Errorf("error while trying to communicate with server: %v", resp.Status)
	}

	return io.ReadAll(resp.Body)
}
