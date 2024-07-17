package main

import (
	"log"
	"net/http"
	_ "time-tracker/docs"
	"time-tracker/internal/db"
	"time-tracker/internal/handlers"
	"time-tracker/internal/migrations"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	httpSwagger "github.com/swaggo/http-swagger"
)

// @title Time Tracker API
// @version 1.0
// @description This is a time tracker server.
// @host localhost:8080
// @BasePath /
func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file: ", err)
	}

	db.Init()

	log.Println("Running database migrations")
	migrations.Migrate(db.Database)
	log.Println("Database migrations completed successfully")

	r := mux.NewRouter()
	r.Use(handlers.LoggingMiddleware)

	r.HandleFunc("/users", handlers.GetUsers).Methods("GET")
	r.HandleFunc("/worklogs", handlers.GetWorkLogs).Methods("GET")
	r.HandleFunc("/tasks/start", handlers.StartTask).Methods("POST")
	r.HandleFunc("/tasks/stop", handlers.StopTask).Methods("POST")
	r.HandleFunc("/users", handlers.AddUser).Methods("POST")
	r.HandleFunc("/users/{id}", handlers.UpdateUser).Methods("PUT")
	r.HandleFunc("/users/{id}", handlers.DeleteUser).Methods("DELETE")

	r.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

	log.Println("Server is starting...")
	http.ListenAndServe(":8080", r)
}
