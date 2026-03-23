package service

import "github.com/nlsnnn/berezhok/internal/modules/auth"

type TokenService interface {
	Generate(claims auth.TokenClaims) (string, error)
	Validate(tokenString string) (*auth.TokenClaims, error)
}
