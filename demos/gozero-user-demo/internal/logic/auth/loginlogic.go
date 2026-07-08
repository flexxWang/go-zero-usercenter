package auth

import (
	"context"
	"database/sql"
	"strconv"
	"strings"

	xauth "gozero-user-demo/internal/auth"
	"gozero-user-demo/internal/svc"
	"gozero-user-demo/internal/types"
	"gozero-user-demo/internal/xerr"

	"github.com/zeromicro/go-zero/core/logx"
)

type LoginLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *LoginLogic) Login(req *types.LoginReq) (resp *types.LoginResp, err error) {
	user, err := l.svcCtx.UserModel.FindByEmail(l.ctx, strings.TrimSpace(req.Email))
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, xerr.NewUnauthorized("email or password is incorrect")
		}

		return nil, xerr.NewInternal("failed to query user")
	}

	if err := xauth.ComparePassword(user.PasswordHash, req.Password); err != nil {
		return nil, xerr.NewUnauthorized("email or password is incorrect")
	}

	token, expireAt, refreshAfter, err := xauth.GenerateToken(
		l.svcCtx.Config.Auth.AccessSecret,
		l.svcCtx.Config.Auth.AccessExpire,
		user.ID,
		user.Email,
	)
	if err != nil {
		return nil, xerr.NewInternal("failed to generate access token")
	}

	if err := l.svcCtx.Redis.SetexCtx(
		l.ctx,
		l.svcCtx.Config.Session.TokenPrefix+token,
		strconv.FormatInt(user.ID, 10),
		int(l.svcCtx.Config.Auth.AccessExpire),
	); err != nil {
		return nil, xerr.NewInternal("failed to persist login session")
	}

	return &types.LoginResp{
		AccessToken:  token,
		AccessExpire: expireAt,
		RefreshAfter: refreshAfter,
	}, nil
}
