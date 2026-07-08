package auth

import (
	"context"
	"net/http"

	xauth "gozero-user-demo/internal/auth"
	"gozero-user-demo/internal/svc"
	"gozero-user-demo/internal/types"
	"gozero-user-demo/internal/xerr"

	"github.com/zeromicro/go-zero/core/logx"
)

type LogoutLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewLogoutLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LogoutLogic {
	return &LogoutLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *LogoutLogic) Logout(r *http.Request) (resp *types.EmptyResp, err error) {
	token := xauth.BearerToken(r.Header.Get("Authorization"))
	if token == "" {
		return nil, xerr.NewUnauthorized("missing bearer token")
	}

	if _, err := l.svcCtx.Redis.DelCtx(l.ctx, l.svcCtx.Config.Session.TokenPrefix+token); err != nil {
		return nil, xerr.NewInternal("failed to clear login session")
	}

	return &types.EmptyResp{}, nil
}
