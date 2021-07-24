package metrics

import "time"

// Recorder knows how to measure different kind of metrics.
type Recorder interface {
	// WithID will set the ID name to the recorde and every metric
	// measured with the obtained recorder will be identified with
	// the name.
	WithID(id string) Recorder
	// IncAttempt will increment the number of attempts.
	IncAttempt(target string)
	// IncRetry will increment the number of retries.
	IncRetry(target string)
	// IncTimeout will increment the number of timeouts.
	IncTimeout(target string)
	// IncBulkheadQueued increments the number of queued Funcs to execute.
	IncBulkheadQueued(target string)
	// IncBulkheadProcessed increments the number of processed Funcs to execute.
	IncBulkheadProcessed(target string)
	// IncBulkheadProcessed increments the number of timeouts Funcs waiting  to execute.
	IncBulkheadTimeout(target string)
	// IncCircuitbreakerState increments the number of state change.
	IncCircuitbreakerState(target string, state string)
	// IncChaosInjectedFailure increments the number of times injected failure.
	IncChaosInjectedFailure(target string, kind string)
	// SetConcurrencyLimitInflightExecutions sets the number of queued and executions at a given moment.
	SetConcurrencyLimitInflightExecutions(target string, q int)
	// SetConcurrencyLimitExecutingExecutions sets the number of executions at a given moment.
	SetConcurrencyLimitExecutingExecutions(target string, q int)
	// IncConcurrencyLimitResult increments the results obtained by the executions after applying the
	// limiter result policy.
	IncConcurrencyLimitResult(target string, result string)
	// SetConcurrencyLimitLimiterLimit sets the current limit the limiter algorithm has calculated.
	SetConcurrencyLimitLimiterLimit(target string, limit int)
	// ObserveConcurrencyLimitQueuedTime will measure the duration of a function waiting on a queue until it's executed.
	ObserveConcurrencyLimitQueuedTime(target string, start time.Time)
}
