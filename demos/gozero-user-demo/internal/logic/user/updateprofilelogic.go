package user

import (
	"context"
	"database/sql"
	"strings"

	"gozero-user-demo/internal/assembler"
	xauth "gozero-user-demo/internal/auth"
	"gozero-user-demo/internal/model"
	"gozero-user-demo/internal/svc"
	"gozero-user-demo/internal/types"
	"gozero-user-demo/internal/xerr"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateProfileLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateProfileLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateProfileLogic {
	return &UpdateProfileLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateProfileLogic) UpdateProfile(req *types.UpdateProfileReq) (resp *types.UserResp, err error) {
	userID, err := xauth.UserIDFromContext(l.ctx)
	if err != nil {
		return nil, xerr.NewUnauthorized("invalid auth context")
	}

	patch := model.UserPatch{
		Nickname: trimPtr(req.Nickname),
		Bio:      trimPtr(req.Bio),
		Avatar:   trimPtr(req.Avatar),
	}

	user, err := l.svcCtx.UserModel.UpdateProfile(l.ctx, userID, patch)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, xerr.NewNotFound("user not found")
		}

		return nil, xerr.NewInternal("failed to update profile")
	}

	return assembler.UserResp(user), nil
}

func trimPtr(value *string) *string {
	if value == nil {
		return nil
	}

	trimmed := strings.TrimSpace(*value)
	return &trimmed
}
