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

//signature of the validation function which you can inject to handler to validate your request
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
We have one resource which is app metadata. So conceptually  /apps can return all app metadata
This means that, we can filter it with url query parameters. It also prevent us to define paths
for all possible combinations

so one path with multiple search params can handle all combinations

to get application metadata, we are basically "filtering" which means we are trying to find a subset among all
possible records. The best way to handle "filtering" is to define our GET method so that it can query data
using url search parameters,

If no search paramater passed via URL then it means server should return all data without filtering
TODO: One point to pay attention here it is not feasible to return all data especially if data is huge
TODO: so there should be additional solutions like "paging" in order to prevent performance issues

**GET - /api/v1/apps**
Returns all records

**GET - /api/v1/apps?version=1.0.0**
Returns the record with version 1.0.0

**GET - /api/v1/apps?version=1.0.0&title=my%20app**
Returns the record with version 1.0.0 if exists.Does not check other parameters as version number is unique.

**GET - /api/v1/apps?company=mycompany.com&title=my%20app**
Returns record(s) with company name "mycompany.com" and title **contains** "my app"

**GET - /api/v1/apps?description=latest**
Returns record(s) with description **contains** "latest"

**GET - /api/v1/apps?maintainers.name=Bill&maintainers.name=Joe**
Returns record(s) which have/has maintainers name "Bill" and "Joe"

**GET - /api/v1/apps?maintainers.email=bill@hotmail.com&license=Apache-2.1**
Returns record(s) which have/has maintainers email "bill@hotmail.com" with licence "Apache-2.1"


I dont want to inject search params as json or yaml inside body and send with POST due to following reasons,
1-cache issues
2-cannot be bookmarked
3-url has 2000 char capacity so we still have ways
*/
//routes inits handlers for mux
func (s *Server) routes() {
	s.Routers.HandleFunc("/api/v1/apps", s.Chain(s.createAppMetadataHandler,
		s.withValidation(validator.ValidateRequest),
		s.withLog())).Methods("POST")

	s.Routers.HandleFunc("/api/v1/apps", s.Chain(s.searchAppMetadataHandler,
		s.withLog())).Methods("GET")

}

//searchAppMetadataHandler returns the related records matching url query parameters
//if url query params are empty then returns all records
//Default content type is yaml. However, if client explicetly requires json format
//then server returns the response in json
func (s *Server) searchAppMetadataHandler(w http.ResponseWriter, r *http.Request) {

	queryStr := r.URL.Query() //map[string][]string
	result := s.DB.ReadWithParams(queryStr)
	if r.Header.Get("Accept") == "application/json" {
		s.Log.LogInfo("<-- appliation/json has been requested by client")
		json.NewEncoder(w).Encode(result)
	} else {
		yaml.NewEncoder(w).Encode(result)
	}
}

//searchAppMetadataHandler creates the appliation metadata sent via body payload
//supports both yaml and json payloads
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
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "%s", err.Error())
		return
	}
	s.DB.Insert(m.Version, m)
}

//withValidation middleware performs validation check on request body
func (s *Server) withValidation(validator validatorFunc) middleware {

	s.Log.LogInfo("withValidation called")

	return func(h http.HandlerFunc) http.HandlerFunc {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			if isValid, errorStr := validator(r); !isValid {
				s.Log.LogError("Request is not valid - ", errorStr)
				w.WriteHeader(http.StatusBadRequest)
				fmt.Fprintf(w, "%s", "Request is not valid -", errorStr)
				return
			} else {
				h(w, r)
			}
		})

	}
}

//withLog middleware logs messages for the handler.
func (s *Server) withLog() middleware {

	s.Log.LogInfo("withLog called")

	return func(h http.HandlerFunc) http.HandlerFunc {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			s.Log.LogInfo("<--",r.Method)
			s.Log.LogInfo("<--",r.Header.Get("Accept"))
			s.Log.LogInfo("<--",r.URL.Path)
			s.Log.LogInfo("<--",r.URL.RawQuery)
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
