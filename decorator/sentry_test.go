package decorator_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/getsentry/sentry-go"
	"github.com/yuppyweb/cakelog/decorator"
)

type mockSentryTransport struct {
	Event *sentry.Event
}

func (m *mockSentryTransport) Flush(time.Duration) bool {
	return true
}

func (m *mockSentryTransport) FlushWithContext(context.Context) bool {
	return true
}

func (m *mockSentryTransport) Configure(sentry.ClientOptions) {}

func (m *mockSentryTransport) SendEvent(event *sentry.Event) {
	m.Event = event
}

func (m *mockSentryTransport) Close() {}

var _ sentry.Transport = (*mockSentryTransport)(nil)

func TestSentryLogger_Debug(t *testing.T) {
	t.Parallel()

	mockLogger := new(mockLogger)
	mockTransport := new(mockSentryTransport)

	client, err := sentry.NewClient(sentry.ClientOptions{
		Dsn:       "https://examplePublicKey@o0.ingest.sentry.io/0",
		Transport: mockTransport,
	})
	if err != nil {
		t.Fatalf("Failed to create Sentry client: %v", err)
	}

	debugScope := sentry.NewScope()
	debugScope.SetLevel(sentry.LevelDebug)

	hub := decorator.SentryLoggerHub{
		Debug: sentry.NewHub(client, debugScope),
	}

	msg := "debug message"

	logger := decorator.NewSentryLogger(mockLogger, hub)

	logger.Debug(context.Background(), msg, "debug", 42)

	if len(mockLogger.debugIn) != 1 {
		t.Fatalf("Expected Debug to be called once, got %d", len(mockLogger.debugIn))
	}

	if mockLogger.debugIn[0].msg != msg {
		t.Errorf("Expected message to be '%s', got '%s'", msg, mockLogger.debugIn[0].msg)
	}

	if len(mockLogger.debugIn[0].args) != 3 {
		t.Fatalf("Expected 3 arguments, got %d", len(mockLogger.debugIn[0].args))
	}

	if mockLogger.debugIn[0].args[0] != "debug" {
		t.Errorf("Expected first argument to be 'debug', got '%v'", mockLogger.debugIn[0].args[0])
	}

	if mockLogger.debugIn[0].args[1] != 42 {
		t.Errorf("Expected second argument to be 42, got '%v'", mockLogger.debugIn[0].args[1])
	}

	sentryArgs, ok := mockLogger.debugIn[0].args[2].(map[string]*sentry.EventID)
	if !ok {
		t.Fatalf(
			"Expected third argument to be a map[string]*sentry.EventID, got '%T'",
			mockLogger.debugIn[0].args[2],
		)
	}

	eventID, ok := sentryArgs[logger.SentryEventIDKey]
	if !ok {
		t.Fatalf(
			"Expected Sentry event ID key '%s' to be present in log arguments",
			logger.SentryEventIDKey,
		)
	}

	if eventID != &mockTransport.Event.EventID {
		t.Errorf(
			"Expected Sentry event ID '%s', got '%s'",
			mockTransport.Event.EventID,
			string(*eventID),
		)
	}

	if mockTransport.Event.Message != msg {
		t.Errorf(
			"Expected Sentry event message to be '%s', got '%s'",
			msg,
			mockTransport.Event.Message,
		)
	}

	if mockTransport.Event.Level != sentry.LevelDebug {
		t.Errorf("Expected Sentry event level to be 'debug', got '%s'", mockTransport.Event.Level)
	}
}

