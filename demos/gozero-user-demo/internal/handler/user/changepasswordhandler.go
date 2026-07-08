package user

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"gozero-user-demo/internal/logic/user"
	"gozero-user-demo/internal/svc"
	"gozero-user-demo/internal/types"
)

func ChangePasswordHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.ChangePasswordReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := user.NewChangePasswordLogic(r.Context(), svcCtx)
		resp, err := l.ChangePassword(r, &req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
