package tests

import (
	"github.com/facebookincubator/go-belt"
	"github.com/facebookincubator/go-belt/pkg/field"
	"github.com/facebookincubator/go-belt/tool/experimental/metrics"
	"github.com/facebookincubator/go-belt/tool/experimental/tracer"
	metricsimpl "github.com/xaionaro-go/belt/tool/tracer/implementation/metrics"
)

type dummyMetrics struct {
	onSend func(span tracer.Span)
	tracer *metricsimpl.Tracer
}

var _ metrics.Metrics = (*dummyMetrics)(nil)

func (l *dummyMetrics) Flush() {}
func (l *dummyMetrics) WithContextFields(allFields *field.FieldsChain, newFieldsCount int) belt.Tool {
	panic("unexpected call of WithContextFields")
}
func (l *dummyMetrics) WithTraceIDs(traceIDs belt.TraceIDs, newTraceIDsCount int) belt.Tool {
	panic("unexpected call of WithTraceIDs")
}

func (l *dummyMetrics) Gauge(key string) metrics.Gauge {
	panic("unexpected call of Gauge")
}
func (l *dummyMetrics) GaugeFields(key string, additionalFields field.AbstractFields) metrics.Gauge {
	panic("unexpected call of GaugeFields")
}

type dummyIntGauge struct{}

func (m dummyIntGauge) Value() any {
	panic("unexpected call of Value")
}
func (m dummyIntGauge) Add(v int64) metrics.IntGauge {
	return m
}
func (m dummyIntGauge) WithResetFields(field.AbstractFields) metrics.IntGauge {
	panic("unexpected call of WithResetFields")
}

func (l *dummyMetrics) IntGauge(key string) metrics.IntGauge {
	return &dummyIntGauge{}
}
func (l *dummyMetrics) IntGaugeFields(key string, additionalFields field.AbstractFields) metrics.IntGauge {
	panic("unexpected call of IntGaugeFields")
}

type dummyCount struct{}

func (m dummyCount) Value() any {
	panic("unexpected call of Value")
}
func (m dummyCount) Add(v uint64) metrics.Count {
	return m
}
func (m dummyCount) WithResetFields(field.AbstractFields) metrics.Count {
	panic("unexpected call of WithResetFields")
}

func (l *dummyMetrics) Count(key string) metrics.Count {
	span := &metricsimpl.Span{}
	span.SetName(key[len("spanDuration_"):])
	l.onSend(span)
	return &dummyCount{}
}
func (l *dummyMetrics) CountFields(key string, additionalFields field.AbstractFields) metrics.Count {
	panic("unexpected call of CountFields")
}

func (l *dummyMetrics) OnSend(onSend func(span tracer.Span)) {
	l.onSend = onSend
}

func init() {
	implementations = append(implementations, implementationCase{
		Name: "metrics",
		Factory: func() (tracer.Tracer, DummyReporter) {
			reporter := &dummyMetrics{}
			tracer := metricsimpl.New(reporter)
			reporter.tracer = tracer
			if tracer == nil {
				panic("tracer == nil")
			}
			return tracer, reporter
		},
	})
}
