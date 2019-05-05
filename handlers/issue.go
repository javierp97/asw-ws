package handlers

import (
	"asw-project/models"
	"encoding/json"
	"fmt"
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
	key := r.Header.Get("Authentication")
	users, _ := models.GetAllUsers()
	exists := false
	for _, s := range users {
		if s.FirebaseID == key {
			exists = true
		}
	}
	return exists
}

func GetIssue(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	} else {
		vars := mux.Vars(r)

		id, _ := strconv.Atoi(vars["id"])
		issue, _ := models.FindIssueByID(uint(id))

		w.Header().Set("Content-Type", "application/json")
		j, _ := json.Marshal(issue)
		w.WriteHeader(http.StatusOK)
		w.Write(j)
	}
}

//TODO: CHECK NULLS
func CreateIssue(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	} else {
		exists := authenticate(r)
		if exists == true {
			var issue models.Issue
			decoder := json.NewDecoder(r.Body)
			err := decoder.Decode(&issue)
			if err != nil {
				panic(err)
			}
			models.CreateIssue(issue)
			w.Header().Set("Content-Type", "application/json")
			j, _ := json.Marshal(issue)
			fmt.Println(issue)
			w.WriteHeader(http.StatusOK)
			w.Write(j)
		} else {
			w.Write([]byte("{403:Forbidden}"))
		}
	}
}

func DeleteIssue(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	} else {
		exists := authenticate(r)
		if exists == true {
			vars := mux.Vars(r)

			id, _ := strconv.Atoi(vars["id"])
			var issue models.Issue
			issue.ID = uint(id)
			models.DeleteIssue(issue)
			w.WriteHeader(http.StatusOK)
		} else {
			w.Write([]byte("{403:Forbidden}"))
		}
	}
}
