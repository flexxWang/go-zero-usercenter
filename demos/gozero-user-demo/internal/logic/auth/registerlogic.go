package auth

import (
	"context"
	"strings"

	"gozero-user-demo/internal/assembler"
	xauth "gozero-user-demo/internal/auth"
	"gozero-user-demo/internal/svc"
	"gozero-user-demo/internal/types"
	"gozero-user-demo/internal/xerr"

	"github.com/zeromicro/go-zero/core/logx"
)

type RegisterLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRegisterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RegisterLogic {
	return &RegisterLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RegisterLogic) Register(req *types.RegisterReq) (resp *types.UserResp, err error) {
	hash, err := xauth.HashPassword(req.Password)
	if err != nil {
		return nil, xerr.NewInternal("failed to hash password")
	}

	user, err := l.svcCtx.UserModel.Create(
		l.ctx,
		strings.TrimSpace(req.Email),
		hash,
		strings.TrimSpace(req.Nickname),
	)
	if err != nil {
		if isDuplicateEntry(err) {
			return nil, xerr.NewConflict("email already registered")
		}

		return nil, xerr.NewInternal("failed to create user")
	}

	return assembler.UserResp(user), nil
}

func isDuplicateEntry(err error) bool {
	return err != nil && strings.Contains(err.Error(), "Duplicate entry")
}
