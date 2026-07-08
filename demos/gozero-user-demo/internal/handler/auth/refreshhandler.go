package auth

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"gozero-user-demo/internal/logic/auth"
	"gozero-user-demo/internal/svc"
	"gozero-user-demo/internal/types"
)

func RefreshHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.RefreshTokenReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := auth.NewRefreshLogic(r.Context(), svcCtx)
		resp, err := l.Refresh(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
