package decorator_test

import (
	"context"
	"testing"

	"github.com/prometheus/client_golang/prometheus"
	io_prometheus_client "github.com/prometheus/client_model/go"
	"github.com/yuppyweb/cakelog/decorator"
)

type mockPrometheusCounter struct {
	inc float64
}

func (mockPrometheusCounter) Collect(chan<- prometheus.Metric) {}

func (mockPrometheusCounter) Desc() *prometheus.Desc {
	return nil
}

func (mockPrometheusCounter) Describe(chan<- *prometheus.Desc) {}

func (mockPrometheusCounter) Write(*io_prometheus_client.Metric) error {
	return nil
}

func (m *mockPrometheusCounter) Inc() {
	m.inc++
}

func (m *mockPrometheusCounter) Add(num float64) {
	m.inc += num
}

func TestPrometheusLogger_Debug(t *testing.T) {
	debugCounter := new(mockPrometheusCounter)
	infoCounter := new(mockPrometheusCounter)
	warnCounter := new(mockPrometheusCounter)
	errorCounter := new(mockPrometheusCounter)
	mockLogger := new(mockLogger)

	logger := decorator.NewPrometheusLogger(mockLogger, decorator.PrometheusLoggerCounter{
		Debug: debugCounter,
		Info:  infoCounter,
		Warn:  warnCounter,
		Error: errorCounter,
	})

	logger.Debug(context.Background(), "debug message", "debug", 42)

	if len(mockLogger.debugIn) != 1 {
		t.Fatalf("Expected Debug to be called once, got %d", len(mockLogger.debugIn))
	}

	if mockLogger.debugIn[0].msg != "debug message" {
		t.Errorf("Expected message to be 'debug message', got '%s'", mockLogger.debugIn[0].msg)
	}

	if len(mockLogger.debugIn[0].args) != 2 {
		t.Fatalf("Expected 2 arguments, got %d", len(mockLogger.debugIn[0].args))
	}

	if mockLogger.debugIn[0].args[0] != "debug" {
		t.Errorf("Expected first argument to be 'debug', got '%v'", mockLogger.debugIn[0].args[0])
	}

	if mockLogger.debugIn[0].args[1] != 42 {
		t.Errorf("Expected second argument to be 42, got '%v'", mockLogger.debugIn[0].args[1])
	}

	if mockLogger.infoIn != nil {
		t.Errorf("Expected Info not to be called, got %d calls", len(mockLogger.infoIn))
	}

	if mockLogger.warnIn != nil {
		t.Errorf("Expected Warn not to be called, got %d calls", len(mockLogger.warnIn))
	}

	if mockLogger.errorIn != nil {
		t.Errorf("Expected Error not to be called, got %d calls", len(mockLogger.errorIn))
	}

	if debugCounter.inc != 1 {
		t.Errorf("Expected counter to be incremented once, got %f", debugCounter.inc)
	}

	if infoCounter.inc != 0 {
		t.Errorf("Expected info counter to be 0, got %f", infoCounter.inc)
	}

	if warnCounter.inc != 0 {
		t.Errorf("Expected warn counter to be 0, got %f", warnCounter.inc)
	}

	if errorCounter.inc != 0 {
		t.Errorf("Expected error counter to be 0, got %f", errorCounter.inc)
	}
}

