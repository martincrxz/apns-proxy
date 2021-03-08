package net

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

const (
	deviceIDVarName = "deviceToken"
	servicePath     = "/3/device/{{{.deviceIDVarName}}}"
)

// Server listens to not HTTP/2 requests
type Server struct {
	router *mux.Router
}

// NewServer creates and returns a new server
func NewServer() *Server {
	server := &Server{
		router: mux.NewRouter(),
	}

	server.router.HandleFunc(servicePath, server.processNotification).Methods("POST")
	return server
}

// Run puts the server to listen at the specified port
func (server *Server) Run(p string) error {
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
	deviceId := vars[deviceIDVarName]

	client := getClient()
	client.SendNotification()
}
