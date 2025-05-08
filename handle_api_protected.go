package main

import (
	"net/http"
)

func protectedHandler(w http.ResponseWriter, r *http.Request) {
	tokenStr, _ := extractTokenFromHeader(r)
	claims, err := ValidateToken(tokenStr)
	if err != nil {
		respondJSON(w, http.StatusUnauthorized, map[string]string{"message": "Invalid token"})
		return
	}

	respondJSON(w, http.StatusOK, map[string]string{
		"message": "Hello " + claims.Email + "! You have accessed a protected route.",
	})
}