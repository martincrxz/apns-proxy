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

	port := flag.String("p", defaultPort, "port number")
	clientCertFile := flag.String("cc", "", "client certificate file")
	clientKeyFile := flag.String("ck", "", "client key file")
	serverCertFile := flag.String("sc", "", "server certificate file")
	serverKeyFile := flag.String("sk", "", "server key file")
	proxy := flag.String("x", "", "proxy url")
	logFilePath := flag.String("g", "", "log file")
	flag.Parse()

	logFile := os.Stdout
	if *logFilePath != "" {
		var err error
		logFile, err = os.OpenFile(*logFilePath, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0660)
		if err != nil {
			fmt.Printf("could not create log file, error: %v\n", err)
			os.Exit(exitFailCode)
		}
	}

	log.Logger = zerolog.New(logFile).With().Timestamp().
		Logger().With().Caller().Logger()

	if *clientCertFile == "" {
		log.Warn().Msg("no client cert file specified")
	}

	if *clientKeyFile == "" {
		log.Warn().Msg("no client key file specified")
	}

	if *serverCertFile == "" {
		log.Warn().Msg("no server cert file specified")
	}

	if *serverKeyFile == "" {
		log.Warn().Msg("no server key file specified")
	}

	client := net.NewClient(*clientCertFile, *clientKeyFile, *proxy)
	server := net.NewServer(client)
	if err := server.Run(*port, *serverKeyFile, *serverCertFile); err != nil {
		log.Error().Err(err).Msg("error while running server")
		os.Exit(exitFailCode)
	}
}
