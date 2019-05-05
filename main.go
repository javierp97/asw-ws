package main

import (
	"asw-project/models"
	"asw-ws/handlers"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

func main() {
	err := models.LoadDB()
	if err != nil {
		fmt.Println(err)
	}

	r := mux.NewRouter().StrictSlash(false)

	//Functions to implement
	r.HandleFunc("/api/issue/{id}", handlers.Index).Methods("GET", "OPTIONS")

	//END
	server := &http.Server{
		Addr:           ":9092",
		Handler:        r,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	log.Println("working")
	server.ListenAndServe()
}
