package worker

import (
	"fmt"
	"time"

	"github.com/rs/zerolog/log"

	"devcircus.com/cerberus/pkg/target/request"
	"devcircus.com/cerberus/pkg/util/shell"
)

// Worker execution data
type Worker struct {
	rConfig        request.Config
	requestChannel chan bool
}

// NewWorker create a new instance worker
func NewWorker(data request.Config) *Worker {
	w := Worker{}
	w.rConfig = data
	w.requestChannel = make(chan bool)
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
		<-w.requestChannel
		//throttle <- 1
		log.Debug().Msgf("Performing request: %s %s", w.rConfig.RequestType, w.rConfig.URL)

		reqErr := request.PerformRequest(w.rConfig, nil)

		if reqErr != nil {
			log.Warn().Msgf("Error requesting: %s %s", w.rConfig.RequestType, w.rConfig.URL)
			stdout, stderr, _ := shell.RunCommand("ls", "-ltr")
			fmt.Println("--- stdout ---")
			fmt.Println(stdout)
			fmt.Println("--- stderr ---")
			fmt.Println(stderr)
		} else {
			log.Info().Msgf("Epic win requesting: %s %s", w.rConfig.RequestType, w.rConfig.URL)
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
			w.requestChannel <- true
		case <-quit:
			ticker.Stop()
			return
		}
	}
}
