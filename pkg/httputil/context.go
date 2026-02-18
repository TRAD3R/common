// Package httputil provides HTTP utilities for request tracing and context propagation
package httputil

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// contextKey is a custom type for context keys to avoid collisions
type contextKey string

const (
	// requestIDKey is the context key for request ID
	requestIDKey contextKey = "request_id"

	// HeaderRequestID is the standard request ID header
	HeaderRequestID = "X-Request-ID"

	// HeaderCorrelationID is the correlation ID header for distributed tracing
	HeaderCorrelationID = "X-Correlation-ID"
)

// RequestIDKey is the public string constant for gin.Context.Set/Get
// Note: This is safe to use as string because gin.Context uses its own internal storage
const RequestIDKey = "request_id"

// GetRequestID extracts request_id from gin.Context or generates a new one
func GetRequestID(c *gin.Context) string {
	if requestID := c.GetString(RequestIDKey); requestID != "" {
		return requestID
	}
	return uuid.New().String()
}

// GetRequestIDFromContext extracts request_id from context.Context
func GetRequestIDFromContext(ctx context.Context) string {
	if ginCtx, ok := ctx.(*gin.Context); ok {
		return GetRequestID(ginCtx)
	}

	// Try to extract from context value using typed key
	if requestID, ok := ctx.Value(requestIDKey).(string); ok && requestID != "" {
		return requestID
	}

	return uuid.New().String()
}

// PropagateRequestIDFromContext adds request ID headers from context.Context
// Use this when you don't have access to gin.Context but have context with request_id
//
// Usage:
//
//	req, _ := http.NewRequestWithContext(ctx, "GET", url, nil)
//	httputil.PropagateRequestIDFromContext(ctx, req)
//	resp, err := client.Do(req)
func PropagateRequestIDFromContext(ctx context.Context, req *http.Request) {
	requestID := GetRequestIDFromContext(ctx)
	req.Header.Set(HeaderRequestID, requestID)
	req.Header.Set(HeaderCorrelationID, requestID)
}

// ContextWithRequestID creates a new context with request_id value
// Useful for passing request ID to goroutines or async operations
func ContextWithRequestID(ctx context.Context, requestID string) context.Context {
	return context.WithValue(ctx, requestIDKey, requestID)
}

// ContextFromGin creates a new context from gin.Context with request_id propagated
// Use this when calling service methods that need request tracing
//
// Usage:
//
//	ctx := httputil.ContextFromGin(c)
//	result, err := h.service.DoSomething(ctx, params)
func ContextFromGin(c *gin.Context) context.Context {
	requestID := GetRequestID(c)
	return ContextWithRequestID(c.Request.Context(), requestID)
}