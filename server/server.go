package server

import (
	"GO-07_mongoDB_RMS/handlers"
	"GO-07_mongoDB_RMS/middlewares"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func SetupRoutes() *chi.Mux {

	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Login-Logout & add Address
	r.Post("/register", handlers.Register)
	r.Post("/login", handlers.Login)

	r.Route("/", func(logoutRoute chi.Router) {
		logoutRoute.Use(middlewares.UserAuthentication)
		logoutRoute.Post("/logout", handlers.Logout)
	})

	// Admin routes
	r.Route("/admin", func(adminRoute chi.Router) {
		adminRoute.Use(middlewares.UserAuthentication)
		adminRoute.Use(middlewares.AdminAuthorization)
		adminRoute.Post("/sub-admin", handlers.CreateSubAdmin)
		adminRoute.Get("/sub-admin", handlers.GetSubAdminList)
	})

	// SubAdmin routes
	r.Route("/sub-admin", func(subAdminRoute chi.Router) {
		subAdminRoute.Use(middlewares.UserAuthentication)
		subAdminRoute.Use(middlewares.SubAdminAuthorization)
		subAdminRoute.Post("/restaurant", handlers.CreateRestaurant)
		subAdminRoute.Get("/restaurants", handlers.GetRestaurantByOwnerId)
		subAdminRoute.Post("/dish", handlers.CreateDish)
		subAdminRoute.Get("/dishes", handlers.GetDishByCreatedUserId)
		subAdminRoute.Get("/users", handlers.GetUsersList)
	})

	// User Routes
	r.Route("/user", func(customerRoute chi.Router) {
		customerRoute.Use(middlewares.UserAuthentication)
		customerRoute.Route("/restaurant", func(r chi.Router) {
			r.Get("/", handlers.GetRestaurantList)
			r.Get("/{resId}/dishes", handlers.GetRestaurantDishList)
			r.Get("/{resId}/distance", handlers.GetDistance)
		})
		customerRoute.Route("/order", func(orderRoute chi.Router) {
			orderRoute.Use(middlewares.CustomerAuthorization)
			orderRoute.Post("/{dishId}", handlers.AddOrder)
			orderRoute.Delete("/{orderId}", handlers.CancelOrder)
			orderRoute.Put("/{orderId}", handlers.OkOrder)
		})
	})

	return r
}
