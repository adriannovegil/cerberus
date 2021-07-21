package worker

import (
	"context"
	"net/http"
	"os"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog/log"

	"devcircus.com/cerberus/pkg/config"
	"devcircus.com/cerberus/pkg/execute"
	"devcircus.com/cerberus/pkg/metrics"
)

// Supervisor config data
type Supervisor struct {
	MetricsRecorder metrics.Recorder
}

// TickerTime is the time between supervisor checks
const TickerTime = 5

var (
	workers []Worker
)

// NewSupervisor create a new supervisor worker
func NewSupervisor() *Supervisor {
	s := Supervisor{}
	return &s
}

// Run launch the worker jobs
func (s *Supervisor) Run() {
	// Prometheus registry to expose metrics.
	promreg := prometheus.NewRegistry()
	go func() {
		http.ListenAndServe(":8081", promhttp.HandlerFor(promreg, promhttp.HandlerOpts{}))
	}()
	s.MetricsRecorder = metrics.NewPrometheusRecorder(promreg)

	defer func(start time.Time) {
		s.MetricsRecorder.ObserveCommandExecution(start, true)
	}(time.Now())

	ctx := metrics.SetRecorderOnContext(context.TODO(), s.MetricsRecorder)

	data := config.Config.Targets.Requests
	for i, requestConfig := range data {
		log.Debug().Msgf("Launching worker #%d: %s %s", i, requestConfig.RequestType, requestConfig.URL)
		w := NewWorker(requestConfig)
		workers = append(workers, *w)
		w.Start(ctx)
	}

LOOP:
	for {
		// Calling Sleep method
		time.Sleep(TickerTime * time.Second)
		select {
		case <-execute.Done:
			log.Info().Msg("Graceful termination")
			os.Exit(0)
		case <-execute.Stop:
			log.Warn().Msg("Process terminated by external signal")
			break LOOP
		case <-execute.Reload:
			log.Info().Msg("Reloading configuration")
		default:
			log.Debug().Msg("Supervisor loop signal")
		}
	}
	os.Exit(1)
}
