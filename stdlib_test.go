package log

import (
	"bytes"
	"fmt"
	"log"
	"regexp"
	"testing"
	"time"
)

func TestStdlibWriter(t *testing.T) {
	buf := &bytes.Buffer{}
	log.SetOutput(buf)
	log.SetFlags(log.LstdFlags)
	logger := NewLogfmtLogger(StdlibWriter{})
	logger.Log("key", "val")
	timestamp := time.Now().Format("2006/01/02 15:04:05")
	if want, have := timestamp+" key=val\n", buf.String(); want != have {
		t.Errorf("want %q, have %q", want, have)
	}
}

func TestStdlibAdapterUsage(t *testing.T) {
	buf := &bytes.Buffer{}
	logger := NewLogfmtLogger(buf)
	writer := NewStdlibAdapter(logger)
	stdlog := log.New(writer, "", 0)

	now := time.Now()
	date := now.Format("2006/01/02")
	time := now.Format("15:04:05")

	for flag, want := range map[int]string{
		0:                                      "msg=hello\n",
		log.Ldate:                              "ts=" + date + " msg=hello\n",
		log.Ltime:                              "ts=" + time + " msg=hello\n",
		log.Ldate | log.Ltime:                  "ts=\"" + date + " " + time + "\" msg=hello\n",
		log.Lshortfile:                         "caller=stdlib_test.go:45 msg=hello\n",
		log.Lshortfile | log.Ldate:             "ts=" + date + " caller=stdlib_test.go:45 msg=hello\n",
		log.Lshortfile | log.Ldate | log.Ltime: "ts=\"" + date + " " + time + "\" caller=stdlib_test.go:45 msg=hello\n",
	} {
		buf.Reset()
		stdlog.SetFlags(flag)
		stdlog.Print("hello")
		if have := buf.String(); want != have {
			t.Errorf("flag=%d: want %#v, have %#v", flag, want, have)
		}
	}
}

