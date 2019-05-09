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
	key := r.Header.Get("Authorization")
	users, _ := models.GetAllUsers()
	exists := false
	for _, s := range users {
		fmt.Println("firebaseID: " + s.FirebaseID + " key: " + key)
		if s.FirebaseID == key {
			exists = true
		}
	}
	fmt.Println(exists)
	return exists
}

func checkParams(issue models.Issue) bool {
	correct := true
	fmt.Println(models.ExistKind(issue.Type))
	fmt.Println()
	correct = (correct && models.ExistStatus(issue.Status) && models.ExistPriority(issue.Priority) && models.ExistKind(issue.Type))
	fmt.Println("check params: ", correct)
	return correct
}

func checkNulls(issue models.Issue) bool {
	correct := true
	if issue.Title == "" || issue.Description == "" || issue.Priority == "" || issue.Type == "" || issue.Assignee == "" || issue.Reporter == "" {
		correct = false
	}
	fmt.Println("check nuls: ", correct)
	return correct
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
			issue.Status = "New"
			issue.Votes = 0
			fmt.Println("filepath: " + issue.FilePath)
			if checkParams(issue) == false || checkNulls(issue) == false {
				fmt.Println(checkNulls(issue))
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte(`{"Error":"Wrong parameters"}`))
				return
			}
			models.CreateIssue(issue)
			w.Header().Set("Content-Type", "application/json")
			j, _ := json.Marshal(issue)
			fmt.Println(issue)
			w.WriteHeader(http.StatusOK)
			w.Write(j)
		} else {
			w.Write([]byte(`{"403":"Forbbiden"}`))
		}
	}
}

type BodyID struct {
	ID string `json:"id"`
}

func VoteIssue(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	} else {
		w.Header().Set("Content-Type", "application/json")
		exists := authenticate(r)
		if exists == true {
			key := r.Header.Get("Authorization")
			fmt.Println(key)
			vars := mux.Vars(r)
			id, _ := strconv.Atoi(vars["id"])
			var voteIssue models.VotedIssue
			voteIssue.IDIssue = uint(id)
			voteIssue.UserID = key
			fmt.Println("no peto")
			if id == 0 {
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte(`{"Error":"Wrong parameters"}`))
				return
			}
			iss, erriss := models.FindIssueByID(voteIssue.IDIssue)
			if erriss != nil || iss.Title == "" {
				w.WriteHeader(http.StatusNotFound)
				w.Write([]byte(`{"Error":"The issue does not exist"}`))
				return
			}
			if models.VoteIssue(voteIssue) != nil {
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte(`{"Error":"The issue has been already voted"}`))
				return
			}
			models.VoteThisIssue(uint(id))
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"200":"OK"}`))
		} else {
			w.Write([]byte(`{"403":"Forbbiden"}`))
		}
	}
}

func UnVoteIssue(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	} else {
		w.Header().Set("Content-Type", "application/json")
		exists := authenticate(r)
		if exists == true {
			key := r.Header.Get("Authorization")
			vars := mux.Vars(r)
			id, _ := strconv.Atoi(vars["id"])

			b, _ := models.IsVoted(key, uint(id))
			if id == 0 {
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte(`{"Error":"Wrong parameters"}`))
				return
			}
			var voteIssue models.VotedIssue
			voteIssue.IDIssue = uint(id)
			voteIssue.UserID = key
			iss, erriss := models.FindIssueByID(voteIssue.IDIssue)
			if erriss != nil || iss.Title == "" {
				w.WriteHeader(http.StatusNotFound)
				w.Write([]byte(`{"Error":"The issue does not exist"}`))
				return
			}
			if !b {
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte(`{"Error":"The issue is not voted by this user"}`))
				return
			}

			if models.UnvoteIssue(voteIssue) != nil {
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte(`{"Error":"Error voting the issue"}`))
				return
			}
			models.VoteThisIssue(uint(id))
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"200":"OK"}`))
		} else {
			w.Write([]byte(`{"403":"Forbbiden"}`))
		}
	}
}
