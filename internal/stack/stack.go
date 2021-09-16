package stack

import "runtime"

// Caller returns a frame from the stack of the current goroutine. The argument
// skip is the number of frames to ascend, with 0 identifying the
// calling function.
func Caller(skip int) runtime.Frame {
	var pcs [3]uintptr

	n := runtime.Callers(skip+1, pcs[:])
	frames := runtime.CallersFrames(pcs[:n])
	frame, _ := frames.Next()
	frame, _ = frames.Next()

	return frame
}

// Trace returns a call stack for the current goroutine with element 0
// identifying the calling function.
func Trace() []runtime.Frame {
	var pcs [512]uintptr
	n := runtime.Callers(1, pcs[:])

	frames := runtime.CallersFrames(pcs[:n])
	cs := make([]runtime.Frame, 0, n)

	// Skip extra frame retrieved just to make sure the runtime.sigpanic
	// special case is handled.
	frame, more := frames.Next()

	for more {
		frame, more = frames.Next()
		cs = append(cs, frame)
	}

	return cs
}
