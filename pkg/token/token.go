package token

import (
	"errors"
	"time"

	"github.com/dusk-chancellor/dc-sso/internal/models"
	"github.com/dusk-chancellor/dc-sso/pkg/zaplog"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// jwt tokens implementation

// errors
var (
	ErrInvalidToken = errors.New("invalid token")
)

// Custom claims
type Claims struct {
	UserID string
	Email  string
	jwt.RegisteredClaims
}

// generates jwt based on user data, secret and expiration time
func GenerateToken(user *models.User, secret string, expiration time.Duration) (string, error) {
	claims := Claims{
		user.ID.String(),
		user.Email,
		jwt.RegisteredClaims{
			Subject: user.Username,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiration)),
			IssuedAt: jwt.NewNumericDate(time.Now()),
			ID: generateJTI(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// parsing token and returning its claims content
func ParseToken(in, secret string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(in, &Claims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}

		return []byte(secret), nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenSignatureInvalid) {
			return nil, ErrInvalidToken
		}
		return nil, err
	}

	claims, ok := token.Claims.(*Claims)
	if !token.Valid || !ok {
		zaplog.Log("invalid token or claims")
		return nil, ErrInvalidToken
	}

	return claims, nil
}

func generateJTI() string {
	return uuid.NewString()
}
