package types

import (
	"net/mail"
	"strings"

	"gozero-user-demo/internal/xerr"
)

func (r RegisterReq) Validate() error {
	if err := validateEmail(r.Email); err != nil {
		return err
	}
	if len(strings.TrimSpace(r.Password)) < 6 {
		return xerr.NewBadRequest("password must be at least 6 characters")
	}
	if err := validateNickname(r.Nickname); err != nil {
		return err
	}

	return nil
}

func (r LoginReq) Validate() error {
	if err := validateEmail(r.Email); err != nil {
		return err
	}
	if strings.TrimSpace(r.Password) == "" {
		return xerr.NewBadRequest("password is required")
	}

	return nil
}

func (r GetUserReq) Validate() error {
	if r.UserId <= 0 {
		return xerr.NewBadRequest("userId must be a positive integer")
	}

	return nil
}

func (r *ListUsersReq) Validate() error {
	if r.Page == 0 {
		r.Page = 1
	}
	if r.PageSize == 0 {
		r.PageSize = 10
	}
	if r.Page < 1 {
		return xerr.NewBadRequest("page must be greater than 0")
	}
	if r.PageSize < 1 || r.PageSize > 100 {
		return xerr.NewBadRequest("pageSize must be between 1 and 100")
	}

	return nil
}

func (r UpdateProfileReq) Validate() error {
	if r.Nickname == nil && r.Bio == nil && r.Avatar == nil {
		return xerr.NewBadRequest("at least one field must be provided")
	}

	if r.Nickname != nil {
		if err := validateNickname(*r.Nickname); err != nil {
			return err
		}
	}
	if r.Bio != nil && len(strings.TrimSpace(*r.Bio)) > 160 {
		return xerr.NewBadRequest("bio must be at most 160 characters")
	}
	if r.Avatar != nil && len(strings.TrimSpace(*r.Avatar)) > 255 {
		return xerr.NewBadRequest("avatar must be at most 255 characters")
	}

	return nil
}

func (r ChangePasswordReq) Validate() error {
	if strings.TrimSpace(r.OldPassword) == "" {
		return xerr.NewBadRequest("oldPassword is required")
	}
	if len(strings.TrimSpace(r.NewPassword)) < 6 {
		return xerr.NewBadRequest("newPassword must be at least 6 characters")
	}
	if r.OldPassword == r.NewPassword {
		return xerr.NewBadRequest("newPassword must be different from oldPassword")
	}

	return nil
}

func (r RefreshTokenReq) Validate() error {
	if strings.TrimSpace(r.RefreshToken) == "" {
		return xerr.NewBadRequest("refreshToken is required")
	}

	return nil
}

func validateEmail(email string) error {
	email = strings.TrimSpace(email)
	if email == "" {
		return xerr.NewBadRequest("email is required")
	}

	if _, err := mail.ParseAddress(email); err != nil {
		return xerr.NewBadRequest("email format is invalid")
	}

	return nil
}

func validateNickname(nickname string) error {
	nickname = strings.TrimSpace(nickname)
	if nickname == "" {
		return xerr.NewBadRequest("nickname is required")
	}
	if len(nickname) < 2 || len(nickname) > 20 {
		return xerr.NewBadRequest("nickname must be between 2 and 20 characters")
	}

	return nil
}
