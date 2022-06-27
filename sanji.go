package sanji

import (
	"log"
	"os"
	"sync"

	"github.com/bitleak/lmstfy/client"
)

var Logger *log.Logger = log.New(os.Stdout, "workers: ", log.Ldate|log.Lmicroseconds)

type Luffy struct {
	client *client.LmstfyClient
	config *Config

	started  bool
	handlers HandlersChain
	managers map[string]*manager
	access   sync.Mutex
}

type PriorityLevel int

const (
	PriorityLevelHight = iota
	PriorityLevelMid
	PriorityLevelLow
)

func New(token string, namespace string, host string, port int) *Luffy {
	return NewWithConfig(token, namespace, host, port, 10)
}

func NewWithConfig(token string, namespace string, host string, port int, poll_interval int) *Luffy {
	cf := Configure(&Config{
		Namespace:    namespace,
		Token:        token,
		Host:         host,
		Port:         port,
		PollInterval: poll_interval,
	})
	luffy := &Luffy{
		config:   cf,
		managers: make(map[string]*manager),
		handlers: HandlersChain{logger_middleware, recovery_middleware},
	}
	luffy.client = cf.client
	return luffy
}

func (s *Luffy) Use(mids ...JobFunc) {
	s.handlers = append(s.handlers, mids...)
}

func (s *Luffy) Process(queue string, job JobFunc, ttr uint32, concurrency int, is_priority bool) imanager {
	manager := newManager(s, queue, job, concurrency, is_priority, ttr)
	s.managers[queue] = manager
	return manager
}

func (s *Luffy) Run() {
	s.start()
	s.waitForExit()
}

func (s *Luffy) start() {
	s.access.Lock()
	defer s.access.Unlock()

	if s.started {
		return
	}
	s.startManagers()
	s.started = true
}

func (s *Luffy) startManagers() {
	for _, m := range s.managers {
		m.start()
	}
}

func (s *Luffy) Quit() {
	s.access.Lock()
	defer s.access.Unlock()

	if !s.started {
		return
	}
	s.quitManagers()
	s.waitForExit()
	s.started = false
}

func (s *Luffy) quitManagers() {
	for _, m := range s.managers {
		go func(m *manager) {
			m.quit()
		}(m)
	}
}

func (s *Luffy) waitForExit() {
	for _, m := range s.managers {
		m.wait.Wait()
	}
}
