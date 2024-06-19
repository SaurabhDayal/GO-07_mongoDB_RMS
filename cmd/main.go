package main

import (
	"GO-07_mongoDB_RMS/database"
	"GO-07_mongoDB_RMS/server"
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func main() {

	// Connect to MongoDB
	database.ConnectDatabase()

	// Routes set
	r := server.SetupRoutes()

	// Start server
	server := &http.Server{
		Addr:    database.Port,
		Handler: r,
	}

	go func() {
		fmt.Printf("Server is listening on port %s\n", database.Port)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("Could not listen on %s: %v\n", database.Port, err)
		}
	}()

	// Graceful shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
	fmt.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server Shutdown Failed:%+v", err)
	}

	fmt.Println("Closing MongoDB connection...")
	if err := database.DB.Disconnect(database.MongoCtx); err != nil {
		log.Fatalf("Error disconnecting from MongoDB: %v", err)
	}

	fmt.Println("Server exited")
}
