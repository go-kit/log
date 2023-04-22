//go:build windows
// +build windows

// Package eventlog provides a Logger that writes to Windows Event Log.
package eventlog

import (
	"testing"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	goeventlog "golang.org/x/sys/windows/svc/eventlog"
)

func TestLog(t *testing.T) {
	const name = "go-kit-log"

	w, err := goeventlog.Open(name)
	if err != nil {
		t.Fatalf("Open failed: %s", err)
	}
	defer w.Close()

	logger := NewEventLogLogger(w, log.NewLogfmtLogger)

	err = level.Debug(logger).Log("msg", "debug")
	if err != nil {
		t.Fatalf("debug failed: %s", err)
	}
	err = level.Info(logger).Log("msg", "info")
	if err != nil {
		t.Fatalf("Info failed: %s", err)
	}
	err = level.Warn(logger).Log("msg", "warn")
	if err != nil {
		t.Fatalf("Warning failed: %s", err)
	}
	err = level.Error(logger).Log("msg", "err")
	if err != nil {
		t.Fatalf("Error failed: %s", err)
	}

	err = level.Error(logger).Log("msg", "err")
}
