package main

import (
	"context"
	"errors"
	"net/http"
	"strings"
	"github.com/golang-jwt/jwt/v4"
)

type ContextKey string

const UserEmailKey ContextKey = "userEmail"

func authMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tokenStr, err := extractTokenFromHeader(r)
		if err != nil {
			respondJSON(w, http.StatusUnauthorized, map[string]string{"message": err.Error()})
			return
		}

		claims, err := ValidateToken(tokenStr)
		if err != nil {
			if errors.Is(err, jwt.ErrTokenExpired) {
				respondJSON(w, http.StatusUnauthorized, map[string]string{"message": "Token has expired"})
			} else if err.Error() == "token is revoked" {
				respondJSON(w, http.StatusUnauthorized, map[string]string{"message": "Token has been revoked. Please login again."})
			} else {
				respondJSON(w, http.StatusUnauthorized, map[string]string{"message": "Token is invalid"})
			}
			return
		}

		ctx := context.WithValue(r.Context(), UserEmailKey, claims.Email)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}

func extractTokenFromHeader(r *http.Request) (string, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return "", errors.New("Authorization header missing")
	}
	if !strings.HasPrefix(authHeader, "Bearer ") {
		return "", errors.New("Authorization header must start with 'Bearer '")
	}
	return authHeader[len("Bearer "):], nil
}
