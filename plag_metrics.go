package main

import (
	"errors"
	"expvar"
)

var (
	MetricsNotInitializedError = errors.New("Metrics not initialized")
)

type MetricsFunc func()

type Metrics struct {
	initialized bool
	comparison  *expvar.Int
	errors      *expvar.Int
}

func NewMetrics() *Metrics {
	metrics = &Metrics{
		true,
		expvar.NewInt("comparisons"),
		expvar.NewInt("errors"),
	}
	return metrics
}

func (m *Metrics) AddComparison(value int64) error {
	return m.doIfNotNil(func() {
		m.comparison.Add(value)
	})
}

func (m *Metrics) AddError(value int64) error {
	return m.doIfNotNil(func() {
		m.errors.Add(value)
	})
}

func (m *Metrics) ComparisonInc() error {
	return m.doIfNotNil(func() {
		m.comparison.Add(1)
	})
}

func (m *Metrics) ErrorInc() error {
	return m.doIfNotNil(func() {
		m.errors.Add(1)
	})
}

func (m *Metrics) GerErrors() string {
	if m.errors != nil {
		return m.errors.String()
	}
	return ""
}

func (m *Metrics) GerComparisons() string {
	if m.comparison != nil {
		return m.comparison.String()
	}
	return ""
}

func (m *Metrics) doIfNotNil(function MetricsFunc) error {
	if m != nil && m.initialized {
		function()
		return nil
	}
	return MetricsNotInitializedError
}
