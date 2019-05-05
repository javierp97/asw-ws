package handlers

import (
	"asw-project/models"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	(*w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	(*w).Header().Set("Access-Control-Allow-Credentials", "true")

}

func authenticate(r *http.Request) bool {
	key := r.Header.Get("Authorization")
	users, _ := models.GetAllUsers()
	exists := false
	for _, s := range users {
		if s.FirebaseID == key {
			exists = true
		}
	}
	return exists
}

func Index(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	} else {
		auth := authenticate(r)
		if auth == true {
			enableCors(&w)
			vars := mux.Vars(r)

			id, _ := strconv.Atoi(vars["id"])
			issue, _ := models.FindIssueByID(uint(id))

			w.Header().Set("Content-Type", "application/json")
			j, _ := json.Marshal(issue)
			w.WriteHeader(http.StatusOK)
			w.Write(j)
		} else {
			w.Write([]byte("Cannot autenticate"))
		}
	}
}
