package handlers

import (
	"GO-07_mongoDB_RMS/models"
	"GO-07_mongoDB_RMS/service"
	"encoding/json"
	"net/http"
	"strings"
)

func Register(w http.ResponseWriter, r *http.Request) {
	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	userResponse, err := service.RegisterCustomer(user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(userResponse)
}

func Login(w http.ResponseWriter, r *http.Request) {
	var credentials models.LoginCredentials
	if err := json.NewDecoder(r.Body).Decode(&credentials); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	credentials.Email = strings.TrimSpace(strings.ToLower(credentials.Email))

	loginResponse, err := service.LoginUser(credentials)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(loginResponse)
}

func Logout(w http.ResponseWriter, r *http.Request) {
	tokenString := r.Header.Get("Authorization")
	err := service.LogoutUser(tokenString)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	} else {
		w.WriteHeader(http.StatusOK)
	}
}
