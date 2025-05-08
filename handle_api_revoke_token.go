package main

import (
	"encoding/json"
	"net/http"
)

func revokeHandler(w http.ResponseWriter, r *http.Request) {
	type request struct {
		RefreshToken string `json:"refresh_token,omitempty"`
	}
	var req request
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		respondJSON(w, http.StatusBadRequest, map[string]string{"message": "Invalid request body"})
		return
	}

	tokenStr, err := extractTokenFromHeader(r)
	if err != nil {
		respondJSON(w, http.StatusUnauthorized, map[string]string{"message": err.Error()})
		return
	}

	RevokeToken(tokenStr)
	if req.RefreshToken != "" {
		RevokeToken(req.RefreshToken)
	}

	respondJSON(w, http.StatusOK, map[string]string{"message": "Token(s) revoked successfully"})
}