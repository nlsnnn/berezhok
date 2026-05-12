package jwt

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"

	"github.com/nlsnnn/berezhok/internal/modules/auth"
)

type jwtTokenService struct {
	secret []byte
}

func NewTokenService(secret []byte) *jwtTokenService {
	return &jwtTokenService{
		secret: secret,
	}
}

func (s *jwtTokenService) Generate(claims auth.TokenClaims) (string, error) {
	jwtClaims := jwt.MapClaims{
		"user_id":   claims.UserID,
		"user_type": claims.UserType,
		"exp":       time.Now().Add(24 * time.Hour).Unix(),
	}

	if claims.Role != "" {
		jwtClaims["role"] = claims.Role
	}
	if claims.PartnerID != nil {
		jwtClaims["partner_id"] = claims.PartnerID.String()
	}
	if claims.LocationID != nil {
		jwtClaims["location_id"] = claims.LocationID.String()
	}
	if claims.UserData != nil {
		jwtClaims["user_data"] = claims.UserData
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwtClaims)
	return token.SignedString(s.secret)
}

func (s *jwtTokenService) Validate(tokenString string) (*auth.TokenClaims, error) {
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method")
		}
		return s.secret, nil
	})

	if err != nil || !token.Valid {
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	claims := token.Claims.(jwt.MapClaims)

	userID, err := uuid.Parse(claims["user_id"].(string))
	if err != nil {
		return nil, err
	}

	result := &auth.TokenClaims{
		UserID:   userID,
		UserType: claims["user_type"].(string),
	}

	if role, ok := claims["role"].(string); ok {
		result.Role = role
	}
	if partnerIDValue, ok := claims["partner_id"].(string); ok && partnerIDValue != "" {
		partnerID, err := uuid.Parse(partnerIDValue)
		if err != nil {
			return nil, err
		}
		result.PartnerID = &partnerID
	}
	if locationIDValue, ok := claims["location_id"].(string); ok && locationIDValue != "" {
		locationID, err := uuid.Parse(locationIDValue)
		if err != nil {
			return nil, err
		}
		result.LocationID = &locationID
	}
	if userData, ok := claims["user_data"]; ok {
		result.UserData = userData
	}

	return result, nil
}
