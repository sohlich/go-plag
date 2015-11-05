package main

import (
	"errors"
	"expvar"
)

const (
	//DatabaseOK  is databse state if everything is OK.
	DatabaseOK = "OK"
	//DatabaseNotConnected  is database
	//state if some errors in connection occurs
	DatabaseNotConnected = "NOT CONNECTED"
)

var (
	//ErrMetricsNotInitialized occurs if
	//metrics are not analyzed
	ErrMetricsNotInitialized = errors.New("Metrics not initialized")
)

//MetricFunc is function
//executed if metrics are initi
type metricsFunc func()

//Metrics struct
//provides interface to
//hold metrics
type Metrics struct {
	initialized bool
	comparison  *expvar.Int
	errors      *expvar.Int
	database    *expvar.String
}

//NewMetrics creates
//Metrics object
func NewMetrics() *Metrics {
	metrics = &Metrics{
		true,
		expvar.NewInt("comparisons"),
		expvar.NewInt("errors"),
		expvar.NewString("database"),
	}
	return metrics
}

//AddComparison increments
//the number of comparisons in
//metrics object by given value
func (m *Metrics) AddComparison(value int64) error {
	return m.doIfNotNil(func() {
		m.comparison.Add(value)
	})
}

//AddError increments the
//value of errors
func (m *Metrics) AddError(value int64) error {
	return m.doIfNotNil(func() {
		m.errors.Add(value)
	})
}

//ComparisonInc increments the
//value of comparisons
//in matrics object by 1
func (m *Metrics) ComparisonInc() error {
	return m.doIfNotNil(func() {
		m.comparison.Add(1)
	})
}

//ErrorInc increments the value
//of error in metrics struc by 1
func (m *Metrics) ErrorInc() error {
	return m.doIfNotNil(func() {
		m.errors.Add(1)
	})
}

//SetDatabaseState sets the state of database
//for displaying in health check
func (m *Metrics) SetDatabaseState(state string) error {
	return m.doIfNotNil(func() {
		m.database.Set(state)
	})
}

//GetDatabaseState retrieves the state of database
//for displaying in health check
func (m *Metrics) GetDatabaseState() string {
	if m.errors != nil {
		return m.database.String()
	}
	return ""
}

//GetErrors returns the number of
//errors
func (m *Metrics) GetErrors() string {
	if m.errors != nil {
		return m.errors.String()
	}
	return ""
}

//GetComparisons returns the string value
//of comparison count
func (m *Metrics) GetComparisons() string {
	if m.comparison != nil {
		return m.comparison.String()
	}
	return ""
}

func (m *Metrics) doIfNotNil(function metricsFunc) error {
	if m != nil && m.initialized {
		function()
		return nil
	}
	return ErrMetricsNotInitialized
}
