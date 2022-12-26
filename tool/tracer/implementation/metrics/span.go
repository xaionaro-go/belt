package metrics

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/facebookincubator/go-belt"
	"github.com/facebookincubator/go-belt/pkg/field"
	"github.com/facebookincubator/go-belt/tool/experimental/tracer"
)

func init() {
	if tracer.Default == nil {
		tracer.Default = func() tracer.Tracer {
			return Default()
		}
	}
}

// Span is the implementation of tracer.Span.
type Span struct {
	sendOnce sync.Once
	parent   tracer.Span
	name     string
	tracer   *Tracer
	startTS  time.Time
}

var _ tracer.Span = (*Span)(nil)

func (t *Tracer) newSpanBelt(
	belt *belt.Belt,
	name string,
	parent tracer.Span,
	options ...tracer.SpanOption,
) (tracer.Span, *belt.Belt) {
	if tracer.IsNoopSpan(parent) {
		return tracer.NewNoopSpan(name, parent, time.Now()), belt
	}

	span := &Span{
		parent:  parent,
		tracer:  t,
		name:    name,
		startTS: time.Now(),
	}

	if !t.PreHooks.ProcessSpan(span) {
		return tracer.NewNoopSpan(name, parent, time.Now()), belt
	}
	if belt == nil {
		return span, nil
	}
	return span, tracer.BeltWithSpan(belt, span)
}

func (t *Tracer) newSpanCtx(
	ctx context.Context,
	name string,
	parent tracer.Span,
	options ...tracer.SpanOption,
) (tracer.Span, context.Context) {
	if tracer.IsNoopSpan(parent) {
		return tracer.NewNoopSpan(name, parent, time.Now()), ctx
	}

	span := &Span{
		parent:  parent,
		tracer:  t,
		name:    name,
		startTS: time.Now(),
	}

	if !t.PreHooks.ProcessSpan(span) {
		return tracer.NewNoopSpan(name, parent, time.Now()), ctx
	}
	if ctx == nil {
		return span, nil
	}
	return span, tracer.CtxWithSpan(ctx, span)
}

// ID implements tracer.Span.
func (span *Span) ID() any {
	// not implemented
	return fmt.Sprintf("%s-%d", span.name, span.startTS)
}

// Name implements tracer.Span.
func (span *Span) Name() string {
	return span.name
}

// TraceIDs implements tracer.Span.
func (span *Span) TraceIDs() belt.TraceIDs {
	// not implemented
	return nil
}

// StartTS implements tracer.Span.
func (span *Span) StartTS() time.Time {
	return span.startTS
}

// Fields implements tracer.Span.
func (span *Span) Fields() field.AbstractFields {
	// not implemented
	return nil
}

// Parent implements tracer.Span.
func (span *Span) Parent() tracer.Span {
	return span.parent
}

// SetName implements tracer.Span.
func (span *Span) SetName(name string) {
	span.name = name
}

// Annotate implements tracer.Span.
func (span *Span) Annotate(ts time.Time, event string) {
	// not implemented
}

// SetField implements tracer.Span.
func (span *Span) SetField(k field.Key, v field.Value) {
	// not implemented
}

// SetFields implements tracer.Span.
func (span *Span) SetFields(fields field.AbstractFields) {
	span.tracer = span.tracer.WithContextFields((*field.FieldsChain)(nil).WithFields(fields), fields.Len()).(*Tracer)
}

// Finish implements tracer.Span.
func (span *Span) Finish() {
	duration := time.Since(span.startTS)
	span.sendOnce.Do(func() {
		if !span.tracer.Hooks.ProcessSpan(span) {
			return
		}
		span.send(duration)
	})
}

// FinishWithDuration implements tracer.Span.
func (span *Span) FinishWithDuration(duration time.Duration) {
	span.sendOnce.Do(func() {
		if !span.tracer.Hooks.ProcessSpan(span) {
			return
		}
		span.send(duration)
	})
}

func (span *Span) send(duration time.Duration) {
	// TODO: use histograms
	key := "spanDuration_" + span.name
	span.tracer.Backend.IntGauge(key).Add(duration.Nanoseconds())
	span.tracer.Backend.Count(key).Add(1)
}

// Flush implements tracer.Span.
func (span *Span) Flush() {
	span.tracer.Flush()
}
