package log

import stdctx "context"

// WithContext returns a shallow copy of l with its context changed
// to ctx. The provided ctx must be non-nil.
func WithContext(ctx stdctx.Context, l Logger) Logger {
	if c, ok := l.(*context); ok {
		return &context{
			logger:     c.logger,
			keyvals:    c.keyvals,
			sKeyvals:   c.sKeyvals,
			hasValuer:  c.hasValuer,
			sHasValuer: c.sHasValuer,
			ctx:        ctx,
		}
	}

	ret := newContext(l)
	ret.ctx = ctx
	return ret
}
