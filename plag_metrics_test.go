package main

import (
	"testing"
)

func TestMetricsNotInitialized(t *testing.T) {
	metrics := &Metrics{}
	err := metrics.AddError(1)
	if err == nil {
		t.Error("Validation of not initialized Metrics failed")
	}
}

func TestNewMetrics(t *testing.T) {
	metrics = NewMetrics()
	if metrics.GerErrors() != "0" {
		t.Error("Metrics does not increment errors")
	}

	if metrics.GerComparisons() != "0" {
		t.Error("Metrics does not increment comparisons")
	}
}