func TestSentryLogger_Info(t *testing.T) {
	t.Parallel()

	mockLogger := new(mockLogger)
	mockTransport := new(mockSentryTransport)

	client, err := sentry.NewClient(sentry.ClientOptions{
		Dsn:       "https://examplePublicKey@o0.ingest.sentry.io/0",
		Transport: mockTransport,
	})
	if err != nil {
		t.Fatalf("Failed to create Sentry client: %v", err)
	}

	infoScope := sentry.NewScope()
	infoScope.SetLevel(sentry.LevelInfo)

	hub := decorator.SentryLoggerHub{
		Info: sentry.NewHub(client, infoScope),
	}

	msg := "info message"

	logger := decorator.NewSentryLogger(mockLogger, hub)

	logger.Info(context.Background(), msg, "info", 57)

	if len(mockLogger.infoIn) != 1 {
		t.Fatalf("Expected Info to be called once, got %d", len(mockLogger.infoIn))
	}

	if mockLogger.infoIn[0].msg != msg {
		t.Errorf("Expected message to be '%s', got '%s'", msg, mockLogger.infoIn[0].msg)
	}

	if len(mockLogger.infoIn[0].args) != 3 {
		t.Fatalf("Expected 3 arguments, got %d", len(mockLogger.infoIn[0].args))
	}

	if mockLogger.infoIn[0].args[0] != "info" {
		t.Errorf("Expected first argument to be 'info', got '%v'", mockLogger.infoIn[0].args[0])
	}

	if mockLogger.infoIn[0].args[1] != 57 {
		t.Errorf("Expected second argument to be 57, got '%v'", mockLogger.infoIn[0].args[1])
	}

	sentryArgs, ok := mockLogger.infoIn[0].args[2].(map[string]*sentry.EventID)
	if !ok {
		t.Fatalf(
			"Expected third argument to be a map[string]*sentry.EventID, got '%T'",
			mockLogger.infoIn[0].args[2],
		)
	}

	eventID, ok := sentryArgs[logger.SentryEventIDKey]
	if !ok {
		t.Fatalf(
			"Expected Sentry event ID key '%s' to be present in log arguments",
			logger.SentryEventIDKey,
		)
	}

	if eventID != &mockTransport.Event.EventID {
		t.Errorf(
			"Expected Sentry event ID '%s', got '%s'",
			mockTransport.Event.EventID,
			string(*eventID),
		)
	}

	if mockTransport.Event.Message != msg {
		t.Errorf(
			"Expected Sentry event message to be '%s', got '%s'",
			msg,
			mockTransport.Event.Message,
		)
	}

	if mockTransport.Event.Level != sentry.LevelInfo {
		t.Errorf("Expected Sentry event level to be 'info', got '%s'", mockTransport.Event.Level)
	}
}

func TestSentryLogger_Warn(t *testing.T) {
	t.Parallel()

	mockLogger := new(mockLogger)
	mockTransport := new(mockSentryTransport)

	client, err := sentry.NewClient(sentry.ClientOptions{
		Dsn:       "https://examplePublicKey@o0.ingest.sentry.io/0",
		Transport: mockTransport,
	})
	if err != nil {
		t.Fatalf("Failed to create Sentry client: %v", err)
	}

	warnScope := sentry.NewScope()
	warnScope.SetLevel(sentry.LevelWarning)

	hub := decorator.SentryLoggerHub{
		Warn: sentry.NewHub(client, warnScope),
	}

	msg := "warn message"

	logger := decorator.NewSentryLogger(mockLogger, hub)

	logger.Warn(context.Background(), msg, "warn", 69)

	if len(mockLogger.warnIn) != 1 {
		t.Fatalf("Expected Warn to be called once, got %d", len(mockLogger.warnIn))
	}

	if mockLogger.warnIn[0].msg != msg {
		t.Errorf("Expected message to be '%s', got '%s'", msg, mockLogger.warnIn[0].msg)
	}

	if len(mockLogger.warnIn[0].args) != 3 {
		t.Fatalf("Expected 3 arguments, got %d", len(mockLogger.warnIn[0].args))
	}

	if mockLogger.warnIn[0].args[0] != "warn" {
		t.Errorf("Expected first argument to be 'warn', got '%v'", mockLogger.warnIn[0].args[0])
	}

	if mockLogger.warnIn[0].args[1] != 69 {
		t.Errorf("Expected second argument to be 69, got '%v'", mockLogger.warnIn[0].args[1])
	}

	sentryArgs, ok := mockLogger.warnIn[0].args[2].(map[string]*sentry.EventID)
	if !ok {
		t.Fatalf(
			"Expected third argument to be a map[string]*sentry.EventID, got '%T'",
			mockLogger.warnIn[0].args[2],
		)
	}

	eventID, ok := sentryArgs[logger.SentryEventIDKey]
	if !ok {
		t.Fatalf(
			"Expected Sentry event ID key '%s' to be present in log arguments",
			logger.SentryEventIDKey,
		)
	}

	if eventID != &mockTransport.Event.EventID {
		t.Errorf(
			"Expected Sentry event ID '%s', got '%s'",
			mockTransport.Event.EventID,
			string(*eventID),
		)
	}

	if mockTransport.Event.Message != msg {
		t.Errorf(
			"Expected Sentry event message to be '%s', got '%s'",
			msg,
			mockTransport.Event.Message,
		)
	}

	if mockTransport.Event.Level != sentry.LevelWarning {
		t.Errorf("Expected Sentry event level to be 'warning', got '%s'", mockTransport.Event.Level)
	}
}

