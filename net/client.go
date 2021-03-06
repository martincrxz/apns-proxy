package net

import (
	"crypto/tls"
	"crypto/x509"
	"net/http"
	"net/url"
	"strconv"

	"github.com/rs/zerolog/log"
	"golang.org/x/net/http2"
)

type Client = http.Client

// NewClient creates and returns a new client
func NewClient(certFile, keyFile, proxy string) (*http.Client, error) {

	transport := &http.Transport{}

	if certFile != "" && keyFile != "" {
		cert, err := tls.LoadX509KeyPair(certFile, keyFile)
		if err != nil {
			log.Err(err).Msg("could not load certificate")
			return nil, err
		}
		x509Cert, err := x509.ParseCertificate(cert.Certificate[0])
		if err != nil {
			log.Err(err).Msg("could not load certificate")
			return nil, err
		}
		log.Info().Msg("certificate loaded, issuer common name: " + x509Cert.Issuer.CommonName)

		tlsConfig := &tls.Config{
			Certificates: []tls.Certificate{cert},
		}

		tlsConfig.BuildNameToCertificate()

		transport.TLSClientConfig = tlsConfig
	}

	if proxy != "" {
		proxyUrl, err := url.Parse(proxy)
		if err != nil {
			log.Error().Err(err).Msg("invalid proxy")
			return nil, err
		} else {
			log.Info().Msg("using proxy: " + proxy)
		}
		transport.Proxy = http.ProxyURL(proxyUrl)
	} else {
		log.Info().Msg("using no proxy")
	}

	if err := http2.ConfigureTransport(transport); err != nil {
		log.Error().Err(err).Msg("could not configure transport for HTTP/2 connections")
		return nil, err
	} else {
		log.Info().Msg("transport configured for HTTP/2 connections")
	}

	return &http.Client{Transport: transport}, nil
}

// SendNotification sends a notification using the specified client
func SendNotification(client *Client, req *http.Request) (*http.Response, error) {
	resp, err := client.Do(req)
	if err != nil {
		log.Err(err).Msg("could not get apple response")
		return nil, err
	}
	log.Info().Msg("apple response with status code " + strconv.Itoa(resp.StatusCode))
	return resp, nil
}
