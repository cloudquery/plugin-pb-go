//go:build !windows

package managedplugin

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"os/signal"
	"testing"
	"time"
)

const (
	// interruptedExitCode is the code the helper process exits with once it
	// receives os.Interrupt, chosen so it cannot be confused with a normal exit
	// or with the -1 reported for a process killed by a signal.
	interruptedExitCode = 42

	// helperReadyMarker is printed by the helper process once its signal
	// handler is installed, so the parent never signals it too early.
	helperReadyMarker = "cq-helper-ready"

	helperProcessEnv = "CQ_WANT_HELPER_PROCESS"
)

// TestHelperProcessInterrupt is not a real test. It is re-executed as the child
// process by the tests below (the standard os/exec helper-process pattern): it
// waits for os.Interrupt and exits with interruptedExitCode, so the parent can
// tell a graceful interrupt apart from a kill.
func TestHelperProcessInterrupt(t *testing.T) {
	if os.Getenv(helperProcessEnv) != "1" {
		t.Skip("not a real test; only run as a child process")
	}
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)
	fmt.Println(helperReadyMarker)
	<-ch
	os.Exit(interruptedExitCode)
}

// startInterruptibleChild starts TestHelperProcessInterrupt in a child process
// and returns once that child has installed its signal handler. configure is
// applied to the command before it is started.
func startInterruptibleChild(ctx context.Context, t *testing.T, configure func(*exec.Cmd)) *exec.Cmd {
	t.Helper()

	cmd := exec.CommandContext(ctx, os.Args[0], "-test.run=^TestHelperProcessInterrupt$", "-test.timeout=60s")
	cmd.Env = append(os.Environ(), helperProcessEnv+"=1")
	cmd.SysProcAttr = getSysProcAttr()
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		t.Fatalf("failed to get stdout pipe: %v", err)
	}
	if configure != nil {
		configure(cmd)
	}
	if err := cmd.Start(); err != nil {
		t.Fatalf("failed to start child: %v", err)
	}
	t.Cleanup(func() {
		if cmd.ProcessState == nil && cmd.Process != nil {
			_ = cmd.Process.Kill()
		}
	})

	ready := make(chan error, 1)
	go func() { ready <- awaitMarker(stdout) }()
	select {
	case err := <-ready:
		if err != nil {
			t.Fatalf("child never reported ready: %v", err)
		}
	case <-time.After(30 * time.Second):
		t.Fatal("timed out waiting for child to install its signal handler")
	}

	return cmd
}

func awaitMarker(r io.Reader) error {
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		if scanner.Text() == helperReadyMarker {
			return nil
		}
	}
	if err := scanner.Err(); err != nil {
		return err
	}
	return errors.New("child exited before printing the ready marker")
}

// waitForExit waits for the child and returns its exit code, failing the test
// if it does not exit promptly.
func waitForExit(t *testing.T, cmd *exec.Cmd) int {
	t.Helper()

	waited := make(chan error, 1)
	go func() { waited <- cmd.Wait() }()

	select {
	case err := <-waited:
		var exitErr *exec.ExitError
		if err != nil && !errors.As(err, &exitErr) && !errors.Is(err, context.Canceled) {
			t.Fatalf("unexpected Wait error: %v", err)
		}
	case <-time.After(30 * time.Second):
		_ = cmd.Process.Kill()
		t.Fatal("child did not exit within 30s")
	}
	return cmd.ProcessState.ExitCode()
}

// TestConfigureCmdCancel asserts that cancelling the context interrupts the
// plugin process gracefully rather than killing it, and that WaitDelay bounds
// how long we wait for it afterwards.
func TestConfigureCmdCancel(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cmd := startInterruptibleChild(ctx, t, configureCmdCancel)

	if cmd.WaitDelay != 30*time.Second {
		t.Errorf("WaitDelay = %v, want %v", cmd.WaitDelay, 30*time.Second)
	}
	if cmd.Cancel == nil {
		t.Fatal("Cancel was not set")
	}

	cancel()

	if got := waitForExit(t, cmd); got != interruptedExitCode {
		t.Errorf("child exit code = %d, want %d (child was killed rather than interrupted)", got, interruptedExitCode)
	}
}

// TestConfigureCmdCancelIsRequired documents the default behaviour the fix
// replaces: without configureCmdCancel, exec kills the plugin on cancellation.
func TestConfigureCmdCancelIsRequired(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cmd := startInterruptibleChild(ctx, t, nil)

	cancel()

	if got := waitForExit(t, cmd); got == interruptedExitCode {
		t.Errorf("child exit code = %d, want anything else: the default cancel should kill, not interrupt", got)
	}
}

// TestInterruptProcess asserts the platform helper delivers a graceful
// interrupt on this OS.
func TestInterruptProcess(t *testing.T) {
	cmd := startInterruptibleChild(context.Background(), t, nil)

	if err := interruptProcess(cmd.Process); err != nil {
		t.Fatalf("interruptProcess: %v", err)
	}

	if got := waitForExit(t, cmd); got != interruptedExitCode {
		t.Errorf("child exit code = %d, want %d", got, interruptedExitCode)
	}
}
