package request

import (
	"app/src/utils/debug"
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"os"
)

// Request wrapper
func Do(token, method, url string, payload, unmarshal interface{}) string {
	// Marshal payload to bytes
	var body io.ReadWriter

	if payload != nil {
		buf, err := json.Marshal(payload)
		if err != nil {
			panic(err)
		}
		body = bytes.NewBuffer(buf)
	}

	// Create client
	client := &http.Client{}

	// Create request
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		panic(err)
	}

	// Set content type
	req.Header.Set("Content-Type", "application/json")

	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	// Fetch request
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	// Read response body
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	// Display results
	if os.Getenv("Debug") == "true" {
		debug.Info(resp, respBody)
	}

	if unmarshal != nil {
		err = json.Unmarshal([]byte(respBody), &unmarshal)
		//err = json.NewDecoder(resp.Body).Decode(&unmarshal)
		if err != nil {
			panic(err)
		}
	}

	// Return string body
	return string(respBody)
}
