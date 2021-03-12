package net

import (
	"crypto/tls"
	"crypto/x509"
	"net/http"
	"strconv"
	"sync"

	"github.com/rs/zerolog/log"
	"golang.org/x/net/http2"
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
		log.Err(err).Msg("could not load certificate")
	}
	x509Cert, err := x509.ParseCertificate(cert.Certificate[0])
	log.Info().Msg("certificate loaded, issuer common name: " + x509Cert.Issuer.CommonName)

	tlsConfig := &tls.Config{
		Certificates:       []tls.Certificate{cert},
		InsecureSkipVerify: true,
	}

	tlsConfig.BuildNameToCertificate()

	transport := &http2.Transport{TLSClientConfig: tlsConfig}

	return &Client{
		httpClient: http.Client{Transport: transport},
	}
}

// SendNotification sends a notification using the specified client
func (client *Client) SendNotification(req *http.Request) (*http.Response, error) {
	resp, err := client.httpClient.Do(req)
	if err != nil {
		log.Err(err).Msg("could not get apple response")
		return nil, err
	}
	log.Info().Msg("apple response with status code " + strconv.Itoa(resp.StatusCode))
	return resp, nil
}

// ClientsPool is a pool of HTTP/2 clients than can comunicate with APNS server
type ClientsPool struct {
	pool *sync.Pool
}

// NewClientsPool creates and returns a new clients pool
func NewClientsPool(num int, certFile, keyFile string) *ClientsPool {
	log.Info().Msg("creating clients pool with " + strconv.Itoa(num) + " clients")
	clientsPool := ClientsPool{
		pool: &sync.Pool{},
	}

	for i := 0; i < num; i++ {
		clientsPool.pool.Put(NewClient(certFile, keyFile))
	}

	return &clientsPool
}

// GetClient removes and returns a random client from the clients pool
func (clientsPool *ClientsPool) GetClient() *Client {
	return clientsPool.pool.Get().(*Client)
}

// GetClientBack puts a client back in the clients pool
func (clientsPool *ClientsPool) GetClientBack(client *Client) {
	clientsPool.pool.Put(client)
}
