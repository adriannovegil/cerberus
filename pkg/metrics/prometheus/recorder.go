package prometheus

import (
	"fmt"
	"time"

	cPrometheus "github.com/prometheus/client_golang/prometheus"

	"devcircus.com/cerberus/pkg/metrics"
)

var (
	promNamespace = "cerberus"

	promCommandSubsystem          = "command"
	promRetrySubsystem            = "retry"
	promTimeoutSubsystem          = "timeout"
	promBulkheadSubsystem         = "bulkhead"
	promCBSubsystem               = "circuitbreaker"
	promChaosSubsystem            = "chaos"
	promConcurrencyLimitSubsystem = "concurrencylimit"
)

type recorder struct {
	// Metrics.
	cmdExecutionDuration           *cPrometheus.HistogramVec
	retryRetries                   *cPrometheus.CounterVec
	timeoutTimeouts                *cPrometheus.CounterVec
	bulkQueued                     *cPrometheus.CounterVec
	bulkProcessed                  *cPrometheus.CounterVec
	bulkTimeouts                   *cPrometheus.CounterVec
	cbStateChanges                 *cPrometheus.CounterVec
	chaosFailureInjections         *cPrometheus.CounterVec
	concurrencyLimitInflights      *cPrometheus.GaugeVec
	concurrencyLimitExecuting      *cPrometheus.GaugeVec
	concurrencyLimitResult         *cPrometheus.CounterVec
	concurrencyLimitLimit          *cPrometheus.GaugeVec
	concurrencyLimitQueuedDuration *cPrometheus.HistogramVec

	id  string
	reg cPrometheus.Registerer
}

// NewRecorder returns a new Recorder that knows how to measure
// using Prometheus kind metrics.
func NewRecorder(reg cPrometheus.Registerer) metrics.Recorder {
	r := &recorder{
		reg: reg,
	}

	r.registerMetrics()
	return r
}

func (r recorder) WithID(id string) metrics.Recorder {
	return &recorder{
		cmdExecutionDuration:           r.cmdExecutionDuration,
		retryRetries:                   r.retryRetries,
		timeoutTimeouts:                r.timeoutTimeouts,
		bulkQueued:                     r.bulkQueued,
		bulkProcessed:                  r.bulkProcessed,
		bulkTimeouts:                   r.bulkTimeouts,
		cbStateChanges:                 r.cbStateChanges,
		chaosFailureInjections:         r.chaosFailureInjections,
		concurrencyLimitInflights:      r.concurrencyLimitInflights,
		concurrencyLimitExecuting:      r.concurrencyLimitExecuting,
		concurrencyLimitResult:         r.concurrencyLimitResult,
		concurrencyLimitLimit:          r.concurrencyLimitLimit,
		concurrencyLimitQueuedDuration: r.concurrencyLimitQueuedDuration,

		id:  id,
		reg: r.reg,
	}
}

