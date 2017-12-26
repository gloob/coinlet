package ticker

import (
	"time"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
)

type Ticker struct {
	logger log.Logger
	stop   chan chan struct{}
}

func NewTicker(logger log.Logger) *Ticker {
	return &Ticker{
		logger: logger,
		stop:   make(chan chan struct{}),
	}
}

// Run returns when Stop is invoked.
func (t *Ticker) Run() {
	level.Info(t.logger).Log("msg", "Starting ticker...")
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			level.Debug(t.logger).Log("msg", "tick")
		case q := <-t.stop:
			close(q)
			return
		}
	}
}

// Stop the ticker.
func (t *Ticker) Stop() {
	level.Info(t.logger).Log("msg", "Shutting down ticker...")
	q := make(chan struct{})
	t.stop <- q
	<-q
}
