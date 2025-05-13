/*
 * Author: Martin <lmccc.dev@gmail.com>
 * Co-Author: AI Assistant
 * Description: This code was collaboratively developed by Martin and AI Assistant.
 */

package log

import "context"

// 使用非导出类型作为 context key 以避免冲突
// (Using unexported type as context key to avoid collisions)
type contextKey int

const (
	// TraceIDKey 是用于在 context 中存储 Trace ID 的键
	// (TraceIDKey is the key for storing Trace ID in context)
	TraceIDKey contextKey = iota // 使用 iota 保证唯一性 (Use iota for uniqueness)
	// RequestIDKey 是用于在 context 中存储 Request ID 的键
	// (RequestIDKey is the key for storing Request ID in context)
	RequestIDKey
)

// --- Helper functions for context (Optional but recommended) ---

// ContextWithTraceID 将 Trace ID 添加到 context 中
// (ContextWithTraceID adds Trace ID to the context)
func ContextWithTraceID(ctx context.Context, traceID string) context.Context {
	return context.WithValue(ctx, TraceIDKey, traceID)
}

// ContextWithRequestID 将 Request ID 添加到 context 中
// (ContextWithRequestID adds Request ID to the context)
func ContextWithRequestID(ctx context.Context, requestID string) context.Context {
	return context.WithValue(ctx, RequestIDKey, requestID)
}

// TraceIDFromContext 从 context 中提取 Trace ID
// (TraceIDFromContext extracts Trace ID from the context)
func TraceIDFromContext(ctx context.Context) (string, bool) {
	val, ok := ctx.Value(TraceIDKey).(string)
	return val, ok
}

// RequestIDFromContext 从 context 中提取 Request ID
// (RequestIDFromContext extracts Request ID from the context)
func RequestIDFromContext(ctx context.Context) (string, bool) {
	val, ok := ctx.Value(RequestIDKey).(string)
	return val, ok
} 