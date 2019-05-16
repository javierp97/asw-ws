package main

import (
	"asw-project/models"
	"fmt"
	"log"
	"net/http"
	"time"

	"asw-ws/handlers"

	"github.com/gorilla/mux"
)

func main() {
	err := models.LoadDB()
	if err != nil {
		fmt.Println(err)
	}

	r := mux.NewRouter().StrictSlash(false)

	//Functions to implement
	r.HandleFunc("/api/issue/{id}", handlers.GetIssue).Methods("GET", "OPTIONS")
	r.HandleFunc("/api/issue", handlers.GetAllIssues).Methods("GET", "OPTIONS")
	r.HandleFunc("/api/issue", handlers.CreateIssue).Methods("POST", "OPTIONS")
	r.HandleFunc("/api/issue/{id}/vote", handlers.VoteIssue).Methods("POST", "OPTIONS")
	r.HandleFunc("/api/issue/{id}/vote", handlers.UnVoteIssue).Methods("DELETE", "OPTIONS")
	r.HandleFunc("/api/issue/{id}", handlers.DeleteIssue).Methods("DELETE", "OPTIONS")
	r.HandleFunc("/api/issue/{id}/attach", handlers.PostAttachment).Methods("POST", "OPTIONS")
	r.HandleFunc("/api/issue/attach/{idattach}", handlers.DeleteAttachment).Methods("DELETE", "OPTIONS")
	r.HandleFunc("/api/issue/{id}", handlers.UpdateIssue).Methods("PUT", "OPTIONS")
	r.HandleFunc("/api/issue/{id}/state", handlers.UpdateState).Methods("PUT", "OPTIONS")
	r.HandleFunc("/api/issue/{id}/comments", handlers.CreateComment).Methods("POST", "OPTIONS")
	r.HandleFunc("/api/comments/{commentId}", handlers.DeleteComment).Methods("DELETE", "OPTIONS")
	r.HandleFunc("/api/comments/{commentId}", handlers.EditComment).Methods("PUT", "OPTIONS")
	r.HandleFunc("/api/issue/{id}/watch", handlers.WatchIssue).Methods("POST", "OPTIONS")
	r.HandleFunc("/api/issue/{id}/watch", handlers.UnWatchIssue).Methods("DELETE", "OPTIONS")
	r.HandleFunc("/api/comments/{commentId}", handlers.GetComment).Methods("GET", "OPTIONS")
	r.HandleFunc("/api/users", handlers.GetUsers).Methods("GET", "OPTIONS")

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
