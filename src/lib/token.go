package lib

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

var (
	ErrInvalidToken     = errors.New("invalid token")
	ErrTokenExpired     = errors.New("token expired")
	ErrTokenMalformed   = errors.New("token malformed")
	ErrMissingSecretKey = errors.New("missing JWT secret key")
)

type Claims struct {
	UserId uuid.UUID `json:"user_id"`
	jwt.RegisteredClaims
}

var jwtSecret []byte

func InitialiseJWT(secret string) {
	jwtSecret = []byte(secret)
}

func GenerateToken(id uuid.UUID) (string, error) {
	if len(jwtSecret) == 0 {
		return "", ErrMissingSecretKey
	}

	claims := Claims{
		UserId: id,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "adire-apparel",
			Subject:   id.String(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

func ValidateToken(tokenString string) (*Claims, error) {
	if len(jwtSecret) == 0 {
		return nil, ErrMissingSecretKey
	}

	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrTokenMalformed
		}
		return jwtSecret, nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrTokenExpired
		}
		if errors.Is(err, jwt.ErrTokenMalformed) {
			return nil, ErrTokenMalformed
		}
		return nil, ErrInvalidToken
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, ErrInvalidToken
	}

	return claims, nil
}

func RefreshToken(tokenString string) (string, error) {
	claims, err := ValidateToken(tokenString)
	if err != nil && !errors.Is(err, ErrTokenExpired) {
		return "", err
	}

	if claims == nil {
		return "", ErrInvalidToken
	}

	return GenerateToken(claims.UserId)
}

func ExtractUserID(tokenString string) (uuid.UUID, error) {
	claims, err := ValidateToken(tokenString)
	if err != nil {
		return uuid.Nil, err
	}

	return claims.UserId, nil
}
