package user

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"gozero-user-demo/internal/logic/user"
	"gozero-user-demo/internal/svc"
)

func GetMeHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := user.NewGetMeLogic(r.Context(), svcCtx)
		resp, err := l.GetMe()
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
