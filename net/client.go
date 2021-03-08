package net

import (
	"encoding/json"
	"net/http"
)

const (
	apnsURL = "api.push.apple.com:443"
)

// Client sends HTTP/2 requests to APNS server
type Client struct {
	httpClient http.Client
}

// SendNotification sends a notification using the specified client
func (client *Client) SendNotification() {
	req, err := http.NewRequest("POST", apnsURL, nil)
	if err != nil {
		// handle error
	}
	resp, err := client.httpClient.Do(req)
	if err != nil {
		// handle error
	}
	defer resp.Body.Close()

	var data APNSResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		// handle error
	}
}

func getClient() *Client {
	return nil
}
