package user

import (
	"context"
	"database/sql"

	"gozero-user-demo/internal/assembler"
	xauth "gozero-user-demo/internal/auth"
	"gozero-user-demo/internal/svc"
	"gozero-user-demo/internal/types"
	"gozero-user-demo/internal/xerr"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetMeLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetMeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetMeLogic {
	return &GetMeLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetMeLogic) GetMe() (resp *types.UserResp, err error) {
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

	return assembler.UserResp(user), nil
}
