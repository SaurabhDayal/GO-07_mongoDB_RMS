package handlers

import (
	"GO-07_mongoDB_RMS/service"
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"strconv"
)

func GetRestaurantList(w http.ResponseWriter, r *http.Request) {
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

	response, err := service.GetRestaurantList(limit, offset)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Set response headers and encode the result as JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func GetRestaurantDishList(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "resId")
	restaurantID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

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

	response, err := service.GetRestaurantDishList(limit, offset, restaurantID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Set response headers and encode the result as JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func GetDistance(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "resId")
	restaurantID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	strLatitude := r.URL.Query().Get("latitude")
	latitude, err := strconv.ParseFloat(strLatitude, 64)
	if err != nil {
		latitude = 23.5937
	}

	strLongitude := r.URL.Query().Get("longitude")
	longitude, err := strconv.ParseFloat(strLongitude, 64)
	if err != nil {
		longitude = 78.9629
	}

	response, err := service.GetDistance(latitude, longitude, restaurantID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Set response headers and encode the result as JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)

}

func AddOrder(w http.ResponseWriter, r *http.Request) {
	//var order models.Orders
	//err := json.NewDecoder(r.Body).Decode(&order)
	//if err != nil {
	//	w.WriteHeader(http.StatusInternalServerError)
	//}
	//ord, err := dbHelper.CreateNewOrder(&order)
	//if err != nil {
	//	w.WriteHeader(http.StatusNoContent)
	//} else {
	//	w.WriteHeader(http.StatusOK)
	//	err = json.NewEncoder(w).Encode(ord)
	//	if err != nil {
	//		w.WriteHeader(http.StatusInternalServerError)
	//	}
	//}
}

func CancelOrder(w http.ResponseWriter, r *http.Request) {
	// id distance
	//err := dbHelper.CancelOrder()
	//if err != nil {
	//	w.WriteHeader(http.StatusNoContent)
	//} else {
	//	w.WriteHeader(http.StatusOK)
	//	if err != nil {
	//		w.WriteHeader(http.StatusInternalServerError)
	//	}
	//}
}

func OkOrder(w http.ResponseWriter, r *http.Request) {
	// id distance
	//order, err := dbHelper.OkOrder()
	//if err != nil {
	//	w.WriteHeader(http.StatusNoContent)
	//} else {
	//	w.WriteHeader(http.StatusOK)
	//	err = json.NewEncoder(w).Encode(order)
	//	if err != nil {
	//		w.WriteHeader(http.StatusInternalServerError)
	//	}
	//}
}
