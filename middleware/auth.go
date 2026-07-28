package middleware

import (
	"context"
	"net/http"

	"github.com/justKody/taskboard-go-api/service/auth"
	"github.com/justKody/taskboard-go-api/utils"
)

type contextKey string

const UserIDKey contextKey = "userID"

func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenString, err := auth.GetTokenFromHeader(r)
		if err != nil {
			utils.WriteError(w, http.StatusUnauthorized, err)
			return
		}

		userID, err := auth.VerifyJWT(tokenString)
		if err != nil {
			utils.WriteError(w, http.StatusUnauthorized, err)
			return
		}

		ctx := context.WithValue(r.Context(), UserIDKey, userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func GetUserID(ctx context.Context) (string, bool) {
	userID, ok := ctx.Value(UserIDKey).(string)
	return userID, ok
}
