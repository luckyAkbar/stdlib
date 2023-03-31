// Package echomiddleware contains middleware function to be used by echo framework
package echomiddleware

import (
	"context"

	"github.com/labstack/echo/v4"
	"github.com/luckyAkbar/stdlib/helper"
)

// ReqIDCtxKeyType is the type for context key for request ID
type ReqIDCtxKeyType string

// ReqIDCtxKey is the key for request ID in context
const ReqIDCtxKey ReqIDCtxKeyType = "github.com/luckyAkbar/stdlib:echo_middleware:ReqIDCtxKey"

// RequestID is a middleware to generate request ID and set it to context
// if in the request header already have request ID in key `X-Request-ID` and `overrideHeader` is true,
// then the request ID will be override with new generated ID
// otherwise, both on context and response header will use supplied request header `X-Request-ID`
func RequestID(overrideHeader bool) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			req := c.Request()
			res := c.Response()
			rid := req.Header.Get(echo.HeaderXRequestID)

			if overrideHeader {
				rid = helper.GenerateID()
			}

			ctx := setRequestIDToContext(c.Request().Context(), rid)
			c.SetRequest(c.Request().WithContext(ctx))
			res.Header().Set(echo.HeaderXRequestID, rid)

			return next(c)
		}
	}
}

// GetRequestIDFromCtx is a helper function to get request ID from context
func GetRequestIDFromCtx(ctx context.Context) string {
	id, ok := ctx.Value(ReqIDCtxKey).(string)
	if !ok {
		return ""
	}

	return id
}

// can only set request ID from this middleware. Other can only read
func setRequestIDToContext(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, ReqIDCtxKey, id)
}
