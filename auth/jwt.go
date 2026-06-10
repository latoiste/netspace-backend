package auth

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Claim struct {
	UserId       string `json:"userId"`
	ActorType    string `json:"actorType"`
	LocationSlug string `json:"locationSlug,omitempty"`
	jwt.RegisteredClaims
}

type Auth struct {
	jwtKey []byte
}

func NewAuth(jwtKey []byte) *Auth {
	return &Auth{
		jwtKey: jwtKey,
	}
}

func (a *Auth) generateJWT(userId string, actorType string, locationSlug string) (string, error) {
	claim := Claim{
		UserId:       userId,
		ActorType:    actorType,
		LocationSlug: locationSlug,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 8)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		claim,
	)

	tokenString, err := token.SignedString(a.jwtKey)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func (a *Auth) GenerateUserJWT(userId string) (string, error) {
	return a.generateJWT(userId, "user", "")
}

func (a *Auth) GenerateAdminJWT(adminId string, locationSlug string) (string, error) {
	return a.generateJWT(adminId, "admin", locationSlug)
}

func (a *Auth) VerifyToken(tokenString string) (*jwt.Token, error) {
	claims := &Claim{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Mismatched signing method")
		}
		return a.jwtKey, nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, fmt.Errorf("Invalid token")
	}

	return token, nil
}

func (a *Auth) ExtractTokenFromHeader(r *http.Request) (string, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return "", errors.New("Authorization header not given")
	}

	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
		return "", errors.New("No token given")
	}

	return parts[1], nil
}
