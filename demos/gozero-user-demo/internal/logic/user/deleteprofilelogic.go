package user

import (
	"context"
	"database/sql"
	"net/http"

	xauth "gozero-user-demo/internal/auth"
	"gozero-user-demo/internal/svc"
	"gozero-user-demo/internal/types"
	"gozero-user-demo/internal/xerr"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteProfileLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteProfileLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteProfileLogic {
	return &DeleteProfileLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteProfileLogic) DeleteProfile(r *http.Request) (resp *types.EmptyResp, err error) {
	userID, err := xauth.UserIDFromContext(l.ctx)
	if err != nil {
		return nil, xerr.NewUnauthorized("invalid auth context")
	}

	if err := l.svcCtx.UserModel.DeleteByID(l.ctx, userID); err != nil {
		if err == sql.ErrNoRows {
			return nil, xerr.NewNotFound("user not found")
		}

		return nil, xerr.NewInternal("failed to delete user")
	}

	token := xauth.BearerToken(r.Header.Get("Authorization"))
	if token != "" {
		if _, err := l.svcCtx.Redis.DelCtx(l.ctx, l.svcCtx.Config.Session.TokenPrefix+token); err != nil {
			return nil, xerr.NewInternal("failed to clear login session")
		}
	}

	refreshUserKey := xauth.RefreshUserKey(l.svcCtx.Config.Session.RefreshUserPrefix, userID)
	if refreshToken, err := l.svcCtx.Redis.GetCtx(l.ctx, refreshUserKey); err == nil {
		if _, err := l.svcCtx.Redis.DelCtx(l.ctx, xauth.RefreshTokenKey(l.svcCtx.Config.Session.RefreshTokenPrefix, refreshToken)); err != nil {
			return nil, xerr.NewInternal("failed to clear refresh session")
		}
	}
	if _, err := l.svcCtx.Redis.DelCtx(l.ctx, refreshUserKey); err != nil {
		return nil, xerr.NewInternal("failed to clear refresh session index")
	}

	return &types.EmptyResp{}, nil
}
