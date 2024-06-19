package handlers

import (
	"GO-07_mongoDB_RMS/models"
	"GO-07_mongoDB_RMS/service"
	"encoding/json"
	"net/http"
	"strconv"
)

func CreateSubAdmin(w http.ResponseWriter, r *http.Request) {
	var subAdmin models.User
	if err := json.NewDecoder(r.Body).Decode(&subAdmin); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	userResponse, err := service.CreateSubAdmin(subAdmin)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(userResponse)
}

func GetSubAdminList(w http.ResponseWriter, r *http.Request) {
	strLimit := r.URL.Query().Get("limit")
	limit, err := strconv.Atoi(strLimit)
	if err != nil {
		limit = 10
	}

	strOffset := r.URL.Query().Get("offset")
	offset, err := strconv.Atoi(strOffset)
	if err != nil {
		offset = 0
	}

	subAdminsResponse, err := service.GetSubAdminList(limit, offset)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Set response headers and encode the result as JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(subAdminsResponse)
}
