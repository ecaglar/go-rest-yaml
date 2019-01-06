package main

import (
	"../pkg/context"
	"../pkg/logger"
	"../pkg/memstore"
	"../pkg/server"
	"../pkg/workpool"
	"net/http"
)

const (
	MaxWorker = 3  //os.Getenv("MAX_WORKERS")
	MaxQueue  = 20 //os.Getenv("MAX_QUEUE")
)

func main() {

	//create async logger
	asyncLogger := logger.CreateAsyncLogger()

	storage := memstore.CreateInMemDB()
	storage.SetLogger(asyncLogger)

	//create application context
	appContext := context.AppContext{
		Storage: storage,
		Logger:  asyncLogger,
	}

	//initialize dispatcher and pools
	workQueue := make(chan workpool.WorkRequest, MaxQueue)
	dispatcher := workpool.NewDispatcher(workQueue, MaxWorker, &appContext)
	dispatcher.StartDispatcher()

	//create server
	server := server.CreateServer(&appContext, workQueue)

	http.Handle("/", server.Routers)
	asyncLogger.Log(logger.INFO, "Listening localhost 8080...")
	http.ListenAndServe("localhost:8080", server.Routers)

}
