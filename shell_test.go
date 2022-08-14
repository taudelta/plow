package plow

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestRunCmd(t *testing.T) {
	buf := bytes.NewBuffer([]byte{})
	RunCmd(10*time.Second, buf, "echo test")
	result := strings.TrimRight(buf.String(), "\r\n")
	if result != "test" {
		t.Errorf("%q != test", result)
	}
}

func TestRunScript(t *testing.T) {
	buf := bytes.NewBuffer([]byte{})
	curDir, _ := os.Getwd()
	RunScript(filepath.Join(curDir, "test.ps1"), 10*time.Second, buf)
	result := strings.TrimRight(buf.String(), "\r\n")
	if result == "" {
		t.Error("result is empty")
	}
}
