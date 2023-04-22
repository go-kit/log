//go:build windows
// +build windows

package eventlog_test

import (
	"fmt"

	"github.com/go-kit/log"
	"github.com/go-kit/log/eventlog"
	"github.com/go-kit/log/level"
	goeventlog "golang.org/x/sys/windows/svc/eventlog"
)

func ExampleNewEventLogLogger_defaultPrioritySelector() {
	// Normal eventlog writer
	w, err := goeventlog.Open("go-kit-log")
	if err != nil {
		fmt.Println(err)
		return
	}

	logger := eventlog.NewEventLogLogger(w, log.NewLogfmtLogger)

	type Task struct {
		ID int
	}

	RunTask := func(task Task, logger log.Logger) {
		logger.Log("taskID", task.ID, "event", "starting task")

		logger.Log("taskID", task.ID, level.Key(), level.DebugValue(), "msg", "debug because of explicit level")
		logger.Log("taskID", task.ID, level.Key(), level.WarnValue(), "msg", "warn because of explicit level")
		logger.Log("taskID", task.ID, level.Key(), level.ErrorValue(), "msg", "error because of explicit level")

		logger.Log("taskID", task.ID, "event", "task complete")
	}

	RunTask(Task{ID: 1}, logger)

	// Output:
}
