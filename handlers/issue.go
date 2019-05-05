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
}

func authenticate(key string) bool {
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

	key := r.Header.Get("apiKey")
	auth := authenticate(key)
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
