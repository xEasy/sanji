package sanji

import (
	"time"

	"github.com/bitleak/lmstfy/client"
)

const (
	defaultConsumeTimeout = 5
)

type ifetcher interface {
	Queue() string
	Fetch()
	Acknowledge(*client.Job)
	Ready() chan struct{}
	FinishedWork() chan struct{}
	Messages() chan *client.Job
	Close()
	Closed() bool
}

type fetcher struct {
	manager      *manager
	queue        string
	queues       string
	messages     chan *client.Job
	ready        chan struct{}
	finishedwork chan struct{}
	stop         chan struct{}
	exit         chan struct{}
	closed       chan struct{}
}

func newFetcher(manager *manager, queue string, messages chan *client.Job, ready chan struct{}) ifetcher {
	return &fetcher{
		manager:      manager,
		queue:        queue,
		messages:     messages,
		ready:        ready,
		finishedwork: make(chan struct{}),
		stop:         make(chan struct{}),
		exit:         make(chan struct{}),
		closed:       make(chan struct{}),
	}
}

func (f *fetcher) Queue() string {
	return f.queue
}

func (f *fetcher) Fetch() {
	go func() {
		for {
			if f.Closed() {
				break
			}
			<-f.Ready()
			f.tryFetchMessage()
		}
	}()

	for {
		select {
		case <-f.stop:
			close(f.closed)
			close(f.exit)
			break
		}
	}
}

func (f *fetcher) Acknowledge(job *client.Job) {
	err := f.getClient().Ack(job.Queue, job.ID)
	if err != nil {
		Logger.Println("ack job err, queue:", job.Queue, " id:, ", job.ID, " err: ", err)
	}
}

func (f *fetcher) Ready() chan struct{} {
	return f.ready
}

func (f *fetcher) FinishedWork() chan struct{} {
	return f.finishedwork
}

func (f *fetcher) Messages() chan *client.Job {
	return f.messages
}

func (f *fetcher) Close() {
	f.stop <- struct{}{}
	<-f.exit
}

func (f *fetcher) Closed() bool {
	select {
	case <-f.closed:
		return true
	default:
		return false
	}
}

func (f *fetcher) getClient() *client.LmstfyClient {
	return f.manager.captain.client
}

func (f *fetcher) tryFetchMessage() {
	job, err := f.getClient().ConsumeFromQueues(f.manager.jobTTR, defaultConsumeTimeout, f.manager.queues...)
	if err != nil {
		Logger.Println("fetching message ERR: ", err)
		time.Sleep(time.Second * time.Duration(f.manager.captain.config.PollInterval))
	} else {
		if job != nil {
			f.sendJob(job)
		}

		// 任务队列为空，重新获取任务
		if job == nil {
			time.Sleep(time.Second * time.Duration(f.manager.captain.config.PollInterval))
			f.tryFetchMessage()
		}
	}
}

func (f *fetcher) sendJob(job *client.Job) {
	f.messages <- job
}
