package jwt

import (
	"github.com/golang-jwt/jwt/v5"
)

type JwtPayload struct {
	Email string
}

type jwtService struct {
	secret string
}

func NewJwtService(secret string) *jwtService {
	return &jwtService{
		secret: secret,
	}
}

func (j *jwtService) Create(payload JwtPayload) (string, error) {
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email": payload.Email,
	})
	tokenStr, err := t.SignedString([]byte(j.secret))
	if err != nil {
		return "", err
	}
	return tokenStr, nil
}

func (j *jwtService) Parse(token string) (bool, *JwtPayload) {
	t, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		return []byte(j.secret), nil
	})
	if err != nil {
		return false, nil
	}
	email := t.Claims.(jwt.MapClaims)["email"]

	return t.Valid, &JwtPayload{Email: email.(string)}
}
