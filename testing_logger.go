package log

import (
	"testing"

	"github.com/go-logfmt/logfmt"
)

type testingLogger struct {
	tb testing.TB
}

// NewTestingLogger returns a logger that encodes keyvals to tb.Log in
// logfmt format. It is meant to be used in tests.
func NewTestingLogger(tb testing.TB) Logger {
	return testingLogger{tb}
}

func (t testingLogger) Log(keyvals ...interface{}) error {
	t.tb.Helper()
	buf, err := logfmt.MarshalKeyvals(keyvals...)
	if err != nil {
		return err
	}
	t.tb.Log(string(buf))
	return nil
}
