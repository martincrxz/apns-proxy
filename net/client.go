package net

import (
	"crypto/tls"
	"encoding/json"
	"math/rand"
	"net/http"
	"strconv"

	"github.com/rs/zerolog/log"
)

const (
	apnsURL = "api.push.apple.com:443"
)

// Client sends HTTP/2 requests to APNS server
type Client struct {
	httpClient http.Client
}

// NewClient creates and returns a new client
func NewClient(certFile, keyFile string) *Client {

	if certFile == "" || keyFile == "" {
		return &Client{
			httpClient: http.Client{},
		}
	}

	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		// handle error
	}

	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
	}

	tlsConfig.BuildNameToCertificate()

	transport := &http.Transport{TLSClientConfig: tlsConfig}

	return &Client{
		httpClient: http.Client{Transport: transport},
	}
}

// SendNotification sends a notification using the specified client
func (client *Client) SendNotification(deviceID string) (int, string, error) {
	req, err := http.NewRequest("POST", apnsURL, nil)
	if err != nil {
		// handle error
	}
	resp, err := client.httpClient.Do(req)
	if err != nil {
		// handle error
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var data APNSResponse
		if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
			// handle error
		}
	}

	return 200, "", nil
}

// ClientsPool is a pool of HTTP/2 clients than can comunicate with APNS server
type ClientsPool struct {
	clients []*Client
}

// NewClientsPool creates and returns a new clients pool
func NewClientsPool(num int, certFile, keyFile string) *ClientsPool {
	log.Info().Msg("creating clients pool with " + strconv.Itoa(num) + " clients")
	clientsPool := ClientsPool{
		clients: []*Client{},
	}

	for i := 0; i < num; i++ {
		clientsPool.clients = append(clientsPool.clients, NewClient(certFile, keyFile))
	}

	return &clientsPool
}

// GetClient removes and returns a random client from the clients pool
func (clientsPool *ClientsPool) GetClient() *Client {
	index := rand.Intn(len(clientsPool.clients))
	client := clientsPool.clients[index]
	clientsPool.clients = append(clientsPool.clients[:index], clientsPool.clients[index+1:]...)
	return client
}

// GetClientBack puts a client back in the clients pool
func (clientsPool *ClientsPool) GetClientBack(client *Client) {
	clientsPool.clients = append(clientsPool.clients, client)
}
