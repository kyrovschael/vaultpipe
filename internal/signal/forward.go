// Package signal handles OS signal forwarding to child processes.
package signal

import (
	"os"
	"os/signal"
	"syscall"

	"golang.org/x/sys/unix"
)

// Forwarder listens for OS signals and forwards them to a target process.
type Forwarder struct {
	proc   *os.Process
	signals chan os.Signal
	done   chan struct{}
}

// New creates a Forwarder that will forward signals to proc.
func New(proc *os.Process) *Forwarder {
	return &Forwarder{
		proc:    proc,
		signals: make(chan os.Signal, 8),
		done:    make(chan struct{}),
	}
}

// Start begins listening for signals and forwarding them to the child process.
// It returns immediately; call Stop to clean up.
func (f *Forwarder) Start() {
	signal.Notify(f.signals,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGHUP,
		syscall.SIGUSR1,
		syscall.SIGUSR2,
	)
	go f.loop()
}

// Stop unregisters signal handlers and waits for the forwarding goroutine to exit.
func (f *Forwarder) Stop() {
	signal.Stop(f.signals)
	close(f.signals)
	<-f.done
}

func (f *Forwarder) loop() {
	defer close(f.done)
	for sig := range f.signals {
		if us, ok := sig.(unix.Signal); ok {
			_ = f.proc.Signal(us)
		} else {
			_ = f.proc.Signal(sig)
		}
	}
}
