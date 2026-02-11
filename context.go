package xlogging

import "context"

// ContextKey is a type for context keys used by xlogging.
type ContextKey string

// Predefined context keys.
const (
	KeyRequestID ContextKey = "request_id"
	KeyTraceID   ContextKey = "trace_id"
	KeySpanID    ContextKey = "span_id"
	KeyUserID    ContextKey = "user_id"
)

// WithRequestID adds a request ID to the context.
func WithRequestID(ctx context.Context, requestID string) context.Context {
	return context.WithValue(ctx, KeyRequestID, requestID)
}

// WithTraceID adds a trace ID to the context.
func WithTraceID(ctx context.Context, traceID string) context.Context {
	return context.WithValue(ctx, KeyTraceID, traceID)
}

// WithSpanID adds a span ID to the context.
func WithSpanID(ctx context.Context, spanID string) context.Context {
	return context.WithValue(ctx, KeySpanID, spanID)
}

// WithUserID adds a user ID to the context.
func WithUserID(ctx context.Context, userID string) context.Context {
	return context.WithValue(ctx, KeyUserID, userID)
}

// GetRequestID retrieves the request ID from the context.
func GetRequestID(ctx context.Context) (string, bool) {
	v, ok := ctx.Value(KeyRequestID).(string)
	return v, ok
}

// GetTraceID retrieves the trace ID from the context.
func GetTraceID(ctx context.Context) (string, bool) {
	v, ok := ctx.Value(KeyTraceID).(string)
	return v, ok
}

// GetSpanID retrieves the span ID from the context.
func GetSpanID(ctx context.Context) (string, bool) {
	v, ok := ctx.Value(KeySpanID).(string)
	return v, ok
}

// GetUserID retrieves the user ID from the context.
func GetUserID(ctx context.Context) (string, bool) {
	v, ok := ctx.Value(KeyUserID).(string)
	return v, ok
}
