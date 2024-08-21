package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"user/internal/user-service/handlers"
	"user/internal/user-service/service"
	"user/internal/user-service/storage"

	"github.com/gorilla/mux"
)

func main() {
	fmt.Println("User Service Project!")
	connStr := os.Getenv("POSTGRES_CONN_STR")
	if connStr == "" {
		log.Fatal("Environment variable POSTGRES_CONN_STR is required")
	}
	store, err := storage.New(connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer store.DB.Close()

	r := mux.NewRouter()
	service := service.Service{S: store}
	handlers := handlers.Handlers{S: service}
	r.HandleFunc("/users", handlers.SignUp).Methods("POST")
	r.HandleFunc("/login", handlers.SignIn).Methods("POST")
	fmt.Println("Starting server at :8080")
	errServ := http.ListenAndServe(":8080", r)
	if errServ != nil {
		fmt.Println("Error happened, %v", errServ.Error)
		return
	}
}