func TestSentryLogger_Error(t *testing.T) {
	t.Parallel()

	mockLogger := new(mockLogger)
	mockTransport := new(mockSentryTransport)

	client, err := sentry.NewClient(sentry.ClientOptions{
		Dsn:       "https://examplePublicKey@o0.ingest.sentry.io/0",
		Transport: mockTransport,
	})
	if err != nil {
		t.Fatalf("Failed to create Sentry client: %v", err)
	}

	errorScope := sentry.NewScope()
	errorScope.SetLevel(sentry.LevelError)

	hub := decorator.SentryLoggerHub{
		Error: sentry.NewHub(client, errorScope),
	}

	errMsg := "test error"
	expectedErr := errors.New(errMsg)

	logger := decorator.NewSentryLogger(mockLogger, hub)

	logger.Error(context.Background(), expectedErr, "error", 85)

	if len(mockLogger.errorIn) != 1 {
		t.Fatalf("Expected Error to be called once, got %d", len(mockLogger.errorIn))
	}

	if !errors.Is(mockLogger.errorIn[0].err, expectedErr) {
		t.Errorf("Expected error to be '%v', got '%v'", expectedErr, mockLogger.errorIn[0].err)
	}

	if len(mockLogger.errorIn[0].args) != 3 {
		t.Fatalf("Expected 3 arguments, got %d", len(mockLogger.errorIn[0].args))
	}

	if mockLogger.errorIn[0].args[0] != "error" {
		t.Errorf("Expected first argument to be 'error', got '%v'", mockLogger.errorIn[0].args[0])
	}

	if mockLogger.errorIn[0].args[1] != 85 {
		t.Errorf("Expected second argument to be 85, got '%v'", mockLogger.errorIn[0].args[1])
	}

	sentryArgs, ok := mockLogger.errorIn[0].args[2].(map[string]*sentry.EventID)
	if !ok {
		t.Fatalf(
			"Expected third argument to be a map[string]*sentry.EventID, got '%T'",
			mockLogger.errorIn[0].args[2],
		)
	}

	eventID, ok := sentryArgs[logger.SentryEventIDKey]
	if !ok {
		t.Fatalf(
			"Expected Sentry event ID key '%s' to be present in log arguments",
			logger.SentryEventIDKey,
		)
	}

	if eventID != &mockTransport.Event.EventID {
		t.Errorf(
			"Expected Sentry event ID '%s', got '%s'",
			mockTransport.Event.EventID,
			string(*eventID),
		)
	}

	if len(mockTransport.Event.Exception) == 0 {
		t.Fatalf("Expected Sentry event to contain an exception, got none")
	}

	exception := mockTransport.Event.Exception[0]

	if exception.Value != errMsg {
		t.Errorf(
			"Expected Sentry event exception value to be '%s', got '%s'",
			errMsg,
			exception.Value,
		)
	}

	if mockTransport.Event.Level != sentry.LevelError {
		t.Errorf("Expected Sentry event level to be 'error', got '%s'", mockTransport.Event.Level)
	}
}

