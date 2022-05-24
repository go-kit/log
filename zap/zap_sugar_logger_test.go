package zap_test

import (
	"encoding/json"
	"strings"
	"testing"

	kitlog "github.com/go-kit/log"
	"github.com/go-kit/log/level"
	kitzap "github.com/go-kit/log/zap"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func TestZapSugarLogger(t *testing.T) {
	// logger config
	encoderConfig := zap.NewDevelopmentEncoderConfig()
	encoder := zapcore.NewJSONEncoder(encoderConfig)
	levelKey := encoderConfig.LevelKey
	// basic test cases
	type testCase struct {
		level zapcore.Level
		log   func(kitlog.Logger)
		want  map[string]string
	}
	testCases := []testCase{
		//default level test
		{level: zapcore.DebugLevel,
			log: func(l kitlog.Logger) {
				l.Log("key1", "value1")
			},
			want: map[string]string{levelKey: "DEBUG", "key1": "value1"}},

		{level: zapcore.InfoLevel,
			log: func(l kitlog.Logger) {
				l.Log("key2", "value2")
			},
			want: map[string]string{levelKey: "INFO", "key2": "value2"}},

		{level: zapcore.WarnLevel,
			log: func(l kitlog.Logger) {
				l.Log("key3", "value3")
			},
			want: map[string]string{levelKey: "WARN", "key3": "value3"}},

		{level: zapcore.ErrorLevel,
			log: func(l kitlog.Logger) {
				l.Log("key4", "value4")
			},
			want: map[string]string{levelKey: "ERROR", "key4": "value4"}},

		{level: zapcore.DPanicLevel,
			log: func(l kitlog.Logger) {
				l.Log("key5", "value5")
			},
			want: map[string]string{levelKey: "DPANIC", "key5": "value5"}},

		{level: zapcore.PanicLevel,
			log: func(l kitlog.Logger) {
				l.Log("key6", "value6")
			},
			want: map[string]string{levelKey: "PANIC", "key6": "value6"}},

		// kitlog level test
		{level: zapcore.ErrorLevel,
			log: func(l kitlog.Logger) {
				level.Debug(l).Log("key1", "value1")
			},
			want: map[string]string{levelKey: "DEBUG", "key1": "value1"}},

		{level: zapcore.ErrorLevel,
			log: func(l kitlog.Logger) {
				level.Info(l).Log("key2", "value2")
			},
			want: map[string]string{levelKey: "INFO", "key2": "value2"}},

		{level: zapcore.InfoLevel,
			log: func(l kitlog.Logger) {
				level.Warn(l).Log("key3", "value3")
			},
			want: map[string]string{levelKey: "WARN", "key3": "value3"}},

		{level: zapcore.InfoLevel,
			log: func(l kitlog.Logger) {
				level.Error(l).Log("key4", "value4")
			},
			want: map[string]string{levelKey: "ERROR", "key4": "value4"}},
	}
	// test
	for _, testCase := range testCases {
		t.Run(testCase.level.String(), func(t *testing.T) {
			// make logger
			writer := &tbWriter{tb: t}
			logger := zap.New(
				zapcore.NewCore(encoder, zapcore.AddSync(writer), zap.DebugLevel),
				zap.Development())
			// check panic
			shouldPanic := testCase.level >= zapcore.DPanicLevel && (testCase.want[levelKey] == "PANIC" || testCase.want[levelKey] == "DPANIC")
			kitLogger := kitzap.NewZapSugarLogger(logger, testCase.level)
			defer func() {
				isPanic := recover() != nil
				if shouldPanic != isPanic {
					t.Errorf("test level %v should panic(%v), but %v", testCase.level, shouldPanic, isPanic)
				}
				// check log kvs
				logMap := make(map[string]string)
				err := json.Unmarshal([]byte(writer.sb.String()), &logMap)
				if err != nil {
					t.Errorf("unmarshal error: %v", err)
				} else {
					for k, v := range testCase.want {
						vv, ok := logMap[k]
						if !ok || v != vv {
							t.Error("error log")
						}
					}
				}
			}()
			testCase.log(kitLogger)
		})
	}
}

type tbWriter struct {
	tb testing.TB
	sb strings.Builder
}

func (w *tbWriter) Write(b []byte) (n int, err error) {
	w.tb.Logf(string(b))
	w.sb.Write(b)
	return len(b), nil
}
