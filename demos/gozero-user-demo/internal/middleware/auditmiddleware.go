package middleware

import (
	"fmt"
	"net/http"
	"time"

	xauth "gozero-user-demo/internal/auth"
	"gozero-user-demo/internal/xerr"

	redisstore "github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/rest/httpx"
)

type AuditMiddleware struct {
	redis       *redisstore.Redis
	tokenPrefix string
}

func NewAuditMiddleware(redis *redisstore.Redis, tokenPrefix string) *AuditMiddleware {
	return &AuditMiddleware{
		redis:       redis,
		tokenPrefix: tokenPrefix,
	}
}

func (m *AuditMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		userID, _ := xauth.UserIDFromContext(r.Context())

		token := xauth.BearerToken(r.Header.Get("Authorization"))
		if token == "" {
			httpx.ErrorCtx(r.Context(), w, xerr.NewUnauthorized("missing bearer token"))
			return
		}

		if _, err := m.redis.GetCtx(r.Context(), m.tokenPrefix+token); err != nil {
			if err == redisstore.Nil {
				httpx.ErrorCtx(r.Context(), w, xerr.NewUnauthorized("login session expired"))
				return
			}

			httpx.ErrorCtx(r.Context(), w, xerr.NewInternal("failed to verify login session"))
			return
		}

		next(w, r)

		fmt.Printf("[audit] method=%s path=%s userId=%d cost=%s\n",
			r.Method,
			r.URL.Path,
			userID,
			time.Since(start).String(),
		)
	}
}
