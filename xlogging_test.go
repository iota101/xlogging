package xlogging

import (
	"bytes"
	"context"
	"encoding/json"
	"strings"
	"testing"
)

func TestParseLevel(t *testing.T) {
	tests := []struct {
		input    string
		expected Level
	}{
		{"debug", LevelDebug},
		{"DEBUG", LevelDebug},
		{"info", LevelInfo},
		{"INFO", LevelInfo},
		{"warn", LevelWarn},
		{"WARN", LevelWarn},
		{"warning", LevelWarn},
		{"error", LevelError},
		{"ERROR", LevelError},
		{"unknown", LevelInfo}, // default
		{"", LevelInfo},        // default
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := ParseLevel(tt.input)
			if result != tt.expected {
				t.Errorf("ParseLevel(%q) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestDefaultLevelForEnv(t *testing.T) {
	tests := []struct {
		env      Env
		expected Level
	}{
		{EnvProduction, LevelInfo},
		{EnvStaging, LevelDebug},
		{EnvDevelopment, LevelDebug},
	}

	for _, tt := range tests {
		t.Run(string(tt.env), func(t *testing.T) {
			result := defaultLevelForEnv(tt.env)
			if result != tt.expected {
				t.Errorf("defaultLevelForEnv(%q) = %v, want %v", tt.env, result, tt.expected)
			}
		})
	}
}

func TestLoggerBasic(t *testing.T) {
	var buf bytes.Buffer
	log := New(
		WithOutput(&buf),
		WithLevel(LevelDebug),
		WithEnv(EnvProduction), // JSON format
	)

	log.Info("test message", "key", "value")

	var entry map[string]any
	if err := json.Unmarshal(buf.Bytes(), &entry); err != nil {
		t.Fatalf("failed to parse JSON: %v", err)
	}

	if entry["msg"] != "test message" {
		t.Errorf("msg = %v, want %q", entry["msg"], "test message")
	}
	if entry["key"] != "value" {
		t.Errorf("key = %v, want %q", entry["key"], "value")
	}
	if entry["level"] != "INFO" {
		t.Errorf("level = %v, want %q", entry["level"], "INFO")
	}
}

func TestLoggerLevels(t *testing.T) {
	var buf bytes.Buffer
	log := New(
		WithOutput(&buf),
		WithLevel(LevelWarn),
		WithEnv(EnvProduction),
	)

	log.Debug("debug message")
	log.Info("info message")
	log.Warn("warn message")
	log.Error("error message")

	output := buf.String()
	if strings.Contains(output, "debug message") {
		t.Error("debug message should be filtered out")
	}
	if strings.Contains(output, "info message") {
		t.Error("info message should be filtered out")
	}
	if !strings.Contains(output, "warn message") {
		t.Error("warn message should be present")
	}
	if !strings.Contains(output, "error message") {
		t.Error("error message should be present")
	}
}

func TestLoggerWith(t *testing.T) {
	var buf bytes.Buffer
	log := New(
		WithOutput(&buf),
		WithLevel(LevelInfo),
		WithEnv(EnvProduction),
	)

	childLog := log.With("service", "api")
	childLog.Info("child message")

	var entry map[string]any
	if err := json.Unmarshal(buf.Bytes(), &entry); err != nil {
		t.Fatalf("failed to parse JSON: %v", err)
	}

	if entry["service"] != "api" {
		t.Errorf("service = %v, want %q", entry["service"], "api")
	}
}

func TestLoggerWithGroup(t *testing.T) {
	var buf bytes.Buffer
	log := New(
		WithOutput(&buf),
		WithLevel(LevelInfo),
		WithEnv(EnvProduction),
	)

	groupLog := log.WithGroup("request")
	groupLog.Info("grouped message", "method", "GET")

	var entry map[string]any
	if err := json.Unmarshal(buf.Bytes(), &entry); err != nil {
		t.Fatalf("failed to parse JSON: %v", err)
	}

	request, ok := entry["request"].(map[string]any)
	if !ok {
		t.Fatal("request group not found")
	}
	if request["method"] != "GET" {
		t.Errorf("request.method = %v, want %q", request["method"], "GET")
	}
}

func TestContextHandler(t *testing.T) {
	var buf bytes.Buffer
	log := New(
		WithOutput(&buf),
		WithLevel(LevelInfo),
		WithEnv(EnvProduction),
		WithContextKeys(KeyRequestID, KeyTraceID),
	)

	ctx := context.Background()
	ctx = WithRequestID(ctx, "req-123")
	ctx = WithTraceID(ctx, "trace-456")

	log.InfoContext(ctx, "context message")

	var entry map[string]any
	if err := json.Unmarshal(buf.Bytes(), &entry); err != nil {
		t.Fatalf("failed to parse JSON: %v", err)
	}

	if entry["request_id"] != "req-123" {
		t.Errorf("request_id = %v, want %q", entry["request_id"], "req-123")
	}
	if entry["trace_id"] != "trace-456" {
		t.Errorf("trace_id = %v, want %q", entry["trace_id"], "trace-456")
	}
}

func TestContextHelpers(t *testing.T) {
	ctx := context.Background()
	ctx = WithRequestID(ctx, "req-123")
	ctx = WithTraceID(ctx, "trace-456")
	ctx = WithSpanID(ctx, "span-789")
	ctx = WithUserID(ctx, "user-abc")

	if v, ok := GetRequestID(ctx); !ok || v != "req-123" {
		t.Errorf("GetRequestID() = %v, %v, want %q, true", v, ok, "req-123")
	}
	if v, ok := GetTraceID(ctx); !ok || v != "trace-456" {
		t.Errorf("GetTraceID() = %v, %v, want %q, true", v, ok, "trace-456")
	}
	if v, ok := GetSpanID(ctx); !ok || v != "span-789" {
		t.Errorf("GetSpanID() = %v, %v, want %q, true", v, ok, "span-789")
	}
	if v, ok := GetUserID(ctx); !ok || v != "user-abc" {
		t.Errorf("GetUserID() = %v, %v, want %q, true", v, ok, "user-abc")
	}
}

func TestColorHandler(t *testing.T) {
	var buf bytes.Buffer
	log := New(
		WithOutput(&buf),
		WithLevel(LevelDebug),
		WithEnv(EnvDevelopment),
		WithColor(true),
	)

	log.Debug("debug message")
	log.Info("info message")
	log.Warn("warn message")
	log.Error("error message")

	output := buf.String()
	if !strings.Contains(output, "debug message") {
		t.Error("debug message should be present")
	}
	if !strings.Contains(output, "info message") {
		t.Error("info message should be present")
	}
	if !strings.Contains(output, "warn message") {
		t.Error("warn message should be present")
	}
	if !strings.Contains(output, "error message") {
		t.Error("error message should be present")
	}

	// Check that ANSI codes are present
	if !strings.Contains(output, "\033[") {
		t.Error("ANSI color codes should be present")
	}
}

func TestTextHandler(t *testing.T) {
	var buf bytes.Buffer
	log := New(
		WithOutput(&buf),
		WithLevel(LevelInfo),
		WithEnv(EnvDevelopment),
		WithColor(false),
	)

	log.Info("text message", "key", "value")

	output := buf.String()
	if !strings.Contains(output, "text message") {
		t.Error("message should be present")
	}
	if !strings.Contains(output, "key=value") {
		t.Error("attribute should be present")
	}
}

func TestTestLogger(t *testing.T) {
	log := NewTestLogger()

	log.Debug("debug")
	log.Info("info")
	log.Warn("warn")
	log.Error("error")

	if log.Len() != 4 {
		t.Errorf("Len() = %d, want 4", log.Len())
	}

	if !log.HasEntry(LevelDebug, "debug") {
		t.Error("should have debug entry")
	}
	if !log.HasEntry(LevelInfo, "info") {
		t.Error("should have info entry")
	}
	if !log.HasEntry(LevelWarn, "warn") {
		t.Error("should have warn entry")
	}
	if !log.HasEntry(LevelError, "error") {
		t.Error("should have error entry")
	}

	if log.Count(LevelInfo) != 1 {
		t.Errorf("Count(LevelInfo) = %d, want 1", log.Count(LevelInfo))
	}
}

func TestTestLoggerWithAttrs(t *testing.T) {
	log := NewTestLogger()

	log.Info("message", "key", "value")

	if !log.HasEntryWithAttr(LevelInfo, "message", "key", "value") {
		t.Error("should have entry with attribute")
	}
}

func TestTestLoggerWith(t *testing.T) {
	log := NewTestLogger()

	childLog := log.With("service", "api")
	childLog.Info("child message")

	if !log.HasEntryWithAttr(LevelInfo, "child", "service", "api") {
		t.Error("should have entry with pre-set attribute")
	}
}

func TestTestLoggerClear(t *testing.T) {
	log := NewTestLogger()

	log.Info("message")
	if log.Len() != 1 {
		t.Errorf("Len() = %d, want 1", log.Len())
	}

	log.Clear()
	if log.Len() != 0 {
		t.Errorf("Len() after Clear() = %d, want 0", log.Len())
	}
}

func TestDiscard(t *testing.T) {
	log := Discard()

	// Should not panic
	log.Debug("debug")
	log.Info("info")
	log.Warn("warn")
	log.Error("error")
	log.DebugContext(context.Background(), "debug")
	log.InfoContext(context.Background(), "info")
	log.WarnContext(context.Background(), "warn")
	log.ErrorContext(context.Background(), "error")
}

func TestConfigOptions(t *testing.T) {
	var buf bytes.Buffer
	enabled := true

	log := New(
		WithOutput(&buf),
		WithEnv(EnvStaging),
		WithLevel(LevelDebug),
		WithSource(true),
		WithColor(false),
		WithContextKeys(KeyRequestID),
	)

	_ = enabled // silence unused variable warning

	if log.Handler() == nil {
		t.Error("Handler() should not be nil")
	}
}
