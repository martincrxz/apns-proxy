package net

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
)

const (
	deviceIDVarName = "deviceToken"
	servicePath     = "/3/device/{{{.deviceIDVarName}}}"
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
	var data APNSRequest
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorMessage{Error: "json couldn't be parsed"})
		return
	}
	vars := mux.Vars(r)
	deviceID := vars[deviceIDVarName]

	client := server.clientsPool.GetClient()
	client.SendNotification(deviceID)
	server.clientsPool.GetClientBack(client)
}
