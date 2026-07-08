package model

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	redisstore "github.com/zeromicro/go-zero/core/stores/redis"
)

const userCachePrefix = "usercenter:user:"

type User struct {
	ID           int64     `json:"id"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"passwordHash"`
	Nickname     string    `json:"nickname"`
	Bio          string    `json:"bio"`
	Avatar       string    `json:"avatar"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
}

type UserPatch struct {
	Nickname *string
	Bio      *string
	Avatar   *string
}

type UserModel struct {
	db       *sql.DB
	redis    *redisstore.Redis
	cacheTTL int
}

func NewUserModel(db *sql.DB, redis *redisstore.Redis, cacheTTL int) *UserModel {
	return &UserModel{
		db:       db,
		redis:    redis,
		cacheTTL: cacheTTL,
	}
}

func (m *UserModel) Create(ctx context.Context, email, passwordHash, nickname string) (*User, error) {
	result, err := m.db.ExecContext(ctx, `
		INSERT INTO users (email, password_hash, nickname, bio, avatar, created_at, updated_at)
		VALUES (?, ?, ?, '', '', NOW(), NOW())
	`, email, passwordHash, nickname)
	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	return m.FindByID(ctx, id)
}

func (m *UserModel) FindByEmail(ctx context.Context, email string) (*User, error) {
	row := m.db.QueryRowContext(ctx, `
		SELECT id, email, password_hash, nickname, bio, avatar, created_at, updated_at
		FROM users
		WHERE email = ?
		LIMIT 1
	`, email)

	return scanUser(row)
}

func (m *UserModel) FindByID(ctx context.Context, id int64) (*User, error) {
	if cached, err := m.getFromCache(ctx, id); err == nil {
		return cached, nil
	}

	row := m.db.QueryRowContext(ctx, `
		SELECT id, email, password_hash, nickname, bio, avatar, created_at, updated_at
		FROM users
		WHERE id = ?
		LIMIT 1
	`, id)

	user, err := scanUser(row)
	if err != nil {
		return nil, err
	}

	_ = m.setCache(ctx, user)
	return user, nil
}

func (m *UserModel) UpdateProfile(ctx context.Context, id int64, patch UserPatch) (*User, error) {
	sets := make([]string, 0, 4)
	args := make([]any, 0, 4)

	if patch.Nickname != nil {
		sets = append(sets, "nickname = ?")
		args = append(args, *patch.Nickname)
	}
	if patch.Bio != nil {
		sets = append(sets, "bio = ?")
		args = append(args, *patch.Bio)
	}
	if patch.Avatar != nil {
		sets = append(sets, "avatar = ?")
		args = append(args, *patch.Avatar)
	}

	sets = append(sets, "updated_at = NOW()")
	args = append(args, id)

	query := fmt.Sprintf("UPDATE users SET %s WHERE id = ?", strings.Join(sets, ", "))
	result, err := m.db.ExecContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return nil, err
	}
	if affected == 0 {
		return nil, sql.ErrNoRows
	}

	_ = m.deleteCache(ctx, id)
	return m.FindByID(ctx, id)
}

func (m *UserModel) UpdatePassword(ctx context.Context, id int64, passwordHash string) error {
	result, err := m.db.ExecContext(ctx, `
		UPDATE users
		SET password_hash = ?, updated_at = NOW()
		WHERE id = ?
	`, passwordHash, id)
	if err != nil {
		return err
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return sql.ErrNoRows
	}

	return m.deleteCache(ctx, id)
}

func (m *UserModel) getFromCache(ctx context.Context, id int64) (*User, error) {
	val, err := m.redis.GetCtx(ctx, userCacheKey(id))
	if err != nil {
		return nil, err
	}

	var user User
	if err := json.Unmarshal([]byte(val), &user); err != nil {
		return nil, err
	}

	return &user, nil
}

func (m *UserModel) setCache(ctx context.Context, user *User) error {
	if m.cacheTTL <= 0 {
		return nil
	}

	data, err := json.Marshal(user)
	if err != nil {
		return err
	}

	return m.redis.SetexCtx(ctx, userCacheKey(user.ID), string(data), m.cacheTTL)
}

func (m *UserModel) deleteCache(ctx context.Context, id int64) error {
	_, err := m.redis.DelCtx(ctx, userCacheKey(id))
	return err
}

func scanUser(scanner interface {
	Scan(dest ...any) error
}) (*User, error) {
	var user User
	if err := scanner.Scan(
		&user.ID,
		&user.Email,
		&user.PasswordHash,
		&user.Nickname,
		&user.Bio,
		&user.Avatar,
		&user.CreatedAt,
		&user.UpdatedAt,
	); err != nil {
		return nil, err
	}

	return &user, nil
}

func userCacheKey(id int64) string {
	return fmt.Sprintf("%s%d", userCachePrefix, id)
}
