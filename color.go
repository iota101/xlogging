package xlogging

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"sync"
	"time"
)

// ANSI color codes.
const (
	colorReset  = "\033[0m"
	colorRed    = "\033[31m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorBlue   = "\033[34m"
	colorCyan   = "\033[36m"
	colorGray   = "\033[90m"
	colorBold   = "\033[1m"
)

// colorHandler is a slog.Handler that outputs colored text.
type colorHandler struct {
	w           io.Writer
	level       slog.Leveler
	addSource   bool
	attrs       []slog.Attr
	groups      []string
	mu          *sync.Mutex
	contextKeys []ContextKey
}

// colorHandlerOptions configures the colorHandler.
type colorHandlerOptions struct {
	Level       slog.Leveler
	AddSource   bool
	ContextKeys []ContextKey
}

// newColorHandler creates a new colorHandler.
func newColorHandler(w io.Writer, opts *colorHandlerOptions) *colorHandler {
	h := &colorHandler{
		w:  w,
		mu: &sync.Mutex{},
	}
	if opts != nil {
		h.level = opts.Level
		h.addSource = opts.AddSource
		h.contextKeys = opts.ContextKeys
	}
	if h.level == nil {
		h.level = slog.LevelInfo
	}
	return h
}

// Enabled reports whether the handler handles records at the given level.
func (h *colorHandler) Enabled(_ context.Context, level slog.Level) bool {
	return level >= h.level.Level()
}

// Handle handles the record, outputting colored text.
func (h *colorHandler) Handle(ctx context.Context, r slog.Record) error {
	h.mu.Lock()
	defer h.mu.Unlock()

	// Time
	timeStr := r.Time.Format(time.TimeOnly)
	fmt.Fprintf(h.w, "%s%s%s ", colorGray, timeStr, colorReset)

	// Level with color
	levelColor := h.levelColor(r.Level)
	levelStr := h.levelString(r.Level)
	fmt.Fprintf(h.w, "%s%s%-5s%s ", levelColor, colorBold, levelStr, colorReset)

	// Message
	fmt.Fprintf(h.w, "%s%s%s", colorBold, r.Message, colorReset)

	// Context values
	if ctx != nil && len(h.contextKeys) > 0 {
		for _, key := range h.contextKeys {
			if v := ctx.Value(key); v != nil {
				if s, ok := v.(string); ok && s != "" {
					fmt.Fprintf(h.w, " %s%s%s=%s%s%s", colorCyan, string(key), colorReset, colorGray, s, colorReset)
				}
			}
		}
	}

	// Pre-set attributes
	for _, attr := range h.attrs {
		h.writeAttr(attr, h.groups)
	}

	// Record attributes
	r.Attrs(func(a slog.Attr) bool {
		h.writeAttr(a, h.groups)
		return true
	})

	fmt.Fprintln(h.w)
	return nil
}

// writeAttr writes a single attribute with proper formatting.
func (h *colorHandler) writeAttr(a slog.Attr, groups []string) {
	if a.Equal(slog.Attr{}) {
		return
	}

	key := a.Key
	for i := len(groups) - 1; i >= 0; i-- {
		key = groups[i] + "." + key
	}

	if a.Value.Kind() == slog.KindGroup {
		attrs := a.Value.Group()
		newGroups := append(groups, a.Key)
		for _, ga := range attrs {
			h.writeAttr(ga, newGroups)
		}
		return
	}

	fmt.Fprintf(h.w, " %s%s%s=%v", colorCyan, key, colorReset, a.Value.Any())
}

// WithAttrs returns a new handler with the given attributes.
func (h *colorHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	newAttrs := make([]slog.Attr, len(h.attrs)+len(attrs))
	copy(newAttrs, h.attrs)
	copy(newAttrs[len(h.attrs):], attrs)
	return &colorHandler{
		w:           h.w,
		level:       h.level,
		addSource:   h.addSource,
		attrs:       newAttrs,
		groups:      h.groups,
		mu:          h.mu,
		contextKeys: h.contextKeys,
	}
}

// WithGroup returns a new handler with the given group name.
func (h *colorHandler) WithGroup(name string) slog.Handler {
	newGroups := make([]string, len(h.groups)+1)
	copy(newGroups, h.groups)
	newGroups[len(h.groups)] = name
	return &colorHandler{
		w:           h.w,
		level:       h.level,
		addSource:   h.addSource,
		attrs:       h.attrs,
		groups:      newGroups,
		mu:          h.mu,
		contextKeys: h.contextKeys,
	}
}

// levelColor returns the ANSI color for the given level.
func (h *colorHandler) levelColor(level slog.Level) string {
	switch {
	case level >= slog.LevelError:
		return colorRed
	case level >= slog.LevelWarn:
		return colorYellow
	case level >= slog.LevelInfo:
		return colorGreen
	default:
		return colorBlue
	}
}

// levelString returns the string representation of the level.
func (h *colorHandler) levelString(level slog.Level) string {
	switch {
	case level >= slog.LevelError:
		return "ERROR"
	case level >= slog.LevelWarn:
		return "WARN"
	case level >= slog.LevelInfo:
		return "INFO"
	default:
		return "DEBUG"
	}
}
