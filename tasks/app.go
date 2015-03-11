package tasks

import (
	"bytes"
	"fmt"
	. "github.com/tbud/bud/context"
	"io"
	"os"
	"os/exec"
	"time"
)

type App struct {
	*exec.Cmd
}

func NewApp(binPath string, port int) *App {
	cmd := exec.Command(binPath,
		fmt.Sprintf("-port=%d", port))
	cmd.Stdout, cmd.Stderr = os.Stdout, os.Stderr
	return &App{cmd}
}

func (a *App) Start() error {
	listeningWriter := startupListeningWriter{os.Stdout, make(chan bool)}
	a.Stdout = listeningWriter
	Log.Trace("Exec app: %s, %v", a.Path, a.Args)
	if err := a.Cmd.Start(); err != nil {
		Log.Error("Error running: %v", err)
		return err
	}

	select {
	case <-a.waitChan():
		return fmt.Errorf("app died")
	case <-time.After(30 * time.Second):
		a.Kill()
		return fmt.Errorf("app timed out")
	case <-listeningWriter.notifyReady:
		return nil
	}
}

func (a *App) Run() {
	Log.Trace("Exec app:", a.Path, a.Args)
	if err := a.Cmd.Run(); err != nil {
		Log.Error("Error running: %v", err)
	}
}

func (a *App) Kill() {
	if a.Cmd != nil && (a.ProcessState == nil || !a.ProcessState.Exited()) {
		Log.Trace("Killing server pid %d", a.Process.Pid)
		err := a.Process.Kill()
		if err != nil {
			Log.Error("Failed to kill server: %v", err)
		}
	}
}

func (a *App) waitChan() <-chan int {
	ch := make(chan int)
	go func() {
		a.Wait()
		ch <- 1
	}()
	return ch
}

type startupListeningWriter struct {
	dest        io.Writer
	notifyReady chan bool
}

func (s startupListeningWriter) Write(p []byte) (n int, err error) {
	if s.notifyReady != nil && bytes.Contains(p, []byte("Listening")) {
		s.notifyReady <- true
		s.notifyReady = nil
	}
	return s.dest.Write(p)
}
