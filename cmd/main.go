package main

import (
	"../pkg/logger"
	"../pkg/memstore"
	"../pkg/server"
	"net/http"
)

func main()  {

	asyncLogger  := logger.CreateAsyncLogger()
	db 		:= memstore.CreateInMemDB()
	server  := server.CreateServer(asyncLogger,db)

	http.Handle("/",server.Routers)
	asyncLogger.Log(logger.INFO,"Listening localhost 8080...")
	http.ListenAndServe("localhost:8080", server.Routers)

}
