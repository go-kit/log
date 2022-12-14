package log_test

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/go-kit/log"
)

const (
	flushPeriod = 10 * time.Millisecond
	bufferSize  = 10e6
)

// BenchmarkLineBuffered creates line-buffered loggers of various capacities to see which perform best.
func BenchmarkLineBuffered(b *testing.B) {

	for i := 1; i <= 2048; i *= 2 {
		f := outFile(b)
		defer os.Remove(f.Name())

		bufLog := log.NewLineBufferedLogger(f, uint32(i),
			log.WithFlushPeriod(flushPeriod),
			log.WithPrellocatedBuffer(bufferSize),
		)
		l := log.NewLogfmtLogger(bufLog)

		b.Run(fmt.Sprintf("capacity:%d", i), func(b *testing.B) {
			b.ReportAllocs()
			b.StartTimer()

			f.Truncate(0)

			logger := log.With(l, "common_key", "common_value")
			for j := 0; j < b.N; j++ {
				logger.Log("foo_key", "foo_value")
			}

			// force a final flush for outstanding lines in buffer
			bufLog.Flush()
			b.StopTimer()

			contents, err := os.ReadFile(f.Name())
			if err != nil {
				b.Errorf("could not read test file: %s", err)
			}
			lines := strings.Split(string(contents), "\n")

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

	l := log.NewLogfmtLogger(f)
	benchmarkRunner(b, l, baseMessage)

	b.StopTimer()

	contents, err := os.ReadFile(f.Name())
	if err != nil {
		b.Errorf("could not read test file: %s", err)
	}
	lines := strings.Split(string(contents), "\n")

	if want, have := b.N, len(lines)-1; want != have {
		b.Errorf("expected %d lines, have %d", want, have)
	}
}

func BenchmarkLineDiscard(b *testing.B) {
	b.ReportAllocs()

	l := log.NewLogfmtLogger(io.Discard)
	benchmarkRunner(b, l, baseMessage)
}

func TestLineBufferedConcurrency(t *testing.T) {
	t.Parallel()
	bufLog := log.NewLineBufferedLogger(io.Discard, 32,
		log.WithFlushPeriod(flushPeriod),
		log.WithPrellocatedBuffer(bufferSize),
	)
	testConcurrency(t, log.NewLogfmtLogger(bufLog), 10000)
}

func TestOnFlushCallback(t *testing.T) {
	var (
		flushCount     uint32
		flushedEntries int
		buf            bytes.Buffer
	)

	callback := func(entries uint32) {
		flushCount++
		flushedEntries += int(entries)
	}

	bufLog := log.NewLineBufferedLogger(&buf, 2,
		log.WithFlushPeriod(flushPeriod),
		log.WithPrellocatedBuffer(bufferSize),
		log.WithFlushCallback(callback),
	)

	l := log.NewLogfmtLogger(bufLog)
	l.Log("line")
	l.Log("line")
	// first flush
	l.Log("line")

	// force a second
	bufLog.Flush()

	if flushCount != 2 {
		t.Errorf("unexpected number of flushes: %d expected %d", flushCount, 2)
	}

	if flushedEntries != len(strings.Split(buf.String(), "\n"))-1 {
		t.Errorf("unexpected number of entries: %d expected %d", flushedEntries, 3)
	}
}

// outFile creates a real OS file for testing.
// We cannot use stdout/stderr since we need to read the contents afterwards to validate, and we have to write to a file
// to benchmark the impact of write() syscalls.
func outFile(b *testing.B) *os.File {
	f, err := os.CreateTemp(os.TempDir(), "linebuffer*")
	if err != nil {
		b.Fatalf("cannot create test file: %s", err)
	}

	return f
}
