package jwt

type JwtService interface {
	Parse(token string) (bool, *JwtPayload)
	Create(payload JwtPayload) (string, error)
}
