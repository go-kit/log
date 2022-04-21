package log

import (
	"bytes"
	stdctx "context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestWithContext(t *testing.T) {
	ctxKey := "key-in-ctx"
	ctx := stdctx.Background()
	ctx = stdctx.WithValue(ctx, ctxKey, "val-in-ctx")

	getVal := func(ctx stdctx.Context) interface{} {
		return ctx.Value(ctxKey)
	}

	baseTime := time.Date(2015, time.February, 3, 10, 0, 0, 0, time.UTC)
	mockTime := func() time.Time {
		baseTime = baseTime.Add(time.Second)
		return baseTime
	}

	buf := &bytes.Buffer{}
	logger := NewLogfmtLogger(buf)

	logger = With(logger, "ts", Timestamp(mockTime), ctxKey, Valuer(getVal))

	logger = WithContext(ctx, logger)

	logger.Log("key", "val")

	got := buf.String()

	want := "ts=2015-02-03T10:00:01Z key-in-ctx=val-in-ctx key=val\n"
	require.Equal(t, want, got)

	ctx = stdctx.WithValue(stdctx.Background(), ctxKey, "val2-in-ctx")
	logger = WithContext(ctx, logger)
	logger.Log("key", "val")

	got = buf.String()
	want = "ts=2015-02-03T10:00:01Z key-in-ctx=val-in-ctx key=val\nts=2015-02-03T10:00:02Z key-in-ctx=val2-in-ctx key=val\n"
	require.Equal(t, want, got)
}
