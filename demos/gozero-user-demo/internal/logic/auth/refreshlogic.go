package auth

import (
	"context"
	"database/sql"
	"strconv"
	"time"

	redisstore "github.com/zeromicro/go-zero/core/stores/redis"
	xauth "gozero-user-demo/internal/auth"
	"gozero-user-demo/internal/svc"
	"gozero-user-demo/internal/types"
	"gozero-user-demo/internal/xerr"

	"github.com/zeromicro/go-zero/core/logx"
)

type RefreshLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRefreshLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RefreshLogic {
	return &RefreshLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RefreshLogic) Refresh(req *types.RefreshTokenReq) (resp *types.LoginResp, err error) {
	refreshTokenKey := xauth.RefreshTokenKey(l.svcCtx.Config.Session.RefreshTokenPrefix, req.RefreshToken)
	userIDRaw, err := l.svcCtx.Redis.GetCtx(l.ctx, refreshTokenKey)
	if err != nil {
		if err == redisstore.Nil {
			return nil, xerr.NewUnauthorized("refresh token expired")
		}

		return nil, xerr.NewInternal("failed to verify refresh session")
	}

	userID, err := strconv.ParseInt(userIDRaw, 10, 64)
	if err != nil {
		return nil, xerr.NewInternal("invalid refresh session")
	}

	user, err := l.svcCtx.UserModel.FindByID(l.ctx, userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, xerr.NewNotFound("user not found")
		}

		return nil, xerr.NewInternal("failed to query user")
	}

	newToken, expireAt, _, err := xauth.GenerateToken(
		l.svcCtx.Config.Auth.AccessSecret,
		l.svcCtx.Config.Auth.AccessExpire,
		user.ID,
		user.Email,
	)
	if err != nil {
		return nil, xerr.NewInternal("failed to generate access token")
	}

	newRefreshToken, err := xauth.GenerateRefreshToken()
	if err != nil {
		return nil, xerr.NewInternal("failed to generate refresh token")
	}

	if err := l.svcCtx.Redis.SetexCtx(
		l.ctx,
		l.svcCtx.Config.Session.TokenPrefix+newToken,
		strconv.FormatInt(user.ID, 10),
		int(l.svcCtx.Config.Auth.AccessExpire),
	); err != nil {
		return nil, xerr.NewInternal("failed to persist refreshed session")
	}

	if _, err := l.svcCtx.Redis.DelCtx(l.ctx, refreshTokenKey); err != nil {
		return nil, xerr.NewInternal("failed to rotate refresh session")
	}

	if err := l.svcCtx.Redis.SetexCtx(
		l.ctx,
		xauth.RefreshTokenKey(l.svcCtx.Config.Session.RefreshTokenPrefix, newRefreshToken),
		strconv.FormatInt(user.ID, 10),
		int(l.svcCtx.Config.Auth.RefreshExpire),
	); err != nil {
		return nil, xerr.NewInternal("failed to persist refresh session")
	}

	if err := l.svcCtx.Redis.SetexCtx(
		l.ctx,
		xauth.RefreshUserKey(l.svcCtx.Config.Session.RefreshUserPrefix, user.ID),
		newRefreshToken,
		int(l.svcCtx.Config.Auth.RefreshExpire),
	); err != nil {
		return nil, xerr.NewInternal("failed to persist refresh session index")
	}

	refreshExpireAt := time.Now().Unix() + l.svcCtx.Config.Auth.RefreshExpire

	return &types.LoginResp{
		AccessToken:   newToken,
		AccessExpire:  expireAt,
		RefreshToken:  newRefreshToken,
		RefreshExpire: refreshExpireAt,
	}, nil
}
