package assembler

import (
	"time"

	"gozero-user-demo/internal/model"
	"gozero-user-demo/internal/types"
)

func UserResp(user *model.User) *types.UserResp {
	return &types.UserResp{
		User: types.UserProfile{
			Id:        user.ID,
			Email:     user.Email,
			Nickname:  user.Nickname,
			Bio:       user.Bio,
			Avatar:    user.Avatar,
			CreatedAt: user.CreatedAt.Format(time.RFC3339),
		},
	}
}
