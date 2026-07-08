package user

import (
	"context"
	"database/sql"

	"gozero-user-demo/internal/assembler"
	"gozero-user-demo/internal/svc"
	"gozero-user-demo/internal/types"
	"gozero-user-demo/internal/xerr"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetUserLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserLogic {
	return &GetUserLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetUserLogic) GetUser(req *types.GetUserReq) (resp *types.UserResp, err error) {
	user, err := l.svcCtx.UserModel.FindByID(l.ctx, req.UserId)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, xerr.NewNotFound("user not found")
		}

		return nil, xerr.NewInternal("failed to query user")
	}

	return assembler.UserResp(user), nil
}
