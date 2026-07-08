package auth

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"gozero-user-demo/internal/logic/auth"
	"gozero-user-demo/internal/svc"
)

func LogoutHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := auth.NewLogoutLogic(r.Context(), svcCtx)
		resp, err := l.Logout(r)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
