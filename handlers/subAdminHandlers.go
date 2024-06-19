package handlers

import (
	"GO-07_mongoDB_RMS/middlewares"
	"GO-07_mongoDB_RMS/models"
	"GO-07_mongoDB_RMS/service"
	"encoding/json"
	"net/http"
	"strconv"
)

func CreateRestaurant(w http.ResponseWriter, r *http.Request) {
	var restaurant models.Restaurant
	// get userId and role
	uc := middlewares.GetUserContext(r)

	if err := json.NewDecoder(r.Body).Decode(&restaurant); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	restaurant.OwnedByUserID = uc.UserID

	restaurantResponse, err := service.CreateRestaurant(restaurant)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(restaurantResponse)
}

func CreateDish(w http.ResponseWriter, r *http.Request) {
	var dish models.Dish
	// get userId and role
	uc := middlewares.GetUserContext(r)

	if err := json.NewDecoder(r.Body).Decode(&dish); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	dish.CreatedByUserID = uc.UserID

	restaurantResponse, err := service.CreateDish(dish)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(restaurantResponse)
}

func GetRestaurantByOwnerId(w http.ResponseWriter, r *http.Request) {
	uc := middlewares.GetUserContext(r)

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

	response, err := service.GetRestaurantByOwnerId(limit, offset, *uc)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Set response headers and encode the result as JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func GetDishByCreatedUserId(w http.ResponseWriter, r *http.Request) {
	uc := middlewares.GetUserContext(r)

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

	subAdminsResponse, err := service.GetDishByCreatedUserId(limit, offset, *uc)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Set response headers and encode the result as JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(subAdminsResponse)
}

func GetUsersList(w http.ResponseWriter, r *http.Request) {
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

	subAdminsResponse, err := service.GetUserList(limit, offset)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Set response headers and encode the result as JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(subAdminsResponse)
}
