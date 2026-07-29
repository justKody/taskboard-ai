package auth

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/justKody/taskboard-go-api/config"
)

const TokenCookieName = "token"

type Claims struct {
	UserId string `json:"id"`
	jwt.RegisteredClaims
}

func CreateJWT(userId string) (string, error) {
	expirationDuration := time.Duration(config.Envs.JWTExpiration) * time.Minute

	claims := &Claims{
		UserId: userId,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expirationDuration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "taskboard-go-api",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(config.Envs.JWTKey))
}

func VerifyJWT(tokenString string) (string, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(config.Envs.JWTKey), nil
	})
	if err != nil {
		return "", err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return "", errors.New("invalid token")
	}

	return claims.UserId, nil
}

func SetTokenCookie(w http.ResponseWriter, token string) {
	http.SetCookie(w, &http.Cookie{
		Name:     TokenCookieName,
		Value:    token,
		Path:     "/",
		MaxAge:   config.Envs.JWTExpiration * 60,
		HttpOnly: true,
		Secure:   false, // set true in production over HTTPS
		SameSite: http.SameSiteLaxMode,
	})
}

func ClearTokenCookie(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:     TokenCookieName,
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
	})
}

func GetTokenFromCookie(r *http.Request) (string, error) {
	cookie, err := r.Cookie(TokenCookieName)
	if err != nil {
		return "", errors.New("auth cookie is required")
	}
	if cookie.Value == "" {
		return "", errors.New("auth cookie is empty")
	}
	return cookie.Value, nil
}

func GetTokenFromHeader(r *http.Request) (string, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return "", errors.New("authorization header is required")
	}

	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") || parts[1] == "" {
		return "", errors.New("authorization header must be Bearer <token>")
	}

	return parts[1], nil
}

func GetToken(r *http.Request) (string, error) {
	if token, err := GetTokenFromCookie(r); err == nil {
		return token, nil
	}
	return GetTokenFromHeader(r)
}
