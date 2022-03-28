package server

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
)

type httpServer struct {
	Log *Log
}

// exported function so we can establish a http server in a seperate file
func NewHTTPServer(addr string) *http.Server {
	httpsrv := newHttpServer()
	r := mux.NewRouter()

	r.HandleFunc("/", httpsrv.handleConsume).Methods("GET")
	r.HandleFunc("/", httpsrv.handleProduce).Methods("POST")

	return &http.Server{
		Addr:    addr,
		Handler: r,
	}

}

func newHttpServer() *httpServer {
	return &httpServer{
		Log: NewLog(),
	}
}

// Defines what will be sent to the api routes to the handlers
type ProduceRequest struct {
	Record Record `json:"record"`
}

type ProduceResponse struct {
	Offset uint64 `json:"offset"`
}

type ConsumeRequest struct {
	Offset uint64 `json:"offset"`
}

type ConsumeResponse struct {
	Record Record `json:"record"`
}

// Define our api handlers that will be called by the routes

// Handle a new entry to the log (produce a new log)
func (s *httpServer) handleProduce(w http.ResponseWriter, r *http.Request) {
	var req ProduceRequest
	// Decode the expected body to a new product request struct
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// The request body has successfully been read to a new instance of the struct (validated), so add it to the log
	off, err := s.Log.Append(req.Record)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// Define the response from the struct
	res := ProduceResponse{Offset: off}
	err = json.NewEncoder(w).Encode(res)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// Handle a request to get an existing log offset
func (s *httpServer) handleConsume(w http.ResponseWriter, r *http.Request) {
	var req ConsumeRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// Validated into a new struct, so pass to service method
	record, err := s.Log.Read(req.Offset)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	res := ConsumeResponse{Record: record}
	err = json.NewEncoder(w).Encode(res)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
