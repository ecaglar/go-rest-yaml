/*
Package workpool implements a worker thread-pool approach to handle POST requests
Actors for this approach are:

WorkRequest - Work item that can processed by a worker. It is request payload in our case.
Worker      - It can process WorkItems assigned by dispatcher using Workers own work queue.
  			  Worker is responsible for registering itself to WorkerQueue
WorkerQueue - It is a buffered channel of channels. Workers use the channels goes into this channel to retrieve  works
WorkQueue 	- WorkRequests are being pushed to that queue so that dispatcher can pick it up and assign to workers.

Especially under heavy load, (e.g. 1M per minute) this method works quite effective and decrease latency / delay dramatically

Server is responsible for creating teh WorkRequest channel and starting dispatcher.
*/

package workpool

import (
	"../context"
	"../logger"
)

type Dispatcher struct {
	WorkerQueue chan chan WorkRequest
	WorkQueue   chan WorkRequest
	Ctx         *context.AppContext
	MaxWorkers  int
}

//NewDispatcher creates the WorkerQueue using max worker number received as argument.
//It also initialize context and work queue.
func NewDispatcher(workQueue chan WorkRequest, maxWorkers int, ctx *context.AppContext) *Dispatcher {

	WorkerQueue := make(chan chan WorkRequest, maxWorkers)

	return &Dispatcher{
		WorkerQueue: WorkerQueue,
		WorkQueue:   workQueue,
		Ctx:         ctx,
		MaxWorkers:  maxWorkers,
	}
}

//StartDispatcher creates the workers and starts them then starts dispatching.
func (d *Dispatcher) StartDispatcher() {

	//First create workers and make them available to work!
	for i := 0; i < d.MaxWorkers; i++ {
		worker := NewWorker(d.WorkerQueue, d.Ctx)
		worker.start()
	}

	go func() {
		for {
			select {

			//Dispatcher checks work queue and whenever it receives a work (which is handled and passed by http handler)
			//it just starts a new goroutine in order not to wait for worker queue for available workers.
			case work := <-d.WorkQueue:

				d.Ctx.Logger.Log(logger.INFO, "Work ", work.ID.String(), " received from WorkQueue", " version: ", work.Payload.Version)
				go func() {

					//get a available worker which can work on this
					worker := <-d.WorkerQueue

					d.Ctx.Logger.Log(logger.INFO, "Available Worker channel received from WorkerQueue")

					//dispatch the job to available worker.
					worker <- work

				}()
			}
		}
	}()

}