func TestSentryLogger_LogChainWithHub(t *testing.T) {
	t.Parallel()

	mockLogger := new(mockLogger)
	mockTransport := new(mockSentryTransport)

	client, err := sentry.NewClient(sentry.ClientOptions{
		Dsn:       "https://examplePublicKey@o0.ingest.sentry.io/0",
		Transport: mockTransport,
	})
	if err != nil {
		t.Fatalf("Failed to create Sentry client: %v", err)
	}

	debugScope := sentry.NewScope()
	debugScope.SetLevel(sentry.LevelDebug)

	infoScope := sentry.NewScope()
	infoScope.SetLevel(sentry.LevelInfo)

	warnScope := sentry.NewScope()
	warnScope.SetLevel(sentry.LevelWarning)

	errorScope := sentry.NewScope()
	errorScope.SetLevel(sentry.LevelError)

	hub := decorator.SentryLoggerHub{
		Debug: sentry.NewHub(client, debugScope),
		Info:  sentry.NewHub(client, infoScope),
		Warn:  sentry.NewHub(client, warnScope),
		Error: sentry.NewHub(client, errorScope),
	}

	logger := decorator.NewSentryLogger(mockLogger, hub)

	logger.Debug(context.Background(), "debug message")

	if mockTransport.Event == nil {
		t.Fatalf("Expected Sentry event to be captured, got nil")
	}

	if mockTransport.Event.Level != sentry.LevelDebug {
		t.Errorf("Expected Sentry event level to be 'debug', got '%s'", mockTransport.Event.Level)
	}

	if mockTransport.Event.Message != "debug message" {
		t.Errorf(
			"Expected Sentry event message to be 'debug message', got '%s'",
			mockTransport.Event.Message,
		)
	}

	if len(mockLogger.debugIn[0].args) != 1 {
		t.Fatalf("Expected 1 argument, got %d", len(mockLogger.debugIn[0].args))
	}

	sentryArgs, ok := mockLogger.debugIn[0].args[0].(map[string]*sentry.EventID)
	if !ok {
		t.Fatalf(
			"Expected third argument to be a map[string]*sentry.EventID, got '%T'",
			mockLogger.debugIn[0].args[0],
		)
	}

	eventID, ok := sentryArgs[logger.SentryEventIDKey]
	if !ok {
		t.Fatalf(
			"Expected Sentry event ID key '%s' to be present in log arguments",
			logger.SentryEventIDKey,
		)
	}

	if eventID != &mockTransport.Event.EventID {
		t.Errorf(
			"Expected Sentry event ID '%s', got '%s'",
			mockTransport.Event.EventID,
			string(*eventID),
		)
	}

	logger.Info(context.Background(), "info message")

	if mockTransport.Event == nil {
		t.Fatalf("Expected Sentry event to be captured, got nil")
	}

	if mockTransport.Event.Level != sentry.LevelInfo {
		t.Errorf("Expected Sentry event level to be 'info', got '%s'", mockTransport.Event.Level)
	}

	if mockTransport.Event.Message != "info message" {
		t.Errorf(
			"Expected Sentry event message to be 'info message', got '%s'",
			mockTransport.Event.Message,
		)
	}

	if len(mockLogger.infoIn[0].args) != 1 {
		t.Fatalf("Expected 1 argument, got %d", len(mockLogger.infoIn[0].args))
	}

	sentryArgs, ok = mockLogger.infoIn[0].args[0].(map[string]*sentry.EventID)
	if !ok {
		t.Fatalf(
			"Expected third argument to be a map[string]*sentry.EventID, got '%T'",
			mockLogger.infoIn[0].args[0],
		)
	}

	eventID, ok = sentryArgs[logger.SentryEventIDKey]
	if !ok {
		t.Fatalf(
			"Expected Sentry event ID key '%s' to be present in log arguments",
			logger.SentryEventIDKey,
		)
	}

	if eventID != &mockTransport.Event.EventID {
		t.Errorf(
			"Expected Sentry event ID '%s', got '%s'",
			mockTransport.Event.EventID,
			string(*eventID),
		)
	}

	logger.Warn(context.Background(), "warn message")

	if mockTransport.Event == nil {
		t.Fatalf("Expected Sentry event to be captured, got nil")
	}

	if mockTransport.Event.Level != sentry.LevelWarning {
		t.Errorf("Expected Sentry event level to be 'warning', got '%s'", mockTransport.Event.Level)
	}

	if mockTransport.Event.Message != "warn message" {
		t.Errorf(
			"Expected Sentry event message to be 'warn message', got '%s'",
			mockTransport.Event.Message,
		)
	}

	if len(mockLogger.warnIn[0].args) != 1 {
		t.Fatalf("Expected 1 argument, got %d", len(mockLogger.warnIn[0].args))
	}

	sentryArgs, ok = mockLogger.warnIn[0].args[0].(map[string]*sentry.EventID)
	if !ok {
		t.Fatalf(
			"Expected third argument to be a map[string]*sentry.EventID, got '%T'",
			mockLogger.warnIn[0].args[0],
		)
	}

	eventID, ok = sentryArgs[logger.SentryEventIDKey]
	if !ok {
		t.Fatalf(
			"Expected Sentry event ID key '%s' to be present in log arguments",
			logger.SentryEventIDKey,
		)
	}

	if eventID != &mockTransport.Event.EventID {
		t.Errorf(
			"Expected Sentry event ID '%s', got '%s'",
			mockTransport.Event.EventID,
			string(*eventID),
		)
	}

	errMsg := "error message"
	logger.Error(context.Background(), errors.New(errMsg))

	if mockTransport.Event == nil {
		t.Fatalf("Expected Sentry event to be captured, got nil")
	}

	if mockTransport.Event.Level != sentry.LevelError {
		t.Errorf("Expected Sentry event level to be 'error', got '%s'", mockTransport.Event.Level)
	}

	if len(mockTransport.Event.Exception) == 0 {
		t.Fatalf("Expected Sentry event to contain an exception, got none")
	}

	exception := mockTransport.Event.Exception[0]

	if exception.Value != errMsg {
		t.Errorf(
			"Expected Sentry event exception value to be '%s', got '%s'",
			errMsg,
			exception.Value,
		)
	}

	if len(mockLogger.errorIn[0].args) != 1 {
		t.Fatalf("Expected 1 argument, got %d", len(mockLogger.errorIn[0].args))
	}

	sentryArgs, ok = mockLogger.errorIn[0].args[0].(map[string]*sentry.EventID)
	if !ok {
		t.Fatalf(
			"Expected third argument to be a map[string]*sentry.EventID, got '%T'",
			mockLogger.errorIn[0].args[0],
		)
	}

	eventID, ok = sentryArgs[logger.SentryEventIDKey]
	if !ok {
		t.Fatalf(
			"Expected Sentry event ID key '%s' to be present in log arguments",
			logger.SentryEventIDKey,
		)
	}

	if eventID != &mockTransport.Event.EventID {
		t.Errorf(
			"Expected Sentry event ID '%s', got '%s'",
			mockTransport.Event.EventID,
			string(*eventID),
		)
	}
}

