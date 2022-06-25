package sanji

import (
	"errors"
	"strings"
	"sync"
	"sync/atomic"

	"github.com/bitleak/lmstfy/client"
)

var ManagerStartedErr = errors.New("Luffy manager has started")

type imanager interface {
	start()
	quitPrepare()
	quit()
	manage()
	loadWorkers()
	prcessingWorkerSize() int
	getQueues() string
	reset()
}

type HandlersChain []JobFunc

type manager struct {
	wait *sync.WaitGroup

	captain *Luffy
	fetch   ifetcher

	queue       string
	queues      []string
	queuesStr   string
	is_priority bool
	is_started  int32

	concurrency  int
	jobTTR       uint32
	job          JobFunc
	handlers     HandlersChain
	workers      []iworker
	workersMutex sync.Mutex

	confirm chan *client.Job
	stop    chan struct{}
	exit    chan struct{}
}

func (m *manager) insertMid(mids ...JobFunc) {
	m.handlers = append(m.handlers, mids...)
}

func (m *manager) start() {
	if m.started() {
		return
	}
	atomic.StoreInt32(&m.is_started, 1)
	m.wait.Add(1)
	m.insertMid(m.job)

	m.loadWorkers()
	go m.manage()
}

func (m *manager) started() bool {
	return atomic.LoadInt32(&m.is_started) > 0
}

func (m *manager) quitPrepare() {
	if !m.fetch.Closed() {
		m.fetch.Close()
	}
}

func (m *manager) quit() {
	Logger.Println("quitting queue", m.queue, ", waiting for ", m.prcessingWorkerSize(), "/", len(m.workers), " workers")
	m.quitPrepare()

	m.workersMutex.Lock()

	for _, worker := range m.workers {
		worker.quit()
	}

	m.workersMutex.Unlock()
	m.stop <- struct{}{}
	<-m.exit
	m.reset()
	m.wait.Done()
}

func (m *manager) manage() {
	Logger.Println("prcessing namespace", m.captain.config.Namespace, " queue", m.queue, " with ", m.concurrency, " workers.")
	for _, worker := range m.workers {
		worker.start()
	}
	go m.fetch.Fetch()

	for {
		select {
		case job := <-m.confirm:
			m.fetch.Acknowledge(job)
		case <-m.stop:
			m.exit <- struct{}{}
			break
		}
	}
}

func (m *manager) loadWorkers() {
	m.workersMutex.Lock()
	for i := 0; i < m.concurrency; i++ {
		m.workers[i] = newWorker(m)
	}
	m.workersMutex.Unlock()
}

func (m *manager) prcessingWorkerSize() int {
	count := 0
	for _, worker := range m.workers {
		if worker.processing() {
			count++
		}
	}
	return count
}

func (m *manager) getQueues() string {
	return m.queuesStr
}

func (m *manager) reset() {
	m.fetch = newFetcher(m, m.queuesStr, make(chan *client.Job), make(chan struct{}))
}

func newManager(captain *Luffy, queue string, job JobFunc, concurrency int, priority bool, ttr uint32) *manager {
	m := &manager{
		wait:        new(sync.WaitGroup),
		captain:     captain,
		queue:       queue,
		job:         job,
		jobTTR:      ttr,
		concurrency: concurrency,
		is_priority: priority,
		workers:     make([]iworker, concurrency),
		confirm:     make(chan *client.Job),
		stop:        make(chan struct{}),
		exit:        make(chan struct{}),
		handlers:    captain.handlers,
	}

	if m.is_priority {
		m.queues = []string{
			m.captain.priorityQueue(m.queue, PriorityLevelHight),
			m.captain.priorityQueue(m.queue, PriorityLevelMid),
			m.captain.priorityQueue(m.queue, PriorityLevelLow),
		}
		m.queuesStr = strings.Join(m.queues, ",")
	} else {
		qStr := m.captain.priorityQueue(m.queue, PriorityLevelHight)
		m.queues = []string{qStr}
		m.queuesStr = qStr
	}
	m.reset()
	return m
}
