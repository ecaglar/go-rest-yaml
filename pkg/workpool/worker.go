package workpool

import (
	"../context"
	"../logger"
	"github.com/google/uuid"
)

//Worker defines a worker unit which can be assigned "Work"
//through its work channel where worker can pick it up.
//worker should also be aware of workerQueue so that
//it can notify it whenever it is available for the next work
//Worker has also an ID and access to context so that it can use
//storage and logger.
type Worker struct {
	workerQueue chan chan WorkRequest
	work        chan WorkRequest
	Ctx         *context.AppContext
	quit        chan bool
	ID          uuid.UUID
}

//NewWorker creates a worker instance
func NewWorker(workerQueue chan chan WorkRequest, ctx *context.AppContext) *Worker {
	return &Worker{
		workerQueue: workerQueue,
		work:        make(chan WorkRequest),
		quit:        make(chan bool),
		Ctx:         ctx,
		ID:          uuid.New(),
	}
}

//start function makes worker available to work on WorkRequests
//it start with worker notifies worker queue that it is available to execute
//So whenever worker queue has a work item to be able to work on it,
//it has been assigned to worker's work channel by dispatcher so that
//worker can pick it up and start working on that.
func (w *Worker) start() {
	go func() {

		for {
			//registers itself to worker queue
			w.workerQueue <- w.work
			w.Ctx.Logger.Log(logger.INFO, "Worker ", w.ID.String(), " registered its own chan WorkRequest to worker queue to say it is available")

			//after notifiying worker queue about availability, waits for an assignment by its own work queue

			select {
			case job := <-w.work:
				w.Ctx.Logger.Log(logger.INFO, "Work ", job.ID.String(), " has been assigned to worker s queue.")
				w.Ctx.Storage.Insert(job.Payload.Version, job.Payload)

			case <-w.quit:
				return
			}
		}

	}()
}

//stop terminates that worker so that it no task picked by it.
func (w *Worker) stop() {
	go func() { w.quit <- true }()
}
