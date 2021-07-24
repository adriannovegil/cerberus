package worker

import (
	"context"
	"time"

	"github.com/rs/zerolog/log"

	"devcircus.com/cerberus/pkg/config"
	"devcircus.com/cerberus/pkg/fallback"
	"devcircus.com/cerberus/pkg/metrics"
	"devcircus.com/cerberus/pkg/target/request"
)

// Worker execution data
type Worker struct {
	ID                  string
	rConfig             request.Config
	requestChannel      chan bool
	timeRecorderChannel chan bool
	MetricsRecorder     metrics.Recorder
}

// NewWorker create a new instance worker
func NewWorker(id string, data request.Config) *Worker {
	w := Worker{}
	w.ID = id
	w.rConfig = data
	w.requestChannel = make(chan bool)
	w.timeRecorderChannel = make(chan bool)
	return &w
}

// Start the job
func (w *Worker) Start(ctx context.Context) {
	go w.doWork(ctx)
}

func (w *Worker) doWork(ctx context.Context) {
	go w.createTicker()
	go w.timeRecorder()
	metricsRecorder, _ := metrics.RecorderFromContext(ctx)

	for {
		<-w.requestChannel
		w.timeRecorderChannel <- true
		//throttle <- 1

		log.Debug().Msgf("Performing request: %s %s", w.rConfig.RequestType, w.rConfig.URL)

		reqErr := request.PerformRequest(w.rConfig, nil)

		metricsRecorder.IncRetry(w.rConfig.ID)
		metricsRecorder.IncAttempt(w.rConfig.ID)

		if reqErr != nil {
			log.Warn().Msgf("Error requesting: %s %s", w.rConfig.RequestType, w.rConfig.URL)

			for _, fallbackActionName := range w.rConfig.Fallbacks {
				fallback.Execute(ctx,
					*config.GetFallbackCOnfigurationByName(fallbackActionName))
			}

		} else {
			log.Info().Msgf("Epic win requesting: %s %s", w.rConfig.RequestType, w.rConfig.URL)
		}
		w.timeRecorderChannel <- true
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

func (w *Worker) timeRecorder() {
	for {
		<-w.timeRecorderChannel
		start := time.Now()
		<-w.timeRecorderChannel
		elapsed := time.Since(start)
		log.Debug().Msgf("Task executed in: %d miliseconds", elapsed.Nanoseconds()/1000000)
	}
}
