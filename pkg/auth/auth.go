package auth

import (
	"context"
	"net/http"
	"strings"

	httptransport "github.com/go-kit/kit/transport/http"

	"github.com/go-kit/kit/endpoint"
)

// Unauthorized represents an authorization error.
type Unauthorized struct {
}

// ErrorCode is an implementation of the StatusCoder interface in go-kit/http.
func (Unauthorized) ErrorCode() int {
	return -32000
}

// Error is an implementation of the Error interface.
func (Unauthorized) Error() string {
	return "unauthorized"
}

type ContextKey int

const (
	// Deprecated: ContextKeyRequestXToolName is populated in the context by
	// ContextKeyRequestXUserID in the context by
	// PopulateRequestContext. Its value is r.Header.Get("X-User-Id").
	ContextKeyRequestXUserID ContextKey = iota

	// ContextKeyRequestXRequestID from nginx request_id
	ContextKeyRequestXRequestID
)

func ExtractFromContext() httptransport.RequestFunc {
	return func(ctx context.Context, request *http.Request) context.Context {
		for ctxKey, hKey := range map[ContextKey]string{
			ContextKeyRequestXUserID:    "X-User-Id",
			ContextKeyRequestXRequestID: "X-Request-Id",
		} {
			v := ctx.Value(ctxKey)
			switch v := v.(type) {
			case string:
				request.Header.Set(hKey, v)
			case []string:
				request.Header.Set(hKey, strings.Join(v, ";"))
			}
		}
		return ctx
	}
}

// PopulateRequestContext is a RequestFunc that populates several values into
// the context from the HTTP request.
func PopulateRequestContext(ctx context.Context, r *http.Request) context.Context {
	for k, v := range map[ContextKey]string{
		ContextKeyRequestXUserID:    r.Header.Get("X-User-Id"),
		ContextKeyRequestXRequestID: r.Header.Get("X-Request-Id"),
	} {
		ctx = context.WithValue(ctx, k, v)
	}
	return ctx
}

// Middleware returns a Authentication middleware.
func Middleware() endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (interface{}, error) {
			personID, ok := ctx.Value(ContextKeyRequestXUserID).(string)
			// fmt.Println(personID)
			if !ok || personID == "" {
				return nil, Unauthorized{}
			}
			return next(ctx, request)
		}
	}
}

// PersonIDFromContext get user id from context.
func PersonIDFromContext(ctx context.Context) string {
	personID, _ := ctx.Value(ContextKeyRequestXUserID).(string)
	return personID
}

// RequestIDFromContext get requestID from context
func RequestIDFromContext(ctx context.Context) string {
	requestID, _ := ctx.Value(ContextKeyRequestXRequestID).(string)
	return requestID
}
