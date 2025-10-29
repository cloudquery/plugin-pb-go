//go:build !windows

package managedplugin

import (
	"fmt"
	"os"
	"syscall"
	"time"
)

func getSysProcAttr() *syscall.SysProcAttr {
	return &syscall.SysProcAttr{
		// launch as new process group so that signals (ex: SIGINT) are not sent to the child process
		Setpgid: true, // UNIX systems
	}
}

func (c *Client) terminateProcess() error {
	c.logger.Debug().Msgf("sending interrupt signal to %s plugin", c.typ.String())
	if err := c.cmd.Process.Signal(os.Interrupt); err != nil {
		c.logger.Error().Err(err).Msgf("failed to send interrupt signal to %s plugin", c.typ.String())
	}
	timer := time.AfterFunc(5*time.Second, func() {
		c.logger.Info().Msgf("sending kill signal to %s plugin", c.typ.String())
		if err := c.cmd.Process.Kill(); err != nil {
			c.logger.Error().Err(err).Msgf("failed to kill %s plugin", c.typ.String())
		}
	})
	c.logger.Info().Msgf("waiting for %s plugin to terminate", c.typ.String())
	st, err := c.cmd.Process.Wait()
	timer.Stop()
	if err != nil {
		return err
	}
	if !st.Success() {
		var additionalInfo string
		status := st.Sys().(syscall.WaitStatus)
		if status.Signaled() && st.ExitCode() != -1 {
			additionalInfo += fmt.Sprintf(" (exit code: %d)", st.ExitCode())
		}
		if st.ExitCode() == 137 {
			additionalInfo = " (Out of Memory)"
		}
		return fmt.Errorf("%s plugin process failed with %s%s", c.typ.String(), st.String(), additionalInfo)
	}

	return nil
}
