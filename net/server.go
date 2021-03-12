package net

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
)

const (
	deviceIDVarName = "deviceToken"
	servicePath     = "/3/device/{" + deviceIDVarName + "}"
	protocol        = "https"
	host            = "api.push.apple.com"
)

// Server listens to not HTTP/2 requests
type Server struct {
	router      *mux.Router
	clientsPool *ClientsPool
}

// NewServer creates and returns a new server
func NewServer(clientsPool *ClientsPool) *Server {
	server := &Server{
		router:      mux.NewRouter(),
		clientsPool: clientsPool,
	}

	server.router.HandleFunc(servicePath, server.processNotification).Methods("POST")
	return server
}

// Run puts the server to listen at the specified port
func (server *Server) Run(p string) error {
	log.Info().Msg("running server on port " + p)
	return http.ListenAndServe(":"+p, server.router)
}

func (server *Server) processNotification(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	vars := mux.Vars(r)
	deviceID := vars[deviceIDVarName]

	log.Info().Msg("sending notification to " + deviceID + " device")

	url := protocol + "://" + host + "/3/device/" + deviceID
	apnsRequest, err := http.NewRequest(http.MethodPost, url, r.Body)
	if err != nil {
		log.Err(err).Msg("could not build request")
	}

	for key, values := range r.Header {
		for _, value := range values {
			log.Info().Msg("bypassing header to apple, key: " + key + ", value: " + value)
			apnsRequest.Header.Add(key, value)
		}
	}

	client := server.clientsPool.GetClient()
	resp, err := client.SendNotification(apnsRequest)
	server.clientsPool.GetClientBack(client)
	defer resp.Body.Close()

	enc := json.NewEncoder(w)
	enc.SetIndent("", "\t")

	if resp == nil {
		w.WriteHeader(http.StatusInternalServerError)
		enc.Encode(ErrorMessage{Error: "could not get apple response"})
		return
	}

	for key, values := range resp.Header {
		for _, value := range values {
			log.Info().Msg("bypassing header from apple, key: " + key + ", value: " + value)
			w.Header().Add(key, value)
		}
	}

	if resp.StatusCode != http.StatusOK {
		var data APNSResponse
		if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			enc.Encode(ErrorMessage{Error: "could not decode apple response"})
			return
		}
		w.WriteHeader(resp.StatusCode)
		enc.Encode(data)
		return
	}

	w.WriteHeader(http.StatusOK)
	return
}
