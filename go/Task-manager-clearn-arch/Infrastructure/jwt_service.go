package infrastructure

import (
	"time"

	"task-manager-clean-arch/domain"

	"github.com/golang-jwt/jwt/v4"
)

type JWTService interface {
	GenerateTokens(user *domain.User) (string, string, error)
	ValidateAccessToken(tokenStr string) (*domain.JwtCustomClaims, error)
	ValidateRefreshToken(tokenStr string) (*domain.JwtCustomRefreshClaims, error)
}

type jwtService struct {
	secret        string
	accessExpiry  time.Duration
	refreshExpiry time.Duration
}

func NewJWTService(secret string, accessExp, refreshExp time.Duration) JWTService {
	return &jwtService{
		secret:        secret,
		accessExpiry:  accessExp,
		refreshExpiry: refreshExp,
	}
}
func (js *jwtService) GenerateTokens(user *domain.User) (string, string, error) {
	accessClaims := &domain.JwtCustomClaims{
		ID:   user.ID.Hex(),
		Role: string(user.Role),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(js.accessExpiry)),
		},
	}

	refreshClaims := &domain.JwtCustomRefreshClaims{
		ID: user.ID.Hex(),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(js.refreshExpiry)),
		},
	}

	accessToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims).SignedString([]byte(js.secret))
	if err != nil {
		return "", "", err
	}

	refreshToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims).SignedString([]byte(js.secret))
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

func (js *jwtService) ValidateAccessToken(tokenStr string) (*domain.JwtCustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &domain.JwtCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrInvalidKey
		}
		return []byte(js.secret), nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*domain.JwtCustomClaims)
	if !ok || !token.Valid {
		return nil, jwt.ErrTokenMalformed
	}

	return claims, nil
}

func (js *jwtService) ValidateRefreshToken(tokenStr string) (*domain.JwtCustomRefreshClaims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &domain.JwtCustomRefreshClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrInvalidKey
		}
		return []byte(js.secret), nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*domain.JwtCustomRefreshClaims)
	if !ok || !token.Valid {
		return nil, jwt.ErrTokenMalformed
	}

	return claims, nil
}
