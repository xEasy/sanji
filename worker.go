package sanji

import (
	"sync/atomic"
	"time"

	"github.com/bitleak/lmstfy/client"
)

type iworker interface {
	start()
	quit()
	work(jobs chan *client.Job)
	process(job *client.Job) bool
	processing() bool
}

type worker struct {
	manager    *manager
	stop       chan struct{}
	exit       chan struct{}
	currentJob *client.Job
	startedAt  int64
}

func (w *worker) start() {
	go w.work(w.manager.fetch.Messages())
}

func (w *worker) quit() {
	w.stop <- struct{}{}
	<-w.exit
}

func (w *worker) work(jobs chan *client.Job) {
	for {
		select {
		case job := <-jobs:
			atomic.StoreInt64(&w.startedAt, time.Now().UTC().Unix())
			w.currentJob = job

			if w.process(job) {
				w.manager.confirm <- job
			}
			atomic.StoreInt64(&w.startedAt, time.Now().UTC().Unix())
			w.currentJob = nil

			select {
			case w.manager.fetch.FinishedWork() <- struct{}{}:
			default:
			}
		case w.manager.fetch.Ready() <- struct{}{}:
		case <-w.stop:
			w.exit <- struct{}{}
			return
		}
	}
}

func (w *worker) process(job *client.Job) bool {
	ctx := &Context{
		Queue:       w.manager.queue,
		Job:         job,
		acknowledge: true,
		handlers:    w.manager.handlers,
	}

	defer func() {
		recover()
	}()

	ctx.Next()
	return ctx.acknowledge
}

func (w *worker) processing() bool {
	return atomic.LoadInt64(&w.startedAt) > 0
}

func newWorker(m *manager) iworker {
	return &worker{m, make(chan struct{}), make(chan struct{}), nil, 0}
}
