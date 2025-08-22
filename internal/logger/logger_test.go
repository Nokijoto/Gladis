package logger

import (
	"fmt"
	"testing"
)

func TestLogLevels(t *testing.T) {
	// Test that log levels are correctly defined
	levels := []LogLevel{
		LevelDebug,
		LevelInfo,
		LevelWarn,
		LevelError,
		LevelFatal,
	}

	expectedLevels := []string{
		"DEBUG",
		"INFO",
		"WARN",
		"ERROR",
		"FATAL",
	}

	for i, level := range levels {
		if string(level) != expectedLevels[i] {
			t.Errorf("Expected level %s, got %s", expectedLevels[i], string(level))
		}
	}
}

func TestSetLevel(t *testing.T) {
	// Save original level
	originalLevel := logLevel
	defer func() {
		logLevel = originalLevel
	}()

	// Test setting level
	SetLevel(LevelError)
	if logLevel != LevelError {
		t.Errorf("Expected log level to be %s, got %s", LevelError, logLevel)
	}

	SetLevel(LevelDebug)
	if logLevel != LevelDebug {
		t.Errorf("Expected log level to be %s, got %s", LevelDebug, logLevel)
	}
}

func TestLogMessage(t *testing.T) {
	// Test LogMessage structure
	msg := LogMessage{
		Level:   LevelInfo,
		Message: "test message",
		Context: "test context",
	}

	if msg.Level != LevelInfo {
		t.Errorf("Expected level %s, got %s", LevelInfo, msg.Level)
	}
	if msg.Message != "test message" {
		t.Errorf("Expected message 'test message', got '%s'", msg.Message)
	}
	if msg.Context != "test context" {
		t.Errorf("Expected context 'test context', got '%s'", msg.Context)
	}
}

func TestTryCatch(t *testing.T) {
	// Test successful execution
	err := TryCatch(func() error {
		return nil
	}, "test context")

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Test error execution
	testError := fmt.Errorf("test error")
	err = TryCatch(func() error {
		return testError
	}, "test context")

	if err != testError {
		t.Errorf("Expected test error, got %v", err)
	}
}

func TestGetLogColor(t *testing.T) {
	tests := []struct {
		level    LogLevel
		expected int
	}{
		{LevelDebug, 0x7289DA},
		{LevelInfo, 0x00FF00},
		{LevelWarn, 0xFFFF00},
		{LevelError, 0xFF0000},
		{LevelFatal, 0xFF0000},
		{"UNKNOWN", 0xFFFFFF},
	}

	for _, test := range tests {
		color := getLogColor(test.level)
		if color != test.expected {
			t.Errorf("Expected color %x for level %s, got %x", test.expected, test.level, color)
		}
	}
}
