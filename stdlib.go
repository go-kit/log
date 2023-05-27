package log

import (
	"bytes"
	"io"
	"log"
	"regexp"
	"strings"
)

// StdlibWriter implements io.Writer by invoking the stdlib log.Print. It's
// designed to be passed to a Go kit logger as the writer, for cases where
// it's necessary to redirect all Go kit log output to the stdlib logger.
//
// If you have any choice in the matter, you shouldn't use this. Prefer to
// redirect the stdlib log to the Go kit logger via NewStdlibAdapter.
type StdlibWriter struct{}

// Write implements io.Writer.
func (w StdlibWriter) Write(p []byte) (int, error) {
	log.Print(strings.TrimSpace(string(p)))
	return len(p), nil
}

// StdlibAdapter wraps a Logger and allows it to be passed to the stdlib
// logger's SetOutput. It will extract date/timestamps, filenames, and
// messages, and place them under relevant keys.
type StdlibAdapter struct {
	Logger
	timestampKey    string
	fileKey         string
	messageKey      string
	prefix          string
	joinPrefixToMsg bool
	logRegexp       *regexp.Regexp
}

// StdlibAdapterOption sets a parameter for the StdlibAdapter.
type StdlibAdapterOption func(*StdlibAdapter)

// TimestampKey sets the key for the timestamp field. By default, it's "ts".
func TimestampKey(key string) StdlibAdapterOption {
	return func(a *StdlibAdapter) { a.timestampKey = key }
}

// FileKey sets the key for the file and line field. By default, it's "caller".
func FileKey(key string) StdlibAdapterOption {
	return func(a *StdlibAdapter) { a.fileKey = key }
}

// MessageKey sets the key for the actual log message. By default, it's "msg".
func MessageKey(key string) StdlibAdapterOption {
	return func(a *StdlibAdapter) { a.messageKey = key }
}

// StdlibRegexp sets the regular expression used to parse stdlib log messages.
// Nil regexps are ignored and will return options that are no-ops. The default
// value is StdlibRegexpFull.
func StdlibRegexp(re *regexp.Regexp) StdlibAdapterOption {
	if re == nil {
		return func(a *StdlibAdapter) {}
	}
	return func(a *StdlibAdapter) { a.logRegexp = re }
}

// Prefix configures the adapter to parse a prefix from stdlib log events. If
// you provide a non-empty prefix to the stdlib logger, then your should provide
// that same prefix to the adapter via this option.
//
// By default, the prefix isn't included in the msg key. Set joinPrefixToMsg to
// true if you want to include the parsed prefix in the msg.
func Prefix(prefix string, joinPrefixToMsg bool) StdlibAdapterOption {
	return func(a *StdlibAdapter) { a.prefix = prefix; a.joinPrefixToMsg = joinPrefixToMsg }
}

// NewStdlibAdapter returns a new StdlibAdapter wrapper around the passed
// logger. It's designed to be passed to log.SetOutput.
func NewStdlibAdapter(logger Logger, options ...StdlibAdapterOption) io.Writer {
	a := StdlibAdapter{
		Logger:       logger,
		timestampKey: "ts",
		fileKey:      "caller",
		messageKey:   "msg",
		logRegexp:    StdlibRegexpFull,
	}
	for _, option := range options {
		option(&a)
	}

	return a
}

func (a StdlibAdapter) Write(p []byte) (int, error) {
	p = a.handlePrefix(p)

	result := a.subexps(p)
	keyvals := []interface{}{}
	var timestamp string
	if date, ok := result["date"]; ok && date != "" {
		timestamp = date
	}
	if time, ok := result["time"]; ok && time != "" {
		if timestamp != "" {
			timestamp += " "
		}
		timestamp += time
	}
	if timestamp != "" {
		keyvals = append(keyvals, a.timestampKey, timestamp)
	}
	if file, ok := result["file"]; ok && file != "" {
		keyvals = append(keyvals, a.fileKey, file)
	}
	if msg, ok := result["msg"]; ok {
		msg = a.handleMessagePrefix(msg)
		keyvals = append(keyvals, a.messageKey, msg)
	}
	if err := a.Logger.Log(keyvals...); err != nil {
		return 0, err
	}
	return len(p), nil
}

func (a StdlibAdapter) handlePrefix(p []byte) []byte {
	if a.prefix != "" {
		p = bytes.TrimPrefix(p, []byte(a.prefix))
	}
	return p
}

func (a StdlibAdapter) handleMessagePrefix(msg string) string {
	if a.prefix == "" {
		return msg
	}

	msg = strings.TrimPrefix(msg, a.prefix)
	if a.joinPrefixToMsg {
		msg = a.prefix + msg
	}
	return msg
}

const (
	stdlibRegexpPatternDate = `(?P<date>[0-9]{4}/[0-9]{2}/[0-9]{2})?[ ]?`
	stdlibRegexpPatternTime = `(?P<time>[0-9]{2}:[0-9]{2}:[0-9]{2}(\.[0-9]+)?)?[ ]?`
	stdlibRegexpPatternFile = `(?P<file>.+?:[0-9]+)?`
	stdlibRegexpPatternMsg  = `(: )?(?P<msg>(?s:.*))`
)

var (
	// StdlibRegexpFull captures date, time, caller (file), and message from stdlib log messages.
	StdlibRegexpFull = regexp.MustCompile(stdlibRegexpPatternDate + stdlibRegexpPatternTime + stdlibRegexpPatternFile + stdlibRegexpPatternMsg)
	// StdlibRegexpDefault captures date, time and message from stdlib log messages.
	StdlibRegexpDefault = regexp.MustCompile(stdlibRegexpPatternDate + stdlibRegexpPatternTime + stdlibRegexpPatternMsg)
)

func (a StdlibAdapter) subexps(line []byte) map[string]string {
	m := a.logRegexp.FindStringSubmatch(string(line))
	n := a.logRegexp.SubexpNames()
	if len(m) < len(n) {
		return map[string]string{}
	}
	result := map[string]string{}
	for i, name := range n {
		result[name] = strings.TrimRight(m[i], "\n")
	}
	return result
}
