package main

import (
	"net/http"
	"encoding/json"
	"time"
)

func signinHandler(w http.ResponseWriter, r *http.Request) {
	type request struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	var req request
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil || req.Email == "" || req.Password == "" {
		respondJSON(w, http.StatusBadRequest, map[string]string{"message": "Email and password are required"})
		return
	}

	if !AuthenticateUser(req.Email, req.Password) {
		respondJSON(w, http.StatusUnauthorized, map[string]string{"message": "Invalid email or password"})
		return
	}

	token, err := GenerateJWT(req.Email, 15*time.Minute)
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, map[string]string{"message": "Failed to generate token"})
		return
	}

	refreshToken, err := GenerateJWT(req.Email, 7*24*time.Hour)
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, map[string]string{"message": "Failed to generate refresh token"})
		return
	}

	respondJSON(w, http.StatusOK, map[string]string{"token": token, "refresh_token": refreshToken})
}