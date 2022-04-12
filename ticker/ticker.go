package ticker

import (
	"fmt"
	"runtime/debug"
	"time"
)

type Runner func()

type Ticker struct {
	interval time.Duration
	ticker   *time.Ticker
	runner   Runner
}

func New(interval time.Duration, runner Runner) *Ticker {
	return &Ticker{
		interval: interval,
		ticker:   time.NewTicker(interval),
		runner:   runner,
	}
}

func (t *Ticker) Start() {
	go func() {
		defer func() {
			if e := recover(); e != nil {
				stack := string(debug.Stack())
				fmt.Println(stack)
				fmt.Println(e)
			}
		}()
		for range t.ticker.C {
			t.runner()
		}
	}()
}

func (t *Ticker) Stop() {
	t.ticker.Stop()
}

func (t *Ticker) Reset(d time.Duration) time.Duration {
	oldInterval := t.interval
	t.interval = d
	t.ticker.Reset(d)
	return oldInterval
}

func (t *Ticker) Interval() time.Duration {
	return t.interval
}
