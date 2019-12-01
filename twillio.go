package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
)

func mustHaveEnv(n string) (string, error) {
	if v, ok := os.LookupEnv(n); ok {
		log.Printf("%s: %s\n", n, v)
		return v, nil
	}
	return "", fmt.Errorf("%s not set", n)
}

// Texter keeps track of Twillio account info
type Texter struct {
	accountSID    string
	authToken     string
	sendingNumber string
	url           string
}

// NewTexter creates a new Texter from using environment variables, or errors
func NewTexter() (*Texter, error) {
	var s, t, n string
	var err error
	if s, err = mustHaveEnv("TWILLIO_ACCOUNT_SID"); err != nil {
		return nil, err
	}
	if t, err = mustHaveEnv("TWILLIO_AUTH_TOKEN"); err != nil {
		return nil, err
	}
	if n, err = mustHaveEnv("TWILLIO_SENDING_NUMBER"); err != nil {
		return nil, err
	}

	u := fmt.Sprintf("https://api.twilio.com/2010-04-01/Accounts/%s/Messages.json", s)

	return &Texter{
		accountSID:    s,
		authToken:     t,
		sendingNumber: n,
		url:           u,
	}, nil
}

// SendText sends the message to the destination number
func (t *Texter) SendText(destination string, message string) error {
	msgData := url.Values{}
	msgData.Set("To", destination)
	msgData.Set("From", t.sendingNumber)
	msgData.Set("Body", message)
	msgDataReader := *strings.NewReader(msgData.Encode())
	client := &http.Client{}
	req, _ := http.NewRequest("POST", t.url, &msgDataReader)
	req.SetBasicAuth(t.accountSID, t.authToken)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	resp, err := client.Do(req)
	if err != nil {
		log.Panic(err)
	}

	log.Printf("Status Code: %d\n", resp.StatusCode)
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		var data map[string]interface{}
		decoder := json.NewDecoder(resp.Body)
		if err := decoder.Decode(&data); err != nil {
			return err
		}
		fmt.Println(data["sid"])
	} else {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		return fmt.Errorf("sending failed: status %d - %s", resp.StatusCode, body)
	}
	return nil
}