func TestSentryLogger_WithoutHub(t *testing.T) {
	t.Parallel()

	mockLogger := new(mockLogger)

	logger := decorator.NewSentryLogger(mockLogger, decorator.SentryLoggerHub{})

	expectedErr := errors.New("error message")

	logger.Debug(context.Background(), "debug message")
	logger.Info(context.Background(), "info message")
	logger.Warn(context.Background(), "warn message")
	logger.Error(context.Background(), expectedErr)

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

	if mockLogger.debugIn[0].msg != "debug message" {
		t.Errorf(
			"Expected Debug message to be 'debug message', got '%s'",
			mockLogger.debugIn[0].msg,
		)
	}

	if mockLogger.infoIn[0].msg != "info message" {
		t.Errorf("Expected Info message to be 'info message', got '%s'", mockLogger.infoIn[0].msg)
	}

	if mockLogger.warnIn[0].msg != "warn message" {
		t.Errorf("Expected Warn message to be 'warn message', got '%s'", mockLogger.warnIn[0].msg)
	}

	if !errors.Is(mockLogger.errorIn[0].err, expectedErr) {
		t.Errorf("Expected Error to be '%v', got '%v'", expectedErr, mockLogger.errorIn[0].err)
	}

	if len(mockLogger.debugIn[0].args) != 0 {
		t.Errorf("Expected Debug args to be empty, got %v", mockLogger.debugIn[0].args)
	}

	if len(mockLogger.infoIn[0].args) != 0 {
		t.Errorf("Expected Info args to be empty, got %v", mockLogger.infoIn[0].args)
	}

	if len(mockLogger.warnIn[0].args) != 0 {
		t.Errorf("Expected Warn args to be empty, got %v", mockLogger.warnIn[0].args)
	}

	if len(mockLogger.errorIn[0].args) != 0 {
		t.Errorf("Expected Error args to be empty, got %v", mockLogger.errorIn[0].args)
	}
}

func TestSentryLogger_WithEmptyClient(t *testing.T) {
	t.Parallel()

	mockLogger := new(mockLogger)

	hub := decorator.SentryLoggerHub{
		Debug: sentry.NewHub(nil, nil),
		Info:  sentry.NewHub(nil, nil),
		Warn:  sentry.NewHub(nil, nil),
		Error: sentry.NewHub(nil, nil),
	}

	logger := decorator.NewSentryLogger(mockLogger, hub)

	logger.Debug(context.Background(), "debug message")
	logger.Info(context.Background(), "info message")
	logger.Warn(context.Background(), "warn message")
	logger.Error(context.Background(), errors.New("error message"))

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
