package pkg

import (
	"time"

	jwt "github.com/golang-jwt/jwt/v5"
	"github.com/hunderaweke/metsasft/config"
	"github.com/hunderaweke/metsasft/internal/domain"
)

const (
	RefreshTokenDuration = time.Minute
	AccessTokenDuration  = time.Second * 10
)

type UserClaims struct {
	jwt.RegisteredClaims
	UserID           string `json:"user_id"`
	Email            string `json:"email"`
	TelegramUsername string `json:"telegram_username"`
	Type             string `json:"type"`
	IsAdmin          bool   `json:"is_admin"`
}
type ErrInvalidClaim struct{}

func (e ErrInvalidClaim) Error() string {
	return "invalid claims"
}

func GenerateToken(user domain.User) (string, string, error) {
	claims := UserClaims{
		UserID:           user.ID,
		Email:            user.Email,
		TelegramUsername: user.TelegramUsername,
		Type:             "refresh",
		IsAdmin:          user.IsAdmin,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "metsasft",
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(RefreshTokenDuration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}
	config, err := config.LoadConfig()
	if err != nil {
		return "", "", err
	}
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	refreshToken, err := t.SignedString([]byte(config.Jwt.Secret))
	claims.ExpiresAt = jwt.NewNumericDate(time.Now().Add(AccessTokenDuration))
	claims.Type = "access"
	t.Claims = claims
	accessToken, err := t.SignedString([]byte(config.Jwt.Secret))
	if err != nil {
		return "", "", err
	}
	return refreshToken, accessToken, nil
}

func ValidateRefreshToken(token string) (*UserClaims, error) {
	config, err := config.LoadConfig()
	if err != nil {
		return nil, err
	}
	t, err := jwt.ParseWithClaims(token, &UserClaims{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(config.Jwt.Secret), nil
	})
	if err != nil {
		return nil, err
	}
	if !t.Valid {
		return nil, jwt.ErrSignatureInvalid
	} else if claims, ok := t.Claims.(*UserClaims); ok {
		if claims.Type != "refresh" || claims.ExpiresAt.Time.Before(time.Now()) {
			return nil, ErrInvalidClaim{}
		}
		return claims, nil
	}
	return nil, ErrInvalidClaim{}
}

func ValidateAccessToken(token string) (*UserClaims, error) {
	config, err := config.LoadConfig()
	if err != nil {
		return nil, err
	}
	t, err := jwt.ParseWithClaims(token, &UserClaims{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(config.Jwt.Secret), nil
	})
	if err != nil {
		return nil, err
	}
	if !t.Valid {
		return nil, jwt.ErrSignatureInvalid
	} else if claims, ok := t.Claims.(*UserClaims); ok {
		if claims.Type != "access" || claims.ExpiresAt.Time.Before(time.Now()) {
			return nil, ErrInvalidClaim{}
		}
		return claims, nil
	}
	return nil, ErrInvalidClaim{}
}
