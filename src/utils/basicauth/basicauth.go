// Package basicauth provides basic authentication method (JWT token)
package basicauth

import (
	"app/src/utils/debug"
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

// Basic authentication for JWT token
func AuthForJWTToken(url, email, password string) string {

	//Encode the data
	type Payload struct {
		Email    string `json:"email,omitempty"`
		Password string `json:"password,omitempty"`
	}

	payload := &Payload{
		Email:    email,
		Password: password,
	}

	putBody, _ := json.Marshal(payload)
	requestBody := bytes.NewBuffer(putBody)

	// Create client
	client := &http.Client{}

	// Create request
	req, err := http.NewRequest(http.MethodPost, url, requestBody)
	if err != nil {
		log.Fatalln(err)
	}

	// Set content type
	req.Header.Set("Content-Type", "application/json")

	// Fetch Request
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalln(err)
	}
	defer resp.Body.Close()

	// Read Response Body
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	// Display Results
	debug.Info(resp, respBody)

	type Response struct {
		Token string `json:"access_token"`
	}

	var jwt Response
	err = json.Unmarshal(respBody, &jwt)
	if err != nil {
		log.Fatalln("Error! Check your Kitsu credentials in conf.toml")
		//log.Fatalln(err)
	}

	return jwt.Token
}
