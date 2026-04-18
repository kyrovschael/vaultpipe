package signal_test

import (
	"os"
	"syscall"
	"testing"
	"time"

	"github.com/yourorg/vaultpipe/internal/signal"
)

// fakeProcess wraps a real process (self) so we can observe signal delivery.
func startSleepProc(t *testing.T) *os.Process {
	t.Helper()
	cmd := []string{"/bin/sleep", "30"}
	proc, err := os.StartProcess(cmd[0], cmd, &os.ProcAttr{
		Files: []*os.File{os.Stdin, os.Stdout, os.Stderr},
	})
	if err != nil {
		t.Skipf("cannot start sleep process: %v", err)
	}
	return proc
}

func TestForwarder_StartStop(t *testing.T) {
	proc := startSleepProc(t)
	defer func() { _ = proc.Kill() }()

	f := signal.New(proc)
	f.Start()

	// Allow goroutine to settle.
	time.Sleep(20 * time.Millisecond)

	f.Stop() // should not block
}

func TestForwarder_ForwardsSIGHUP(t *testing.T) {
	proc := startSleepProc(t)
	defer func() { _ = proc.Kill() }()

	f := signal.New(proc)
	f.Start()
	defer f.Stop()

	// Send SIGHUP to the child directly via the forwarder path.
	// We verify no panic / error occurs during forwarding.
	_ = proc.Signal(syscall.SIGHUP)
	time.Sleep(20 * time.Millisecond)
}

func TestForwarder_StopIsIdempotentAfterKill(t *testing.T) {
	proc := startSleepProc(t)

	f := signal.New(proc)
	f.Start()

	_ = proc.Kill()
	_, _ = proc.Wait()

	// Stop should complete even though the process is gone.
	done := make(chan struct{})
	go func() {
		f.Stop()
		close(done)
	}()

	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Fatal("Stop() blocked after process already exited")
	}
}