func (r *recorder) registerMetrics() {
	r.cmdExecutionDuration = cPrometheus.NewHistogramVec(cPrometheus.HistogramOpts{
		Namespace: promNamespace,
		Subsystem: promCommandSubsystem,
		Name:      "execution_duration_seconds",
		Help:      "The duration of the command execution in seconds.",
		Buckets:   cPrometheus.DefBuckets,
	}, []string{"id", "success"})

	r.retryRetries = cPrometheus.NewCounterVec(cPrometheus.CounterOpts{
		Namespace: promNamespace,
		Subsystem: promRetrySubsystem,
		Name:      "retries_total",
		Help:      "Total number of retries made by the retry runner.",
	}, []string{"id", "worker"})

	r.timeoutTimeouts = cPrometheus.NewCounterVec(cPrometheus.CounterOpts{
		Namespace: promNamespace,
		Subsystem: promTimeoutSubsystem,
		Name:      "timeouts_total",
		Help:      "Total number of timeouts made by the timeout runner.",
	}, []string{"id"})

	r.bulkQueued = cPrometheus.NewCounterVec(cPrometheus.CounterOpts{
		Namespace: promNamespace,
		Subsystem: promBulkheadSubsystem,
		Name:      "queued_total",
		Help:      "Total number of queued funcs made by the bulkhead runner.",
	}, []string{"id"})

	r.bulkProcessed = cPrometheus.NewCounterVec(cPrometheus.CounterOpts{
		Namespace: promNamespace,
		Subsystem: promBulkheadSubsystem,
		Name:      "processed_total",
		Help:      "Total number of processed funcs made by the bulkhead runner.",
	}, []string{"id"})

	r.bulkTimeouts = cPrometheus.NewCounterVec(cPrometheus.CounterOpts{
		Namespace: promNamespace,
		Subsystem: promBulkheadSubsystem,
		Name:      "timeouts_total",
		Help:      "Total number of timeouts funcs waiting for execution made by the bulkhead runner.",
	}, []string{"id"})

	r.cbStateChanges = cPrometheus.NewCounterVec(cPrometheus.CounterOpts{
		Namespace: promNamespace,
		Subsystem: promCBSubsystem,
		Name:      "state_changes_total",
		Help:      "Total number of state changes made by the circuit breaker runner.",
	}, []string{"id", "state"})

	r.chaosFailureInjections = cPrometheus.NewCounterVec(cPrometheus.CounterOpts{
		Namespace: promNamespace,
		Subsystem: promChaosSubsystem,
		Name:      "failure_injections_total",
		Help:      "Total number of failure injectionsmade by the chaos runner.",
	}, []string{"id", "kind"})

	r.concurrencyLimitInflights = cPrometheus.NewGaugeVec(cPrometheus.GaugeOpts{
		Namespace: promNamespace,
		Subsystem: promConcurrencyLimitSubsystem,
		Name:      "inflight_executions",
		Help:      "The number of inflight executions, these are executing and queued.",
	}, []string{"id"})

	r.concurrencyLimitExecuting = cPrometheus.NewGaugeVec(cPrometheus.GaugeOpts{
		Namespace: promNamespace,
		Subsystem: promConcurrencyLimitSubsystem,
		Name:      "executing_executions",
		Help:      "The number of executing executions.",
	}, []string{"id"})

	r.concurrencyLimitResult = cPrometheus.NewCounterVec(cPrometheus.CounterOpts{
		Namespace: promNamespace,
		Subsystem: promConcurrencyLimitSubsystem,
		Name:      "result_total",
		Help:      "Total results of the executions measured by the limiter algorithm.",
	}, []string{"id", "result"})

	r.concurrencyLimitLimit = cPrometheus.NewGaugeVec(cPrometheus.GaugeOpts{
		Namespace: promNamespace,
		Subsystem: promConcurrencyLimitSubsystem,
		Name:      "limiter_limit",
		Help:      "The concurrency limit measured and calculated by the limiter algorithm.",
	}, []string{"id"})

	r.concurrencyLimitQueuedDuration = cPrometheus.NewHistogramVec(cPrometheus.HistogramOpts{
		Namespace: promNamespace,
		Subsystem: promConcurrencyLimitSubsystem,
		Name:      "queued_duration_seconds",
		Help:      "The duration of the command waiting on the queue.",
		Buckets:   []float64{.001, .005, .01, .015, .025, 0.05, 0.1, 0.2, 0.5, 1, 2.5, 5, 10},
	}, []string{"id"})

	r.reg.MustRegister(r.cmdExecutionDuration,
		r.retryRetries,
		r.timeoutTimeouts,
		r.bulkQueued,
		r.bulkProcessed,
		r.bulkTimeouts,
		r.cbStateChanges,
		r.chaosFailureInjections,
		r.concurrencyLimitInflights,
		r.concurrencyLimitExecuting,
		r.concurrencyLimitResult,
		r.concurrencyLimitLimit,
		r.concurrencyLimitQueuedDuration,
	)
}

func (r recorder) ObserveCommandExecution(start time.Time, success bool) {
	secs := time.Since(start).Seconds()
	r.cmdExecutionDuration.WithLabelValues(r.id, fmt.Sprintf("%t", success)).Observe(secs)
}

func (r recorder) IncRetry(worker string) {
	r.retryRetries.WithLabelValues(r.id, worker).Inc()
}

func (r recorder) IncTimeout() {
	r.timeoutTimeouts.WithLabelValues(r.id).Inc()
}

func (r recorder) IncBulkheadQueued() {
	r.bulkQueued.WithLabelValues(r.id).Inc()
}

func (r recorder) IncBulkheadProcessed() {
	r.bulkProcessed.WithLabelValues(r.id).Inc()
}

func (r recorder) IncBulkheadTimeout() {
	r.bulkTimeouts.WithLabelValues(r.id).Inc()
}

func (r recorder) IncCircuitbreakerState(state string) {
	r.cbStateChanges.WithLabelValues(r.id, state).Inc()
}

func (r recorder) IncChaosInjectedFailure(kind string) {
	r.chaosFailureInjections.WithLabelValues(r.id, kind).Inc()
}

func (r recorder) SetConcurrencyLimitInflightExecutions(q int) {
	r.concurrencyLimitInflights.WithLabelValues(r.id).Set(float64(q))
}

func (r recorder) SetConcurrencyLimitExecutingExecutions(q int) {
	r.concurrencyLimitExecuting.WithLabelValues(r.id).Set(float64(q))
}

func (r recorder) IncConcurrencyLimitResult(result string) {
	r.concurrencyLimitResult.WithLabelValues(r.id, result).Inc()
}

func (r recorder) SetConcurrencyLimitLimiterLimit(limit int) {
	r.concurrencyLimitLimit.WithLabelValues(r.id).Set(float64(limit))
}

func (r recorder) ObserveConcurrencyLimitQueuedTime(start time.Time) {
	secs := time.Since(start).Seconds()
	r.concurrencyLimitQueuedDuration.WithLabelValues(r.id).Observe(secs)
}
