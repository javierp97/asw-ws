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
		if s.FirebaseID == key {
			exists = true
		}
	}
	//fmt.Println(exists)
	return exists
}

func checkParams(issue models.Issue) bool {
	correct := true
	correct = (correct && models.ExistStatus(issue.Status) && models.ExistPriority(issue.Priority) && models.ExistKind(issue.Type))
	return correct
}

func checkNulls(issue models.Issue) bool {
	correct := true
	if issue.Title == "" || issue.Description == "" || issue.Priority == "" || issue.Type == "" || issue.Assignee == "" || issue.Reporter == "" {
		correct = false
	}
	//fmt.Println("check nuls: ", correct)
	return correct
}

func GetAllIssues(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	} else {
		issues, _ := models.GetAllIssues()
		j, _ := json.Marshal(issues)
		w.WriteHeader(http.StatusOK)
		w.Write(j)
	}
}

func GetIssue(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	} else {
		vars := mux.Vars(r)
		id, _ := strconv.Atoi(vars["id"])
		issue, err := models.FindIssueByID(uint(id))
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(`{"Error":"The requested issue could not be found"}`))
			return
		}

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
				w.WriteHeader(http.StatusBadRequest)
			}
			issue.Status = "New"
			issue.Votes = 0
			if checkParams(issue) == false || checkNulls(issue) == false {
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte(`{"Error":"Wrong parameters"}`))
				return
			}
			id, _ := models.CreateIssue(issue)
			issueResp, _ := models.FindIssueByID(id)
			w.WriteHeader(http.StatusOK)
			w.Header().Set("Content-Type", "application/json")
			j, _ := json.Marshal(issueResp)
			w.Write(j)
		} else {
			w.WriteHeader(http.StatusForbidden)
			w.Write([]byte(`{"403":"Forbbiden"}`))
		}
	}
}

func updateIssue(actualIssue *models.Issue, newIssue models.Issue) {

	title := newIssue.Title
	if title != "" {
		actualIssue.Title = title
	}

	desc := newIssue.Description
	if desc != "" {
		actualIssue.Description = desc
	}

	prio := newIssue.Priority
	if prio != "" && models.ExistPriority(prio) {
		actualIssue.Priority = prio
	}

	kind := newIssue.Type
	if kind != "" && models.ExistKind(kind) {
		actualIssue.Type = kind
	}

	status := newIssue.Status
	if status != "" && models.ExistStatus(status) {
		actualIssue.Status = status
	}

	assig := newIssue.Assignee

	if assig != "" {
		actualIssue.Assignee = assig
	}

	file := newIssue.FilePath

	if file != "" {
		actualIssue.FilePath = file
	}

}

func PutIssue(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	} else {
		exists := authenticate(r)
		if exists == true {
			decoder := json.NewDecoder(r.Body)
			var newIssue models.Issue
			err := decoder.Decode(&newIssue)
			if err != nil {
				panic(err) //TODO
			}
			vars := mux.Vars(r)
			id, _ := strconv.Atoi(vars["id"])
			actualIssue, err := models.FindIssueByID(uint(id))
			if err != nil {
				w.WriteHeader(http.StatusNotFound)
				w.Write([]byte(`{"Error":"The requested issue could not be found"}`))
				return
			}

			updateIssue(&actualIssue, newIssue)
			models.UpdateIssue(actualIssue)
			w.WriteHeader(http.StatusOK)
			w.Header().Set("Content-Type", "application/json")
			j, _ := json.Marshal(actualIssue)
			w.Write(j)

		} else {
			w.WriteHeader(http.StatusForbidden)
			w.Write([]byte(`{"Error":"You do not have access to do this request"}`))
			return
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
			issue, error := models.FindIssueByID(uint(id))
			fmt.Println(error)
			err := models.DeleteIssue(issue)

			if err != nil || error != nil {
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte(`{"400":"Bad Request"}`))
				return
			}
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"IssueDeleted":"OK"}`))
		} else {
			w.Write([]byte(`{"403":"Forbbiden"}`))
		}
	}
}
