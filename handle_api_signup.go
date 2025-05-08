package main

import (
	"encoding/json"
	"errors"
	"net/http"
)

func signupHandler(w http.ResponseWriter, r *http.Request) {
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

	err = AddUser(req.Email, req.Password)
	if err != nil {
		if errors.Is(err, ErrUserExists) {
			respondJSON(w, http.StatusBadRequest, map[string]string{"message": err.Error()})
		} else {
			respondJSON(w, http.StatusInternalServerError, map[string]string{"message": "Failed to create user"})
		}
		return
	}
	respondJSON(w, http.StatusCreated, map[string]string{"message": "User created successfully"})
}
