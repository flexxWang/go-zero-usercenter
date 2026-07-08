package user

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"gozero-user-demo/internal/logic/user"
	"gozero-user-demo/internal/svc"
)

func DeleteProfileHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := user.NewDeleteProfileLogic(r.Context(), svcCtx)
		resp, err := l.DeleteProfile(r)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
