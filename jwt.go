package auth

import (
	"NoRethink/internal/apperror"
	"NoRethink/internal/config"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/rs/zerolog/log"
)

type Type string

const (
	Access  Type = "ACCESS"
	Refresh Type = "REFRESH"
)

type claim struct {
	Type Type `json:"type"`
	jwt.RegisteredClaims
}

func generateJwt(userID string, t Type, exp time.Time, secret string) (string, *apperror.ErrorResponse) {
	claims := claim{
		Type: t,

		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   userID,
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(exp),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	signedToken, err := token.SignedString([]byte(secret))
	if err != nil {
		log.Error().Err(err).Msg("failed to sign token")
		return "", &apperror.InternalError
	}

	return signedToken, nil
}

func GenerateAccessToken(userID string, jwt *config.JwtConfig) (string, *apperror.ErrorResponse) {
	exp := time.Now().Add(time.Duration(jwt.AccessDurationInMin) * time.Minute)

	return generateJwt(userID, Access, exp, jwt.AccessSecret)
}

func GenerateRefreshToken(userID string, jwt *config.JwtConfig) (string, time.Time, *apperror.ErrorResponse) {
	exp := time.Now().Add(time.Duration(jwt.RefreshDurationInDay) * 24 * time.Hour)
	token, appErr := generateJwt(userID, Refresh, exp, jwt.RefreshSecret)
	return token, exp, appErr
}

func verifyJWT(tokenStr string, secret string) (*claim, *apperror.ErrorResponse) {
	token, err := jwt.ParseWithClaims(tokenStr, &claim{}, func(t *jwt.Token) (interface{}, error) {

		if t.Method != jwt.SigningMethodHS512 {
			return nil, jwt.ErrSignatureInvalid
		}

		return []byte(secret), nil
	})

	if err != nil {
		return nil, &apperror.InvalidToken
	}

	claims, ok := token.Claims.(*claim)
	if !ok || !token.Valid {
		return nil, &apperror.InvalidToken
	}

	return claims, nil
}

func ValidateAccessToken(token string, accessSecret string) (valid bool, userId string, appErr *apperror.ErrorResponse) {
	claim, appErr := verifyJWT(token, accessSecret)
	if appErr != nil {
		return false, "", appErr
	}

	return claim.Type == Access, claim.Subject, nil
}

func ValidateRefreshToken(token string, refreshSecret string) (valid bool, userId string, appErr *apperror.ErrorResponse) {
	claim, appErr := verifyJWT(token, refreshSecret)
	if appErr != nil {
		return false, "", appErr
	}
	return claim.Type == Refresh, claim.Subject, nil
}
