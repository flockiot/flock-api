package logging

import (
	"bytes"
	"log/slog"
	"strings"
	"testing"
)

func TestSetupTextFormat(t *testing.T) {
	var buf bytes.Buffer
	if err := Setup("info", "text", &buf); err != nil {
		t.Fatal(err)
	}
	slog.Info("hello")
	if !strings.Contains(buf.String(), "hello") {
		t.Fatalf("expected 'hello' in output, got: %s", buf.String())
	}
}

func TestSetupJSONFormat(t *testing.T) {
	var buf bytes.Buffer
	if err := Setup("info", "json", &buf); err != nil {
		t.Fatal(err)
	}
	slog.Info("hello")
	if !strings.Contains(buf.String(), `"msg":"hello"`) {
		t.Fatalf("expected JSON msg field, got: %s", buf.String())
	}
}

func TestSetupDebugLevel(t *testing.T) {
	var buf bytes.Buffer
	if err := Setup("debug", "text", &buf); err != nil {
		t.Fatal(err)
	}
	slog.Debug("debug-msg")
	if !strings.Contains(buf.String(), "debug-msg") {
		t.Fatalf("expected debug message in output, got: %s", buf.String())
	}
}

func TestSetupWarnLevelFiltersInfo(t *testing.T) {
	var buf bytes.Buffer
	if err := Setup("warn", "text", &buf); err != nil {
		t.Fatal(err)
	}
	slog.Info("should-not-appear")
	if strings.Contains(buf.String(), "should-not-appear") {
		t.Fatal("info message should be filtered at warn level")
	}
}

func TestSetupInvalidLevel(t *testing.T) {
	var buf bytes.Buffer
	if err := Setup("bogus", "text", &buf); err == nil {
		t.Fatal("expected error for invalid level")
	}
}

func TestSetupInvalidFormat(t *testing.T) {
	var buf bytes.Buffer
	if err := Setup("info", "xml", &buf); err == nil {
		t.Fatal("expected error for invalid format")
	}
}

func TestParseLevelCaseInsensitive(t *testing.T) {
	for _, input := range []string{"DEBUG", "Debug", "dEbUg"} {
		lvl, err := parseLevel(input)
		if err != nil {
			t.Fatalf("parseLevel(%q) returned error: %v", input, err)
		}
		if lvl != slog.LevelDebug {
			t.Fatalf("parseLevel(%q) = %v, want Debug", input, lvl)
		}
	}
}
