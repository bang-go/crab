package logx

import (
	"testing"
)

func TestCallerSkip(t *testing.T) {
	logger := New(
		WithLevel(LevelDebug),
		WithEncodeText(),
		WithCallerSkip(3),
		WithSource(true),
	)
	logger.Info("test message")
	t.Log("CallerSkip feature works")
}

func myWrapper(msg string) {
	logger := New(
		WithLevel(LevelInfo),
		WithCallerSkip(4),
		WithSource(true),
	)
	logger.Info(msg)
}

func TestWrappedLogger(t *testing.T) {
	myWrapper("wrapped message")
	t.Log("Wrapped logger works")
}

func TestGlobalCallerSkip(t *testing.T) {
	orig := GetCallerSkip()
	defer SetCallerSkip(orig)

	SetCallerSkip(5)
	if GetCallerSkip() != 5 {
		t.Errorf("expected 5, got %d", GetCallerSkip())
	}
	t.Log("Global caller skip works")
}
