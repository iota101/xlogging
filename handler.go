package xlogging

import (
	"context"
	"log/slog"
)

// contextHandler wraps a slog.Handler to extract values from context.
type contextHandler struct {
	inner       slog.Handler
	contextKeys []ContextKey
}

// newContextHandler creates a new contextHandler wrapping the given handler.
func newContextHandler(inner slog.Handler, keys []ContextKey) *contextHandler {
	return &contextHandler{
		inner:       inner,
		contextKeys: keys,
	}
}

// Enabled reports whether the handler handles records at the given level.
func (h *contextHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return h.inner.Enabled(ctx, level)
}

// Handle handles the record, extracting context values and adding them as attributes.
func (h *contextHandler) Handle(ctx context.Context, r slog.Record) error {
	if ctx != nil && len(h.contextKeys) > 0 {
		attrs := make([]slog.Attr, 0, len(h.contextKeys))
		for _, key := range h.contextKeys {
			if v := ctx.Value(key); v != nil {
				if s, ok := v.(string); ok && s != "" {
					attrs = append(attrs, slog.String(string(key), s))
				}
			}
		}
		if len(attrs) > 0 {
			r = r.Clone()
			r.AddAttrs(attrs...)
		}
	}
	return h.inner.Handle(ctx, r)
}

// WithAttrs returns a new handler with the given attributes.
func (h *contextHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &contextHandler{
		inner:       h.inner.WithAttrs(attrs),
		contextKeys: h.contextKeys,
	}
}

// WithGroup returns a new handler with the given group name.
func (h *contextHandler) WithGroup(name string) slog.Handler {
	return &contextHandler{
		inner:       h.inner.WithGroup(name),
		contextKeys: h.contextKeys,
	}
}
