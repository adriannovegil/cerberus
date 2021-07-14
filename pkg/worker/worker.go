package worker

import (
	"time"

	"github.com/rs/zerolog/log"

	"devcircus.com/cerberus/pkg/requests"
)

// Worker execution data
type Worker struct {
	rConfig        requests.RequestConfig
	requestChannel chan requests.RequestConfig
}

// NewWorker create a new instance worker
func NewWorker(data requests.RequestConfig) *Worker {
	w := Worker{}
	w.rConfig = data
	w.requestChannel = make(chan requests.RequestConfig, 1)
	return &w
}

// Start the job
func (w *Worker) Start() {
	go w.doWork()
}

func (w *Worker) doWork() {
	//go w.listenToRequestChannel()
	go w.createTicker()
	for {
		select {
		case requect := <-w.requestChannel:
			//throttle <- 1
			log.Debug().Msgf("Performing request: %s %s", w.rConfig.RequestType, w.rConfig.URL)
			go requests.PerformRequest(requect, nil)
		}
	}
}

// createTicker. A time ticker writes data to request channel for every
// request.CheckEvery seconds
func (w *Worker) createTicker() {
	var ticker *time.Ticker = time.NewTicker(w.rConfig.CheckEvery * time.Second)
	quit := make(chan struct{})
	for {
		select {
		case <-ticker.C:
			w.requestChannel <- w.rConfig
		case <-quit:
			ticker.Stop()
			return
		}
	}
}
