package worker

import (
	"devcircus.com/cerberus/pkg/requests"
)

// Worker execution data
type Worker struct {
	rConfig requests.RequestConfig
}

// NewWorker create a new instance worker
func NewWorker(data requests.RequestConfig) *Worker {
	w := Worker{}
	w.rConfig = data
	return &w
}

// Start the job
func (w *Worker) Start() {
	go w.doWork()
}

func (w *Worker) doWork() {
	println("Working on: ", w.rConfig.RequestType, " ", w.rConfig.URL)
}