func TestStdLibAdapter(t *testing.T) {
	var testcases = []struct {
		input      string
		regexp     *regexp.Regexp
		prefix     string
		joinPrefix bool
		want       string
	}{
		{
			regexp: StdlibRegexpFull,
			input:  `2009/01/23 01:23:23.123123 /a/b/c/d.go:23: hello world`,
			want:   `ts="2009/01/23 01:23:23.123123" caller=/a/b/c/d.go:23 msg="hello world"`,
		},
		{
			regexp: StdlibRegexpFull,
			input:  `hello`,
			want:   `msg=hello`,
		},
		{
			regexp: StdlibRegexpFull,
			input:  `2009/01/23: hello`,
			want:   `ts=2009/01/23 msg=hello`,
		},
		{
			regexp: StdlibRegexpFull,
			input:  `2009/01/23 01:23:23: hello`,
			want:   `ts="2009/01/23 01:23:23" msg=hello`,
		},
		{
			regexp: StdlibRegexpFull,
			input:  `01:23:23: hello`,
			want:   `ts=01:23:23 msg=hello`,
		},
		{
			regexp: StdlibRegexpFull,
			input:  `2009/01/23 01:23:23.123123: hello`,
			want:   `ts="2009/01/23 01:23:23.123123" msg=hello`,
		},
		{
			regexp: StdlibRegexpFull,
			input:  `2009/01/23 01:23:23.123123 /a/b/c/d.go:23: hello`,
			want:   `ts="2009/01/23 01:23:23.123123" caller=/a/b/c/d.go:23 msg=hello`,
		},
		{
			regexp: StdlibRegexpFull,
			input:  `01:23:23.123123 /a/b/c/d.go:23: hello`,
			want:   `ts=01:23:23.123123 caller=/a/b/c/d.go:23 msg=hello`,
		},
		{
			regexp: StdlibRegexpFull,
			input:  `2009/01/23 01:23:23 /a/b/c/d.go:23: hello`,
			want:   `ts="2009/01/23 01:23:23" caller=/a/b/c/d.go:23 msg=hello`,
		},
		{
			regexp: StdlibRegexpFull,
			input:  `2009/01/23 /a/b/c/d.go:23: hello`,
			want:   `ts=2009/01/23 caller=/a/b/c/d.go:23 msg=hello`,
		},
		{
			regexp: StdlibRegexpFull,
			input:  `/a/b/c/d.go:23: hello`,
			want:   `caller=/a/b/c/d.go:23 msg=hello`,
		},
{
			regexp: StdlibRegexpFull,
			input:  `some prefix hello`,
			prefix: "some prefix ",
			want:   `msg=hello`,
		},
		{
			regexp: StdlibRegexpFull,
			input:  `some prefix 2009/01/23: hello`,
			prefix: "some prefix ",
			want:   `ts=2009/01/23 msg=hello`,
		},
		{
			regexp: StdlibRegexpFull,
			input:  `some prefix 2009/01/23 01:23:23: hello`,
			prefix: "some prefix ",
			want:   `ts="2009/01/23 01:23:23" msg=hello`,
		},
		{
			regexp: StdlibRegexpFull,
			input:  `some prefix 01:23:23: hello`,
			prefix: "some prefix ",
			want:   `ts=01:23:23 msg=hello`,
		},
		{
			regexp: StdlibRegexpFull,
			input:  `some prefix 2009/01/23 01:23:23.123123: hello`,
			prefix: "some prefix ",
			want:   `ts="2009/01/23 01:23:23.123123" msg=hello`,
		},
		{
			regexp: StdlibRegexpFull,
			input:  `some prefix 2009/01/23 01:23:23.123123 /a/b/c/d.go:23: hello`,
			prefix: "some prefix ",
			want:   `ts="2009/01/23 01:23:23.123123" caller=/a/b/c/d.go:23 msg=hello`,
		},
		{
			regexp: StdlibRegexpFull,
			input:  `some prefix 01:23:23.123123 /a/b/c/d.go:23: hello`,
			prefix: "some prefix ",
			want:   `ts=01:23:23.123123 caller=/a/b/c/d.go:23 msg=hello`,
		},
		{
			regexp: StdlibRegexpFull,
			input:  `some prefix 2009/01/23 01:23:23 /a/b/c/d.go:23: hello`,
			prefix: "some prefix ",
			want:   `ts="2009/01/23 01:23:23" caller=/a/b/c/d.go:23 msg=hello`,
		},
		{
			regexp: StdlibRegexpFull,
			input:  `some prefix 2009/01/23 /a/b/c/d.go:23: hello`,
			prefix: "some prefix ",
			want:   `ts=2009/01/23 caller=/a/b/c/d.go:23 msg=hello`,
		},
		{
			regexp: StdlibRegexpFull,
			input:  `some prefix /a/b/c/d.go:23: hello`,
			prefix: "some prefix ",
			want:   `caller=/a/b/c/d.go:23 msg=hello`,
		},
		{
			regexp: StdlibRegexpFull,
			input:  `/a/b/c/d.go:23: some prefix hello`,
			prefix: "some prefix ",
			want:   `caller=/a/b/c/d.go:23 msg=hello`,
		},
		{
			regexp:     StdlibRegexpFull,
			input:      `some prefix hello`,
			prefix:     "some prefix ",
			joinPrefix: true,
			want:       `msg="some prefix hello"`,
		},
		{
			regexp:     StdlibRegexpFull,
			input:      `some prefix 2009/01/23: hello`,
			prefix:     "some prefix ",
			joinPrefix: true,
			want:       `ts=2009/01/23 msg="some prefix hello"`,
		},
		{
			regexp:     StdlibRegexpFull,
			input:      `some prefix 2009/01/23 01:23:23: hello`,
			prefix:     "some prefix ",
			joinPrefix: true,
			want:       `ts="2009/01/23 01:23:23" msg="some prefix hello"`,
		},
		{
			regexp:     StdlibRegexpFull,
			input:      `some prefix 01:23:23: hello`,
			prefix:     "some prefix ",
			joinPrefix: true,
			want:       `ts=01:23:23 msg="some prefix hello"`,
		},
		{
			regexp:     StdlibRegexpFull,
			input:      `some prefix 2009/01/23 01:23:23.123123: hello`,
			prefix:     "some prefix ",
			joinPrefix: true,
			want:       `ts="2009/01/23 01:23:23.123123" msg="some prefix hello"`,
		},
		{
			regexp:     StdlibRegexpFull,
			input:      `some prefix 2009/01/23 01:23:23.123123 /a/b/c/d.go:23: hello`,
			prefix:     "some prefix ",
			joinPrefix: true,
			want:       `ts="2009/01/23 01:23:23.123123" caller=/a/b/c/d.go:23 msg="some prefix hello"`,
		},
		{
			regexp:     StdlibRegexpFull,
			input:      `some prefix 01:23:23.123123 /a/b/c/d.go:23: hello`,
			prefix:     "some prefix ",
			joinPrefix: true,
			want:       `ts=01:23:23.123123 caller=/a/b/c/d.go:23 msg="some prefix hello"`,
		},
		{
			regexp:     StdlibRegexpFull,
			input:      `some prefix 2009/01/23 01:23:23 /a/b/c/d.go:23: hello`,
			prefix:     "some prefix ",
			joinPrefix: true,
			want:       `ts="2009/01/23 01:23:23" caller=/a/b/c/d.go:23 msg="some prefix hello"`,
		},
		{
			regexp:     StdlibRegexpFull,
			input:      `some prefix 2009/01/23 /a/b/c/d.go:23: hello`,
			prefix:     "some prefix ",
			joinPrefix: true,
			want:       `ts=2009/01/23 caller=/a/b/c/d.go:23 msg="some prefix hello"`,
		},
		{
			regexp:     StdlibRegexpFull,
			input:      `some prefix /a/b/c/d.go:23: hello`,
			prefix:     "some prefix ",
			joinPrefix: true,
			want:       `caller=/a/b/c/d.go:23 msg="some prefix hello"`,
		},
		{
			regexp:     StdlibRegexpFull,
			input:      `/a/b/c/d.go:23: some prefix hello`,
			prefix:     "some prefix ",
			joinPrefix: true,
			want:       `caller=/a/b/c/d.go:23 msg="some prefix hello"`,
		},
		{
			regexp: StdlibRegexpDefault,
			input:  `error encoding and sending metric family: write tcp 127.0.0.1:9182->127.0.0.1:60125: wsasend:`,
			want:   `msg="error encoding and sending metric family: write tcp 127.0.0.1:9182->127.0.0.1:60125: wsasend:"`,
		},
		{
			regexp: StdlibRegexpDefault,
			input:  `2023/04/28 07:28:46 error encoding and sending metric family: write tcp 127.0.0.1:9182->127.0.0.1:60125: wsasend:`,
			want:   `ts="2023/04/28 07:28:46" msg="error encoding and sending metric family: write tcp 127.0.0.1:9182->127.0.0.1:60125: wsasend:"`,
		},
		{
			regexp: StdlibRegexpDefault,
			input:  `2023/04/28 07:28:46 /a/b/c/d.go:23: error encoding and sending metric family: write tcp 127.0.0.1:9182->127.0.0.1:60125: wsasend:`,
			want:   `ts="2023/04/28 07:28:46" msg="/a/b/c/d.go:23: error encoding and sending metric family: write tcp 127.0.0.1:9182->127.0.0.1:60125: wsasend:"`,
		},
		{
			regexp: StdlibRegexpDefault,
			input:  `2009/01/23 01:23:23.123123 /a/b/c/d.go:23: hello`,
			want:   `ts="2009/01/23 01:23:23.123123" msg="/a/b/c/d.go:23: hello"`,
		},
		{
			regexp: StdlibRegexpDefault,
			input:  `1:9182f`,
			want:   `msg=1:9182f`,
		},
	}
	for _, tt := range testcases {
		t.Run(tt.input, func(t *testing.T) {
			var buf bytes.Buffer
			adapter := NewStdlibAdapter(NewLogfmtLogger(&buf), StdlibRegexp(tt.regexp), Prefix(tt.prefix, tt.joinPrefix))
			fmt.Fprint(adapter, tt.input)

			if want, have := tt.want+"\n", buf.String(); want != have {
				t.Errorf("%q: want %q, have %q", tt.input, want, have)
			}
		})
	}
}
