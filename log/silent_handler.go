package log

import (
	"context"
	"log/slog"
)

var _ Handler = (*silentHandler)(nil)

func newSilentHandler() Handler {
	return &silentHandler{}
}

type silentHandler struct {
}

func (s *silentHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return false
}

func (s *silentHandler) Handle(ctx context.Context, record slog.Record) error {
	return nil
}

func (s *silentHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return nil
}

func (s *silentHandler) WithGroup(name string) slog.Handler {
	return nil
}
