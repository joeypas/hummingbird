package main

import (
	"context"
	"encoding/json"
	"net/http"
)

type registerRequest struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func registerHandler(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	var req registerRequest
	if json.NewDecoder(r.Body).Decode(&req) != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	user, err := createUser(ctx, req.Email, req.Username, req.Password)
	if err != nil {
		http.Error(w, "failed", http.StatusInternalServerError)
		return
	}
	token, err := issueToken(user.id)
	if err != nil {
		http.Error(w, "failed", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(map[string]string{"token": token})
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	var req loginRequest
	if json.NewDecoder(r.Body).Decode(&req) != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	user, err := verifyCredentials(ctx, req.Email, req.Password)
	if err != nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	token, err := issueToken(user.id)
	if err != nil {
		http.Error(w, "failed", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(map[string]string{"token": token})
}
