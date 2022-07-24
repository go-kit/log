package log_test

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/go-kit/log"
)

// BenchmarkLineBuffered creates line-buffered loggers of various capacities to see which perform best.
func BenchmarkLineBuffered(b *testing.B) {
	b.ReportAllocs()

	for i := uint32(1); i <= 1024; i *= 2 {
		b.Run(fmt.Sprintf("capacity:%d", i), func(b *testing.B) {
			f := outFile(b)
			defer os.Remove(f.Name())

			bufLog := log.NewLineBufferedLogger(f, i, 10*time.Millisecond)
			l := log.NewLogfmtLogger(log.NewSyncWriter(bufLog))

			benchmarkRunner(b, l, baseMessage)

			// force a final flush for outstanding lines in buffer
			bufLog.Flush()

			b.StopTimer()
			contents, err := ioutil.ReadFile(f.Name())
			if err != nil {
				b.Errorf("could not read test file: %s", err)
			}
			lines := strings.Split(string(contents), "\n")
			b.StartTimer()

			if want, have := b.N, len(lines)-1; want != have {
				b.Errorf("expected %d lines, have %d", want, have)
			}
		})
	}
}

// BenchmarkLineUnbuffered should perform roughly equivalently to a line-buffered logger with a capacity of 1.
func BenchmarkLineUnbuffered(b *testing.B) {
	b.ReportAllocs()

	f := outFile(b)
	defer os.Remove(f.Name())

	l := log.NewLogfmtLogger(log.NewSyncWriter(f))
	benchmarkRunner(b, l, baseMessage)

	b.StopTimer()
	contents, err := ioutil.ReadFile(f.Name())
	if err != nil {
		b.Errorf("could not read test file: %s", err)
	}
	lines := strings.Split(string(contents), "\n")
	b.StartTimer()

	if want, have := b.N, len(lines)-1; want != have {
		b.Errorf("expected %d lines, have %d", want, have)
	}
}

// outFile creates a real OS file for testing.
// We cannot use stdout/stderr since we need to read the contents after to validate, and we have to write to a file
// to benchmark the impact of write() syscalls.
func outFile(b *testing.B) *os.File {
	f, err := os.OpenFile("/tmp/test", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		b.Fatalf("cannot create test file: %s", err)
	}

	return f
}
