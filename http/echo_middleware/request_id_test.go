package echomiddleware

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/luckyAkbar/stdlib/helper"
	"github.com/stretchr/testify/assert"
)

func TestEchoMiddleware_RequestID(t *testing.T) {
	t.Run("ok - use override", func(t *testing.T) {
		ec := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()
		ectx := ec.NewContext(req, rec)

		handler := func(c echo.Context) error {
			ctxID := GetRequestIDFromCtx(c.Request().Context())
			assert.NotEmpty(t, ctxID)
			return c.JSON(http.StatusOK, nil)
		}

		rid := RequestID(true)
		h := rid(handler)
		err := h(ectx)
		assert.NoError(t, err)

		id := rec.Header().Get(echo.HeaderXRequestID)
		assert.NotEmpty(t, id)
	})

	t.Run("ok - don't override", func(t *testing.T) {
		ID := helper.GenerateID()
		ec := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.Header.Set(echo.HeaderXRequestID, ID)
		rec := httptest.NewRecorder()
		ectx := ec.NewContext(req, rec)

		handler := func(c echo.Context) error {
			ctxID := GetRequestIDFromCtx(c.Request().Context())
			assert.NotEmpty(t, ctxID)
			assert.Equal(t, ctxID, ID)
			return c.JSON(http.StatusOK, nil)
		}

		rid := RequestID(false)
		h := rid(handler)
		err := h(ectx)
		assert.NoError(t, err)

		id := rec.Header().Get(echo.HeaderXRequestID)
		assert.NotEmpty(t, id)
		assert.Equal(t, id, ID)
	})
}

func TestEchoMiddleware_GetRequestIDFromContext(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		ctx := setRequestIDToContext(context.TODO(), "123")

		id := GetRequestIDFromCtx(ctx)
		assert.Equal(t, id, "123")
	})

	t.Run("invalid", func(t *testing.T) {
		ctx := context.WithValue(context.TODO(), ReqIDCtxKey, 111222)

		id := GetRequestIDFromCtx(ctx)
		assert.Equal(t, id, "")
	})
}
