package responsex

import (
	"context"
	"net/http"

	"gozero-user-demo/internal/xerr"

	"github.com/zeromicro/go-zero/core/logc"
)

type Body struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

func OkHandler(_ context.Context, data any) any {
	return Body{
		Code:    0,
		Message: "OK",
		Data:    data,
	}
}

func ErrorHandler(ctx context.Context, err error) (int, any) {
	if codeErr, ok := err.(*xerr.CodeError); ok {
		return codeErr.Status, Body{
			Code:    codeErr.Code,
			Message: codeErr.Message,
		}
	}

	logc.Errorf(ctx, "unexpected error: %v", err)
	return http.StatusInternalServerError, Body{
		Code:    xerr.CodeInternal,
		Message: "internal server error",
	}
}
