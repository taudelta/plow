package plow

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"syscall"
	"time"
)

var logger *log.Logger

func init() {
	logger = log.New(io.Discard, "run cmd", log.LstdFlags)
}

func EnableScriptLogger() {
	logger.SetOutput(os.Stdout)
}

func RunCmd(timeout time.Duration, buf *bytes.Buffer, args ...string) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	var err error
	if runtime.GOOS == "windows" {
		err = execCmd(ctx, "powershell.exe", buf, args...)
	} else {
		args = append([]string{"-c"}, args...)
		err = execCmd(ctx, "bash", buf, args...)
	}
	if err != nil {
		log.Fatal(err)
	}
}

func RunScript(scriptPath string, timeout time.Duration, buf *bytes.Buffer) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	scriptPath = strings.TrimRight(scriptPath, "\n")
	var err error
	if runtime.GOOS == "windows" {
		scriptArg := fmt.Sprintf(`"& ""%s"""`, strings.TrimRight(scriptPath, "\n"))
		err = execCmd(ctx, "powershell.exe", buf, "-NonInteractive", "-noexit", scriptArg)
	} else {
		err = execCmd(ctx, scriptPath, buf)
	}
	if err != nil {
		log.Fatal(err)
	}
}

func execCmd(ctx context.Context, shell string, buf *bytes.Buffer, args ...string) error {
	cmd := exec.CommandContext(ctx, shell, args...)
	cmd.Env = os.Environ()

	errBuf := bytes.NewBuffer([]byte{})

	cmd.Stdout = buf
	cmd.Stderr = errBuf

	if err := cmd.Start(); err != nil {
		return err
	}

	done := make(chan error, 1)
	go func() {
		done <- cmd.Wait()
	}()

	select {
	case <-ctx.Done():
		if runtime.GOOS == "windows" {
			return cmd.Process.Kill()
		} else if runtime.GOOS == "linux" || runtime.GOOS == "darwin" {
			if p, err := os.FindProcess(-os.Getpid()); err != nil {
				return err
			} else {
				return p.Signal(syscall.SIGTERM)
			}
		}
	case err := <-done:
		if err != nil {
			return fmt.Errorf("process finished with error = %w", err)
		}
		if errBuf.String() != "" {
			logger.Println("err", errBuf.String())
		}
		logger.Print("Process finished successfully")
		cmd.Process.Release()
	}

	return nil
}
