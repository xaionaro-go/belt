package metrics

import (
	"context"

	"github.com/facebookincubator/go-belt"
	"github.com/facebookincubator/go-belt/pkg/field"
	"github.com/facebookincubator/go-belt/tool/experimental/metrics"
	"github.com/facebookincubator/go-belt/tool/experimental/tracer"
)

type Tracer struct {
	Backend  metrics.Metrics
	Hooks    tracer.Hooks
	PreHooks tracer.Hooks
}

var _ tracer.Tracer = (*Tracer)(nil)

func New(metrics metrics.Metrics) *Tracer {
	return &Tracer{
		Backend: metrics,
	}
}

func (t Tracer) clone() *Tracer {
	return &t
}

// Flush implements belt.Tool.
func (t *Tracer) Flush() {
	t.Backend.Flush()
}

// WithContextFields implements belt.Tool.
func (t *Tracer) WithContextFields(allFields *field.FieldsChain, newFieldsCount int) belt.Tool {
	c := t.clone()
	c.Backend = t.Backend.WithContextFields(allFields, newFieldsCount).(metrics.Metrics)
	return c
}

// WithTraceIDs implements belt.Tool.
func (t *Tracer) WithTraceIDs(traceIDs belt.TraceIDs, newTraceIDsCount int) belt.Tool {
	c := t.clone()
	c.Backend = t.Backend.WithTraceIDs(traceIDs, newTraceIDsCount).(metrics.Metrics)
	return c
}

// Start implements tracer.Tracer.
func (t *Tracer) Start(name string, parent tracer.Span, options ...tracer.SpanOption) tracer.Span {
	span, _ := t.newSpanBelt(nil, name, parent, options...)
	return span
}

// StartWithBelt implements tracer.Tracer.
func (t *Tracer) StartWithBelt(belt *belt.Belt, name string, options ...tracer.SpanOption) (tracer.Span, *belt.Belt) {
	return t.newSpanBelt(belt, name, nil, options...)
}

// StartChildWithBelt implements tracer.Tracer.
func (t *Tracer) StartChildWithBelt(belt *belt.Belt, name string, options ...tracer.SpanOption) (tracer.Span, *belt.Belt) {
	return t.newSpanBelt(belt, name, tracer.SpanFromBelt(belt), options...)
}

// StartWithCtx implements tracer.Tracer.
func (t *Tracer) StartWithCtx(ctx context.Context, name string, options ...tracer.SpanOption) (tracer.Span, context.Context) {
	return t.newSpanCtx(ctx, name, nil, options...)
}

// StartChildWithCtx implements tracer.Tracer.
func (t *Tracer) StartChildWithCtx(ctx context.Context, name string, options ...tracer.SpanOption) (tracer.Span, context.Context) {
	return t.newSpanCtx(ctx, name, tracer.SpanFromCtx(ctx), options...)
}

// WithPreHooks implements tracer.Tracer.
func (t *Tracer) WithPreHooks(newPreHooks ...tracer.Hook) tracer.Tracer {
	c := t.clone()
	if newPreHooks == nil {
		c.PreHooks = nil
	} else {
		var preHooks tracer.Hooks
		preHooks = append(preHooks, t.PreHooks...)
		preHooks = append(preHooks, newPreHooks...)
		c.PreHooks = preHooks
	}
	return c
}

// WithHooks implements tracer.Tracer.
func (t *Tracer) WithHooks(newHooks ...tracer.Hook) tracer.Tracer {
	c := t.clone()
	if newHooks == nil {
		c.Hooks = nil
	} else {
		var hooks tracer.Hooks
		hooks = append(hooks, t.Hooks...)
		hooks = append(hooks, newHooks...)
		c.Hooks = hooks
	}
	return c
}

var Default = func() *Tracer {
	return New(metrics.Default())
}
