package server

import (
	"../memstore"
	"../logger"
	"../model"
	"../validator"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"net/http"
)

type validatorFunc func(h *http.Request) (bool, string)
type middleware func(h http.HandlerFunc) http.HandlerFunc


//Shared dependencies, better to pass lots of parameters to handlers
type Server struct {
	DB      memstore.Database
	Log     *logger.ExtLogger
	Routers *mux.Router
}

func CreateServer(l *logger.ExtLogger, database memstore.Database) *Server {
	server := &Server{
		Log:     l,
		Routers: mux.NewRouter(),
		DB:      database,
	}
	server.routes()
	return server
}

/**
We have one resource which is app metadata. So conceptially /apps can return all app metadata
This means that, we can filter it with query parameters. It also prevent us to define paths
for allpossible combinations

so one path with multiple search params can handle all combinations

we are basically "filtering" which means we are trying to find a subset amoung all
that is why all os these should be search parameters so we should nt send them via body or POST
etc

/api/v1/apps?version=1.0.0&license=Apache-2.0

I dont want to inject search params as json or yaml inside body and send with POST as,
1-cache issues
2-cannot be bookmarked
3-url has 2000 char capasity so we still have ways
*/

func (s *Server) routes() {
	s.Routers.HandleFunc("/api/v1/apps", s.Chain(s.createAppMetadataHandler,
		s.withValidation(validator.ValidateRequest),
		s.withLog())).Methods("POST")

	s.Routers.HandleFunc("/api/v1/apps", s.Chain(s.searchAppMetadataHandler,
		s.withLog())).Methods("GET")

}

func (s *Server) searchAppMetadataHandler(w http.ResponseWriter, r *http.Request) {

	queryStr := r.URL.Query() //map[string][]string
	result := s.DB.ReadWithParams(queryStr)
	if r.Header.Get("Accept") == "application/json" {
		json.NewEncoder(w).Encode(result)
	} else {
		yaml.NewEncoder(w).Encode(result)
	}
}

func (s *Server) createAppMetadataHandler(w http.ResponseWriter, r *http.Request) {
	var bodyBytes []byte
	if r.Body != nil {
		bodyBytes, _ = ioutil.ReadAll(r.Body)
	}
	// Restore the io.ReadCloser to its original state
	r.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))
	bodyString := string(bodyBytes)

	s.Log.LogInfo("body ", bodyString)

	var m = model.Metadata{}
	err := yaml.Unmarshal([]byte(bodyString), &m)
	if err != nil {
		s.Log.LogError(err.Error())
		fmt.Fprintf(w, "%s", err.Error())
		return
	}
	s.DB.Insert(m.Version, m)
}

func (s *Server) withValidation(validator validatorFunc) middleware {

	s.Log.LogInfo("withValidation called")

	return func(h http.HandlerFunc) http.HandlerFunc {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			s.Log.LogInfo("Inside withValidation / HandlerFunc")
			if isValid, errorStr := validator(r); !isValid {
				s.Log.LogWarning("Request is not valid - ", errorStr)
				fmt.Fprintf(w, "%s", "Request is not valid -", errorStr)
				return
			} else {
				h(w, r)
			}
		})

	}
}

func (s *Server) withLog() middleware {

	s.Log.LogInfo("withLog called")

	return func(h http.HandlerFunc) http.HandlerFunc {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			s.Log.LogInfo("Log Before")
			defer s.Log.LogInfo("Log After")
			h(w, r)
		})

	}

}

func (s *Server) Chain(h http.HandlerFunc, m ...middleware) http.HandlerFunc {

	if len(m) < 1 {
		return h
	}

	wrapper := h

	//in order to execute  same as given order
	for i := len(m) - 1; i >= 0; i-- {
		wrapper = m[i](wrapper)
	}
	return wrapper
}
