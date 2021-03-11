package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/martincrxz/apns-proxy/net"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

const (
	defaultClientsNumber = 10
	defaultPort          = "8080"
	exitFailCode         = 1
)

func main() {

	logFile, err := os.OpenFile("./log/apns-proxy.log", os.O_RDWR|os.O_APPEND|os.O_CREATE, 0660)
	if err != nil {
		fmt.Printf("could not create log file, error: %v\n", err)
		os.Exit(exitFailCode)
	}

	log.Logger = zerolog.New(logFile).With().Timestamp().
		Logger().With().Caller().Logger()

	port := flag.String("p", defaultPort, "port number")
	clientsNumber := flag.Int("l", defaultClientsNumber, "number of clients")
	certFile := flag.String("c", "", "certificate file")
	keyFile := flag.String("k", "", "key file")
	flag.Parse()

	if *certFile == "" {
		log.Warn().Msg("no cert file specified")
	}

	if *keyFile == "" {
		log.Warn().Msg("no key file specified")
	}

	clientsPool := net.NewClientsPool(*clientsNumber, *certFile, *keyFile)
	server := net.NewServer(clientsPool)
	if err := server.Run(*port); err != nil {
		log.Error().Err(err).Msg("error while running server")
		os.Exit(exitFailCode)
	}
}
