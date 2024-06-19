package middlewares

import (
	"GO-07_mongoDB_RMS/models"
	"GO-07_mongoDB_RMS/utils"
	"context"
	"errors"
	"github.com/golang-jwt/jwt/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"os"
)

type userContextType string

const (
	UserContextKey userContextType = "user_context"
)

func UserAuthentication(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("authorization")
		claims := &models.UserClaims{}
		if token == "" {
			http.Error(w, "token not sent in header", http.StatusBadRequest)
			return
		} else {
			// for logout functionality
			if !utils.IsTokenValid(token) {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			parseToken, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
				return []byte(os.Getenv("jwtSecret")), nil
			})
			if err != nil {
				if errors.Is(err, jwt.ErrSignatureInvalid) {
					http.Error(w, "Invalid token signature", http.StatusUnauthorized)
					return
				}
				http.Error(w, "Token is expired", http.StatusUnauthorized)
				return
			}
			if !parseToken.Valid {
				http.Error(w, "Invalid token", http.StatusUnauthorized)
				return
			}
			if claims.UserID == primitive.NilObjectID {
				http.Error(w, "Invalid token", http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), UserContextKey, claims)
			next.ServeHTTP(w, r.WithContext(ctx))
		}
	})
}

func GetUserContext(req *http.Request) *models.UserClaims {
	return req.Context().Value(UserContextKey).(*models.UserClaims)
}

func AdminAuthorization(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userCtx := r.Context().Value(UserContextKey).(*models.UserClaims)
		for _, role := range userCtx.Roles {
			if role == models.AdminRole {
				next.ServeHTTP(w, r)
				return
			}
		}

		http.Error(w, "admin only routes", http.StatusForbidden)
		return
	})
}

func SubAdminAuthorization(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userCtx := r.Context().Value(UserContextKey).(*models.UserClaims)
		for _, role := range userCtx.Roles {
			if role == models.AdminRole || role == models.SubAdminRole {
				next.ServeHTTP(w, r)
				return
			}
		}

		http.Error(w, "admin only routes", http.StatusForbidden)
		return
	})
}

func CustomerAuthorization(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userCtx := r.Context().Value(UserContextKey).(*models.UserClaims)
		for _, role := range userCtx.Roles {
			if role == models.CustomerRole {
				next.ServeHTTP(w, r)
				return
			}
		}

		http.Error(w, "admin only routes", http.StatusForbidden)
		return
	})
}
