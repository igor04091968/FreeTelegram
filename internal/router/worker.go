package router

import (
	"context"
	"time"

	"freetelegram/internal/telemetry"
)

type Worker struct {
	updater *Updater
	stats   *telemetry.Stats
	period  time.Duration
}

func NewWorker(updater *Updater, stats *telemetry.Stats, period time.Duration) *Worker {
	if period <= 0 {
		period = 5 * time.Minute
	}
	return &Worker{updater: updater, stats: stats, period: period}
}

func (w *Worker) Run(ctx context.Context) {
	ticker := time.NewTicker(w.period)
	defer ticker.Stop()

	for {
		_, err := w.updater.RunOnce(ctx)
		if err != nil {
			w.stats.MarkError(err.Error())
		} else {
			w.stats.MarkUpdate(0)
		}
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
		}
	}
}
