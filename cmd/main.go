package main

import (
	"../pkg/logger"
	"../pkg/memstore"
	"../pkg/server"
	"net/http"
)

func main()  {

	logger  := logger.CreateLogger()
	db 		:= memstore.CreateInMemDB()
	server  := server.CreateServer(logger,db)

	http.Handle("/",server.Routers)
	logger.LogInfo("Listening localhost 8080...")
	http.ListenAndServe("localhost:8080", server.Routers)

}
