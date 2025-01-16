package log

import (
	"context"
	"fmt"
	"log/slog"
)

var _ slog.Handler = (*multiHandler)(nil)

// newMultiHandler 创建一个新的多处理程序
func newMultiHandler(handlers ...slog.Handler) slog.Handler {
	return &multiHandler{
		handlers: handlers,
	}
}

type multiHandler struct {
	handlers []slog.Handler
}

func (h *multiHandler) Enabled(ctx context.Context, level slog.Level) bool {
	for i := range h.handlers {
		if h.handlers[i].Enabled(ctx, level) {
			return true
		}
	}

	return false
}

func (h *multiHandler) Handle(ctx context.Context, record slog.Record) (err error) {
	for i := range h.handlers {
		if h.handlers[i].Enabled(ctx, record.Level) {
			err = func() error {
				defer func() {
					if v := recover(); v != nil {
						switch v.(type) {
						case error:
							err = v.(error)
						default:
							err = fmt.Errorf("recover from panic: %v", v)
						}
					}
				}()
				return h.handlers[i].Handle(ctx, record.Clone())
			}()
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (h *multiHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	var handlers = make([]slog.Handler, len(h.handlers))
	for i, s := range h.handlers {
		handlers[i] = s.WithAttrs(attrs)
	}
	return newMultiHandler(handlers...)
}

func (h *multiHandler) WithGroup(name string) slog.Handler {
	var handlers = make([]slog.Handler, len(h.handlers))
	for i, s := range h.handlers {
		handlers[i] = s.WithGroup(name)
	}
	return newMultiHandler(handlers...)
}
