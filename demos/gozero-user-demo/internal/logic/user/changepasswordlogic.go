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

type ChangePasswordLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewChangePasswordLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ChangePasswordLogic {
	return &ChangePasswordLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ChangePasswordLogic) ChangePassword(r *http.Request, req *types.ChangePasswordReq) (resp *types.EmptyResp, err error) {
	userID, err := xauth.UserIDFromContext(l.ctx)
	if err != nil {
		return nil, xerr.NewUnauthorized("invalid auth context")
	}

	user, err := l.svcCtx.UserModel.FindByID(l.ctx, userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, xerr.NewNotFound("user not found")
		}

		return nil, xerr.NewInternal("failed to query user")
	}

	if err := xauth.ComparePassword(user.PasswordHash, req.OldPassword); err != nil {
		return nil, xerr.NewUnauthorized("old password is incorrect")
	}

	newHash, err := xauth.HashPassword(req.NewPassword)
	if err != nil {
		return nil, xerr.NewInternal("failed to hash new password")
	}

	if err := l.svcCtx.UserModel.UpdatePassword(l.ctx, userID, newHash); err != nil {
		if err == sql.ErrNoRows {
			return nil, xerr.NewNotFound("user not found")
		}

		return nil, xerr.NewInternal("failed to update password")
	}

	token := xauth.BearerToken(r.Header.Get("Authorization"))
	if token != "" {
		if _, err := l.svcCtx.Redis.DelCtx(l.ctx, l.svcCtx.Config.Session.TokenPrefix+token); err != nil {
			return nil, xerr.NewInternal("failed to clear login session")
		}
	}

	return &types.EmptyResp{}, nil
}