func TestPrometheusLogger_Info(t *testing.T) {
	debugCounter := new(mockPrometheusCounter)
	infoCounter := new(mockPrometheusCounter)
	warnCounter := new(mockPrometheusCounter)
	errorCounter := new(mockPrometheusCounter)
	mockLogger := new(mockLogger)

	logger := decorator.NewPrometheusLogger(mockLogger, decorator.PrometheusLoggerCounter{
		Debug: debugCounter,
		Info:  infoCounter,
		Warn:  warnCounter,
		Error: errorCounter,
	})

	logger.Info(context.Background(), "info message", "info", 75)

	if len(mockLogger.infoIn) != 1 {
		t.Fatalf("Expected Info to be called once, got %d", len(mockLogger.infoIn))
	}

	if mockLogger.infoIn[0].msg != "info message" {
		t.Errorf("Expected message to be 'info message', got '%s'", mockLogger.infoIn[0].msg)
	}

	if len(mockLogger.infoIn[0].args) != 2 {
		t.Fatalf("Expected 2 arguments, got %d", len(mockLogger.infoIn[0].args))
	}

	if mockLogger.infoIn[0].args[0] != "info" {
		t.Errorf("Expected first argument to be 'info', got '%v'", mockLogger.infoIn[0].args[0])
	}

	if mockLogger.infoIn[0].args[1] != 75 {
		t.Errorf("Expected second argument to be 75, got '%v'", mockLogger.infoIn[0].args[1])
	}

	if mockLogger.debugIn != nil {
		t.Errorf("Expected Debug not to be called, got %d calls", len(mockLogger.debugIn))
	}

	if mockLogger.warnIn != nil {
		t.Errorf("Expected Warn not to be called, got %d calls", len(mockLogger.warnIn))
	}

	if mockLogger.errorIn != nil {
		t.Errorf("Expected Error not to be called, got %d calls", len(mockLogger.errorIn))
	}

	if debugCounter.inc != 0 {
		t.Errorf("Expected debug counter to be 0, got %f", debugCounter.inc)
	}

	if infoCounter.inc != 1 {
		t.Errorf("Expected info counter to be incremented once, got %f", infoCounter.inc)
	}

	if warnCounter.inc != 0 {
		t.Errorf("Expected warn counter to be 0, got %f", warnCounter.inc)
	}

	if errorCounter.inc != 0 {
		t.Errorf("Expected error counter to be 0, got %f", errorCounter.inc)
	}
}

func TestPrometheusLogger_Warn(t *testing.T) {
	debugCounter := new(mockPrometheusCounter)
	infoCounter := new(mockPrometheusCounter)
	warnCounter := new(mockPrometheusCounter)
	errorCounter := new(mockPrometheusCounter)
	mockLogger := new(mockLogger)

	logger := decorator.NewPrometheusLogger(mockLogger, decorator.PrometheusLoggerCounter{
		Debug: debugCounter,
		Info:  infoCounter,
		Warn:  warnCounter,
		Error: errorCounter,
	})

	logger.Warn(context.Background(), "warn message", "warn", 85)

	if len(mockLogger.warnIn) != 1 {
		t.Fatalf("Expected Warn to be called once, got %d", len(mockLogger.warnIn))
	}

	if mockLogger.warnIn[0].msg != "warn message" {
		t.Errorf("Expected message to be 'warn message', got '%s'", mockLogger.warnIn[0].msg)
	}

	if len(mockLogger.warnIn[0].args) != 2 {
		t.Fatalf("Expected 2 arguments, got %d", len(mockLogger.warnIn[0].args))
	}

	if mockLogger.warnIn[0].args[0] != "warn" {
		t.Errorf("Expected first argument to be 'warn', got '%v'", mockLogger.warnIn[0].args[0])
	}

	if mockLogger.warnIn[0].args[1] != 85 {
		t.Errorf("Expected second argument to be 85, got '%v'", mockLogger.warnIn[0].args[1])
	}

	if mockLogger.debugIn != nil {
		t.Errorf("Expected Debug not to be called, got %d calls", len(mockLogger.debugIn))
	}

	if mockLogger.infoIn != nil {
		t.Errorf("Expected Info not to be called, got %d calls", len(mockLogger.infoIn))
	}

	if mockLogger.errorIn != nil {
		t.Errorf("Expected Error not to be called, got %d calls", len(mockLogger.errorIn))
	}

	if debugCounter.inc != 0 {
		t.Errorf("Expected debug counter to be 0, got %f", debugCounter.inc)
	}

	if infoCounter.inc != 0 {
		t.Errorf("Expected info counter to be 0, got %f", infoCounter.inc)
	}

	if warnCounter.inc != 1 {
		t.Errorf("Expected warn counter to be incremented once, got %f", warnCounter.inc)
	}

	if errorCounter.inc != 0 {
		t.Errorf("Expected error counter to be 0, got %f", errorCounter.inc)
	}
}

