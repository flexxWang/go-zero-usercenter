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

func UserProfiles(users []*model.User) []types.UserProfile {
	items := make([]types.UserProfile, 0, len(users))
	for _, user := range users {
		items = append(items, types.UserProfile{
			Id:        user.ID,
			Email:     user.Email,
			Nickname:  user.Nickname,
			Bio:       user.Bio,
			Avatar:    user.Avatar,
			CreatedAt: user.CreatedAt.Format(time.RFC3339),
		})
	}

	return items
}
