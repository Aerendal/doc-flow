package logger

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
)

func TestLogLevels(t *testing.T) {
	var buf bytes.Buffer
	log := New(&buf, LevelInfo, FormatText)

	log.Debug("should not appear")
	if buf.Len() > 0 {
		t.Error("debug message should not appear at info level")
	}

	log.Info("info message")
	if !strings.Contains(buf.String(), "[INFO]") {
		t.Error("info message should contain [INFO]")
	}
	if !strings.Contains(buf.String(), "info message") {
		t.Error("info message should contain the message text")
	}

	buf.Reset()
	log.Warn("warning %d", 42)
	if !strings.Contains(buf.String(), "[WARN]") {
		t.Error("warn message should contain [WARN]")
	}
	if !strings.Contains(buf.String(), "warning 42") {
		t.Error("warn message should contain formatted text")
	}

	buf.Reset()
	log.Error("error occurred")
	if !strings.Contains(buf.String(), "[ERROR]") {
		t.Error("error message should contain [ERROR]")
	}
}

func TestLogJSON(t *testing.T) {
	var buf bytes.Buffer
	log := New(&buf, LevelDebug, FormatJSON)

	log.Info("test message")

	var entry jsonEntry
	if err := json.Unmarshal(buf.Bytes(), &entry); err != nil {
		t.Fatalf("failed to parse JSON log: %v", err)
	}

	if entry.Level != "INFO" {
		t.Errorf("expected level INFO, got %q", entry.Level)
	}
	if entry.Message != "test message" {
		t.Errorf("expected message 'test message', got %q", entry.Message)
	}
}

func TestLogWithFields(t *testing.T) {
	var buf bytes.Buffer
	log := New(&buf, LevelDebug, FormatText)

	log.WithDocID("my_doc").WithPath("/docs/test.md").Info("processing")
	output := buf.String()

	if !strings.Contains(output, "doc_id=my_doc") {
		t.Error("expected doc_id field in output")
	}
	if !strings.Contains(output, "path=/docs/test.md") {
		t.Error("expected path field in output")
	}
}

func TestLogWithFieldsJSON(t *testing.T) {
	var buf bytes.Buffer
	log := New(&buf, LevelDebug, FormatJSON)

	log.WithDocID("test_doc").Info("check")

	var entry jsonEntry
	if err := json.Unmarshal(buf.Bytes(), &entry); err != nil {
		t.Fatalf("failed to parse JSON: %v", err)
	}

	if entry.Fields["doc_id"] != "test_doc" {
		t.Errorf("expected doc_id=test_doc, got %q", entry.Fields["doc_id"])
	}
}

func TestParseLevel(t *testing.T) {
	tests := map[string]Level{
		"debug":   LevelDebug,
		"info":    LevelInfo,
		"warn":    LevelWarn,
		"error":   LevelError,
		"unknown": LevelInfo,
	}

	for input, expected := range tests {
		got := ParseLevel(input)
		if got != expected {
			t.Errorf("ParseLevel(%q) = %v, want %v", input, got, expected)
		}
	}
}
