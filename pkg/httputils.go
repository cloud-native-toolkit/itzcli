package pkg

import (
	"bytes"
	b64 "encoding/base64"
	"fmt"
	logger "github.com/sirupsen/logrus"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
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

// Some default handlers

// BasicAuthHandler
func BasicAuthHandler(user string, password string) AuthHandlerFunc {
	return func(req *http.Request) error {
		req.SetBasicAuth(user, password)
		return nil
	}
}

// BearerAuthTandler
func BearerAuthTandler(token string) AuthHandlerFunc {
	return func(req *http.Request) error {
		req.Header.Set("Authorization", fmt.Sprintf("%s %s", Bearer, token))
		return nil
		return nil
	}
}

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

func Exec(svc *ServiceClient) error {
	client := &http.Client{}

	// Add form parameters, if there are any.
	var reqForm url.Values
	var body io.Reader

	if svc.FParams != nil {
		reqForm = make(url.Values)
		fParams := svc.FParams()
		if len(fParams) > 0 {
			if svc.ContentType == ContentTypeMultiPart {
				buf := new(bytes.Buffer)
				w := multipart.NewWriter(buf)
				w.SetBoundary(Boundary)
				for k, v := range fParams {
					w.WriteField(k, v)
				}
				err := w.Close()
				if err != nil {
					return err
				}
				body = strings.NewReader(buf.String())
			} else {
				for k, v := range fParams {
					reqForm[k] = []string{v}
				}
				body = strings.NewReader(reqForm.Encode())
			}
		}
		logger.Tracef("Adding form body: %v", reqForm.Encode())
	} else {
		body = svc.Body
	}

	req, err := http.NewRequest(svc.Method, svc.BaseURL, body)
	if err != nil {
		return err
	}

	if svc.AuthHandler != nil {
		err = svc.AuthHandler(req)
		if err != nil {
			return err
		}
	}

	// Add query string parameters, if there are any.
	if svc.QParams != nil {
		qParams := svc.QParams()
		if len(qParams) > 0 {
			q := req.URL.Query()
			for k, v := range qParams {
				q.Add(k, v)
			}
			req.URL.RawQuery = q.Encode()
		}
	}

	logger.Debugf("Calling %s %s", req.Method, req.URL.String())
	if len(req.Form) > 0 {
		logger.Tracef("Using form values: %v", req.Form)
	}

	if svc.ContentType != "" {
		req.Header.Set("Content-Type", svc.ContentType)
	}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	if svc.ResponseHandler != nil {
		err = svc.ResponseHandler(resp.Body)
		if err != nil {
			return err
		}
	}

	if svc.ExpectedStatusCode != 0 {
		if resp.StatusCode != svc.ExpectedStatusCode {
			return fmt.Errorf("expected status code %d, got %d", svc.ExpectedStatusCode, resp.StatusCode)
		}
	}
	return nil
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
