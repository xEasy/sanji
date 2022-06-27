package sanji

import (
	"os"
	"os/signal"
	"syscall"
)

func (s *Luffy) handleSignals() {
	signals := make(chan os.Signal)
	signal.Notify(signals, syscall.SIGUSR1, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP, syscall.SIGQUIT)

	for sig := range signals {
		switch sig {
		case syscall.SIGUSR1, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP, syscall.SIGQUIT:
			s.Quit()
		}
	}
}
