package pkg

import (
	"fmt"
	"io"
	"net/http"
)

func ReadHttpGet(url string, token string) ([]byte, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", fmt.Sprintf("%s", url), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+token)
	resp, err := client.Do(req)
	defer resp.Body.Close()
	return io.ReadAll(resp.Body)
}
