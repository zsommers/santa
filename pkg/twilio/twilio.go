package twilio

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

func mustHaveEnv(e string) (string, error) {
	if v, ok := os.LookupEnv(e); ok {
		return v, nil
	}
	return "", fmt.Errorf("%s not set", e)
}

// Texter keeps track of Twillio account info
type Texter struct {
	accountSID    string
	authToken     string
	sendingNumber string
	url           string
	reallySend    bool
}

type twilioResponse struct {
	SID string `json:"sid"`
}

// NewTexter creates a new Texter from using environment variables, or errors
func NewTexter(reallySend bool) (*Texter, error) {
	var s, t, n string
	var err error
	if s, err = mustHaveEnv("TWILIO_ACCOUNT_SID"); err != nil {
		return nil, err
	}
	if t, err = mustHaveEnv("TWILIO_AUTH_TOKEN"); err != nil {
		return nil, err
	}
	if n, err = mustHaveEnv("TWILIO_SENDING_NUMBER"); err != nil {
		return nil, err
	}

	u := fmt.Sprintf("https://api.twilio.com/2010-04-01/Accounts/%s/Messages.json", s)

	return &Texter{
		accountSID:    s,
		authToken:     t,
		sendingNumber: n,
		url:           u,
		reallySend:    reallySend,
	}, nil
}

// SendText sends the message to the destination number
func (t *Texter) SendText(destination, message string) error {
	reader := t.buildPayload(destination, message)

	req, err := t.buildRequest(reader)
	if err != nil {
		return err
	}

	if t.reallySend {
		if err := t.sendRequest(req); err != nil {
			return err
		}
	} else {
		log.Printf("Would have sent %s the following\n\t%s", destination, message)
	}
	return nil
}

func (t *Texter) buildPayload(destination, message string) *strings.Reader {
	msgData := url.Values{}
	msgData.Set("To", destination)
	msgData.Set("From", t.sendingNumber)
	msgData.Set("Body", message)
	return strings.NewReader(msgData.Encode())
}

func (t *Texter) buildRequest(r *strings.Reader) (*http.Request, error) {
	req, err := http.NewRequest("POST", t.url, r)
	if err != nil {
		return nil, err
	}
	req.SetBasicAuth(t.accountSID, t.authToken)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	return req, nil
}

func (t *Texter) sendRequest(r *http.Request) error {
	client := &http.Client{
		Timeout: 60 * time.Second,
	}
	resp, err := client.Do(r)
	if err != nil {
		return err
	}

	log.Printf("Status Code: %d", resp.StatusCode)
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		var responseData twilioResponse
		decoder := json.NewDecoder(resp.Body)
		if err := decoder.Decode(&responseData); err != nil {
			return err
		}
		log.Printf("Response SID: %s", responseData.SID)
	} else {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		return fmt.Errorf("sending failed: status %d - %s", resp.StatusCode, body)
	}
	return nil
}
