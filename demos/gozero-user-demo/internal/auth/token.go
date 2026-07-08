package auth

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"gozero-user-demo/internal/xerr"

	"github.com/golang-jwt/jwt/v4"
)

func GenerateToken(secret string, expireSeconds int64, userID int64, email string) (string, int64, int64, error) {
	now := time.Now().Unix()
	expireAt := now + expireSeconds

	claims := jwt.MapClaims{
		"userId": userID,
		"email":  email,
		"exp":    expireAt,
		"iat":    now,
	}

	token := jwt.New(jwt.SigningMethodHS256)
	token.Claims = claims

	signed, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", 0, 0, err
	}

	return signed, expireAt, 0, nil
}

func GenerateRefreshToken() (string, error) {
	buf := make([]byte, 32)
	if _, err := rand.Read(buf); err != nil {
		return "", err
	}

	return base64.RawURLEncoding.EncodeToString(buf), nil
}

func UserIDFromContext(ctx context.Context) (int64, error) {
	raw := ctx.Value("userId")
	if raw == nil {
		return 0, xerr.NewUnauthorized("missing userId claim")
	}

	switch value := raw.(type) {
	case int64:
		return value, nil
	case int:
		return int64(value), nil
	case float64:
		return int64(value), nil
	case json.Number:
		return value.Int64()
	case string:
		return strconv.ParseInt(value, 10, 64)
	default:
		return 0, errors.New("unsupported userId claim type")
	}
}

func BearerToken(header string) string {
	if !strings.HasPrefix(header, "Bearer ") {
		return ""
	}

	return strings.TrimSpace(strings.TrimPrefix(header, "Bearer "))
}

func RefreshTokenKey(prefix, token string) string {
	return prefix + token
}

func RefreshUserKey(prefix string, userID int64) string {
	return fmt.Sprintf("%s%d", prefix, userID)
}
