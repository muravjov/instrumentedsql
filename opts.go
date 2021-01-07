package instrumentedsql

import (
	"context"
	"time"
)

type opts struct {
	Logger
	Tracer
	OpsExcluded map[string]struct{}
	OmitArgs    bool
	TimeoutFunc func() time.Duration
}

// Opt is a functional option type for the wrapped driver
type Opt func(*opts)

func (o *opts) hasOpExcluded(op string) bool {
	_, ok := o.OpsExcluded[op]
	return ok
}

// WithLogger sets the logger of the wrapped driver to the provided logger
func WithLogger(l Logger) Opt {
	return func(o *opts) {
		o.Logger = l
	}
}

// WithOpsExcluded excludes some of OpSQL that are not required
func WithOpsExcluded(ops ...string) Opt {
	return func(o *opts) {
		o.OpsExcluded = make(map[string]struct{})
		for _, op := range ops {
			o.OpsExcluded[op] = struct{}{}
		}
	}
}

// WithTracer sets the tracer of the wrapped driver to the provided tracer
func WithTracer(t Tracer) Opt {
	return func(o *opts) {
		o.Tracer = t
	}
}

// WithOmitArgs will make it so that query arguments are omitted from logging and tracing
func WithOmitArgs() Opt {
	return func(o *opts) {
		o.OmitArgs = true
	}
}

// WithIncludeArgs will make it so that query arguments are included in logging and tracing
// This is the default, but can be used to override WithOmitArgs
func WithIncludeArgs() Opt {
	return func(o *opts) {
		o.OmitArgs = false
	}
}

// WithTimeout sets timeout on every driver's operation except driver.ConnBeginTx.BeginTx
func WithTimeoutFunc(f func() time.Duration) Opt {
	return func(o *opts) {
		o.TimeoutFunc = f
	}
}

func (o *opts) setTimeout(ctx context.Context) (context.Context, context.CancelFunc) {
	if f := o.TimeoutFunc; f != nil {
		if timeout := f(); timeout != 0 {
			return context.WithTimeout(ctx, timeout)
		}
	}
	return ctx, func() {}
}
