package worker

import (
	"context"
	"os"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/segmentio/ksuid"

	"devcircus.com/cerberus/pkg/config"
	"devcircus.com/cerberus/pkg/execute"
	"devcircus.com/cerberus/pkg/metrics"
	"devcircus.com/cerberus/pkg/metrics/prometheus"
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

	s.MetricsRecorder = prometheus.NewRecorder(
		prometheus.StartPrometheusServer()).WithID(s.genKsuid().String())

	//defer func(start time.Time) {
	//	s.MetricsRecorder.ObserveCommandExecution(start, true)
	//}(time.Now())

	ctx := metrics.SetRecorderOnContext(context.TODO(), s.MetricsRecorder)

	data := config.Config.Targets.Requests
	for i, requestConfig := range data {
		log.Debug().Msgf("Launching worker #%d: %s %s", i, requestConfig.RequestType, requestConfig.URL)

		w := NewWorker(s.genKsuid().String(), requestConfig)
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

func (s *Supervisor) genKsuid() ksuid.KSUID {
	return ksuid.New()
}