func TestPrometheusLogger_Error(t *testing.T) {
	debugCounter := new(mockPrometheusCounter)
	infoCounter := new(mockPrometheusCounter)
	warnCounter := new(mockPrometheusCounter)
	errorCounter := new(mockPrometheusCounter)
	mockLogger := new(mockLogger)

	logger := decorator.NewPrometheusLogger(mockLogger, decorator.PrometheusLoggerCounter{
		Debug: debugCounter,
		Info:  infoCounter,
		Warn:  warnCounter,
		Error: errorCounter,
	})

	logger.Error(context.Background(), context.Canceled, "error", 90)

	if len(mockLogger.errorIn) != 1 {
		t.Fatalf("Expected Error to be called once, got %d", len(mockLogger.errorIn))
	}

	if mockLogger.errorIn[0].err != context.Canceled {
		t.Errorf("Expected error to be context.Canceled, got '%v'", mockLogger.errorIn[0].err)
	}

	if len(mockLogger.errorIn[0].args) != 2 {
		t.Fatalf("Expected 2 arguments, got %d", len(mockLogger.errorIn[0].args))
	}

	if mockLogger.errorIn[0].args[0] != "error" {
		t.Errorf("Expected first argument to be 'error', got '%v'", mockLogger.errorIn[0].args[0])
	}

	if mockLogger.errorIn[0].args[1] != 90 {
		t.Errorf("Expected second argument to be 90, got '%v'", mockLogger.errorIn[0].args[1])
	}

	if mockLogger.debugIn != nil {
		t.Errorf("Expected Debug not to be called, got %d calls", len(mockLogger.debugIn))
	}

	if mockLogger.infoIn != nil {
		t.Errorf("Expected Info not to be called, got %d calls", len(mockLogger.infoIn))
	}

	if mockLogger.warnIn != nil {
		t.Errorf("Expected Warn not to be called, got %d calls", len(mockLogger.warnIn))
	}

	if debugCounter.inc != 0 {
		t.Errorf("Expected debug counter to be 0, got %f", debugCounter.inc)
	}

	if infoCounter.inc != 0 {
		t.Errorf("Expected info counter to be 0, got %f", infoCounter.inc)
	}

	if warnCounter.inc != 0 {
		t.Errorf("Expected warn counter to be 0, got %f", warnCounter.inc)
	}

	if errorCounter.inc != 1 {
		t.Errorf("Expected error counter to be incremented once, got %f", errorCounter.inc)
	}
}

func TestPrometheusLogger_WithoutCounter(t *testing.T) {
	mockLogger := new(mockLogger)

	logger := decorator.NewPrometheusLogger(mockLogger, decorator.PrometheusLoggerCounter{})

	logger.Debug(context.Background(), "debug message")
	logger.Info(context.Background(), "info message")
	logger.Warn(context.Background(), "warn message")
	logger.Error(context.Background(), context.Canceled)

	if len(mockLogger.debugIn) != 1 {
		t.Fatalf("Expected Debug to be called once, got %d", len(mockLogger.debugIn))
	}

	if len(mockLogger.infoIn) != 1 {
		t.Fatalf("Expected Info to be called once, got %d", len(mockLogger.infoIn))
	}

	if len(mockLogger.warnIn) != 1 {
		t.Fatalf("Expected Warn to be called once, got %d", len(mockLogger.warnIn))
	}

	if len(mockLogger.errorIn) != 1 {
		t.Fatalf("Expected Error to be called once, got %d", len(mockLogger.errorIn))
	}
}
