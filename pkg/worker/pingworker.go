package worker

import (
	"context"

	"devcircus.com/cerberus/pkg/target/ping"
)

// PingWorker execution data
type PingWorker struct {
	ID string
}

// NewPingWorker create a new instance worker
func NewPingWorker(id string, data ping.Config) *PingWorker {
	pw := PingWorker{}
	return &pw
}

// Start the job
func (pw *PingWorker) Start(ctx context.Context) {
	go pw.doWork(ctx)
}

func (pw *PingWorker) doWork(ctx context.Context) {
}

// StartTicker writes data to request channel for every x seconds
func (pw *PingWorker) StartTicker() {
}

// StartTimeRecorder start the time recorder channel
func (pw *PingWorker) StartTimeRecorder() {
}
