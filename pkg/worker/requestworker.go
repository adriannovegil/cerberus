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

// RequestWorker execution data
type RequestWorker struct {
	ID                  string
	tickerChannel       chan bool
	timeRecorderChannel chan bool
	MetricsRecorder     metrics.Recorder
	rConfig             request.Config
}

// NewRequestWorker create a new instance worker
func NewRequestWorker(id string, data request.Config) *RequestWorker {
	rw := RequestWorker{}
	rw.ID = id
	rw.rConfig = data
	rw.tickerChannel = make(chan bool)
	rw.timeRecorderChannel = make(chan bool)
	return &rw
}

// Start the job
func (rw *RequestWorker) Start(ctx context.Context) {
	go rw.doWork(ctx)
}

func (rw *RequestWorker) doWork(ctx context.Context) {
	go rw.StartTicker()
	go rw.StartTimeRecorder()
	metricsRecorder, _ := metrics.RecorderFromContext(ctx)

	for {
		<-rw.tickerChannel
		rw.timeRecorderChannel <- true
		//throttle <- 1

		log.Debug().Msgf("Performing request: %s %s", rw.rConfig.RequestType, rw.rConfig.URL)

		reqErr := request.PerformRequest(rw.rConfig, nil)

		metricsRecorder.IncRetry(rw.rConfig.ID)
		metricsRecorder.IncAttempt(rw.rConfig.ID)

		if reqErr != nil {
			log.Warn().Msgf("Error requesting: %s %s", rw.rConfig.RequestType, rw.rConfig.URL)

			for _, fallbackActionName := range rw.rConfig.Fallbacks {
				fallback.Execute(ctx,
					*config.GetFallbackCOnfigurationByName(fallbackActionName))
			}

		} else {
			log.Info().Msgf("Epic win requesting: %s %s", rw.rConfig.RequestType, rw.rConfig.URL)
		}
		rw.timeRecorderChannel <- true
	}
}

// StartTicker writes data to request channel for every x seconds
func (rw *RequestWorker) StartTicker() {
	var ticker *time.Ticker = time.NewTicker(rw.rConfig.CheckEvery * time.Second)
	quit := make(chan struct{})
	for {
		select {
		case <-ticker.C:
			rw.tickerChannel <- true
		case <-quit:
			ticker.Stop()
			return
		}
	}
}

// StartTimeRecorder start the time recorder channel
func (rw *RequestWorker) StartTimeRecorder() {
	for {
		<-rw.timeRecorderChannel
		start := time.Now()
		<-rw.timeRecorderChannel
		elapsed := time.Since(start)
		log.Debug().Msgf("Task executed in: %d miliseconds", elapsed.Nanoseconds()/1000000)
	}
}
