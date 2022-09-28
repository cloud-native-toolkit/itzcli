package pkg

import (
	"bytes"
	b64 "encoding/base64"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type ServiceClientAuthType string

const (
	Bearer ServiceClientAuthType = "Bearer"
	Basic  ServiceClientAuthType = "Basic"
)

// ServiceClientOper is a client operation that provides more structure-driven
// interaction with the backend APIs so there don't have to be so many variations
// of HTTP methods.
type ServiceClientOper struct {
	Method   string
	URL      string
	Body     io.Reader
	Result   io.Writer
	AuthType ServiceClientAuthType
	User     string
	Password string
}

func ReadHttpGetT(url string, token string) ([]byte, error) {
	return readHttpGet(url, "Bearer "+strings.TrimSpace(token))
}

func ReadHttpGetB(url string, user string, password string) ([]byte, error) {
	data := fmt.Sprintf("%s:%s", user, password)
	sEnc := b64.StdEncoding.EncodeToString([]byte(data))
	return readHttpGet(url, "Basic "+sEnc)
}

func readHttpGet(url string, auth string) ([]byte, error) {
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
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("error while trying to communicate with server: %v", resp.Status)
	}

	return io.ReadAll(resp.Body)
}

// PostFileToURL posts the given file to the URL
func PostFileToURL(path string, url string) error {
	data, err := ReadFile(path)
	if err != nil {
		return err
	}
	req, err := http.Post(url, "application/zip", bytes.NewReader(data))

	if err != nil {
		return err
	}

	if req.StatusCode != 200 {
		return fmt.Errorf("error while trying to post %s to server: %v", path, req.StatusCode)
	}
	return nil
}

func PostToURLB(url string, user string, pass string, data []byte) error {
	client := &http.Client{}
	req, err := http.NewRequest("POST", fmt.Sprintf("%s", url), nil)
	if err != nil {
		return err
	}
	req.SetBasicAuth(user, pass)
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode > 400 {
		return fmt.Errorf("error while trying to post to <%s>: %v", url, resp.StatusCode)
	}
	return nil
}
