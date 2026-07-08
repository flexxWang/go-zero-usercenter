package user

import (
	"context"

	"gozero-user-demo/internal/assembler"
	"gozero-user-demo/internal/svc"
	"gozero-user-demo/internal/types"
	"gozero-user-demo/internal/xerr"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListUsersLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListUsersLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListUsersLogic {
	return &ListUsersLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListUsersLogic) ListUsers(req *types.ListUsersReq) (resp *types.UserListResp, err error) {
	users, total, err := l.svcCtx.UserModel.List(l.ctx, req.Page, req.PageSize)
	if err != nil {
		return nil, xerr.NewInternal("failed to query users")
	}

	return &types.UserListResp{
		Items:    assembler.UserProfiles(users),
		Page:     req.Page,
		PageSize: req.PageSize,
		Total:    total,
	}, nil
}
