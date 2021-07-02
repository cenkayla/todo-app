package main

import (
	"log"
	"net/http"

	"github.com/cenkayla/todo-app/middleware"
	"github.com/gorilla/mux"
)

func main() {

	r := mux.NewRouter()
	r.HandleFunc("/api/task", middleware.GetAllTasks).Methods("GET")
	r.HandleFunc("/api/task", middleware.Create).Methods("POST")
	r.HandleFunc("/api/deletetask/{id}", middleware.Delete).Methods("DELETE")
	r.HandleFunc("/api/deleteAllTask", middleware.DeleteAll).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":8080", r))
}
