package main

import (
	"encoding/json"
	"time"
	"net/http"
)

func refreshHandler(w http.ResponseWriter, r *http.Request) {
	type request struct {
		RefreshToken string `json:"refresh_token"`
	}
	var req request
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil || req.RefreshToken == "" {
		respondJSON(w, http.StatusBadRequest, map[string]string{"message": "Refresh token is required"})
		return
	}

	claims, err := ValidateToken(req.RefreshToken)
	if err != nil {
		respondJSON(w, http.StatusUnauthorized, map[string]string{"message": "Invalid or expired refresh token"})
		return
	}

	if !CheckUserExists(claims.Email) {
		respondJSON(w, http.StatusUnauthorized, map[string]string{"message": "User not found"})
		return
	}

	newToken, err := GenerateJWT(claims.Email, 15*time.Minute)
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, map[string]string{"message": "Failed to generate new token"})
		return
	}

	respondJSON(w, http.StatusOK, map[string]string{"token": newToken})
}