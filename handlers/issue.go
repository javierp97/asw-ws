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
	//	fmt.Println(models.ExistKind(issue.Type))
	//	fmt.Println()
	correct = (correct && models.ExistStatus(issue.Status) && models.ExistPriority(issue.Priority) && models.ExistKind(issue.Type))
	//	fmt.Println("check params: ", correct)
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
				panic(err) //TODO
			}
			issue.Status = "New"
			issue.Votes = 0
			//fmt.Println("filepath: " + issue.FilePath)
			if checkParams(issue) == false || checkNulls(issue) == false {
				//	fmt.Println(checkNulls(issue))
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

func UpdateState(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	} else {
		w.Header().Set("Content-Type", "application/json")
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
			if models.ExistStatus(newIssue.Status) == false {
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte(`{"Error":"The parameter is incorrect"}`))
				return
			}
			actualIssue.Status = newIssue.Status
			models.UpdateIssue(actualIssue)
			var comment models.Comment
			owner, err := models.GetCommentOwnerById(uint(id))
			comment.Content = "The status of this issue changed to " + newIssue.Status
			comment.OwnerID = r.Header.Get("Authorization")
			comment.OwnerName = owner
			comment.IssueID = uint(id)
			models.CreateComment(comment)

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

func UpdateIssue(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	} else {
		w.Header().Set("Content-Type", "application/json")
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
		w.Header().Set("Content-Type", "application/json")
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
			w.WriteHeader(http.StatusForbidden)
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
			w.WriteHeader(http.StatusForbidden)
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
			w.WriteHeader(http.StatusForbidden)
			w.Write([]byte(`{"403":"Forbbiden"}`))
		}
	}
}

func CreateComment(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	} else {
		w.Header().Set("Content-Type", "application/json")
		exists := authenticate(r)
		if exists == true {
			decoder := json.NewDecoder(r.Body)
			var newComment models.Comment
			err := decoder.Decode(&newComment)

			if err != nil {
				panic("error")
			}

			vars := mux.Vars(r)
			id, _ := strconv.Atoi(vars["id"])

			exists := models.ExistsIssue(uint(id))

			if exists == false {
				w.WriteHeader(http.StatusNotFound)
				w.Write([]byte(`{"Error":"The issue does not exist"}`))
				return
			}

			user, _ := models.FindUserByID(r.Header.Get("Authorization"))
			newComment.OwnerID = user.FirebaseID
			newComment.OwnerName = user.Username
			newComment.IssueID = uint(id)

			commentId, _ := models.CreateCommentAndReturnId(newComment)
			commentIdString := strconv.Itoa(int(commentId))
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("{Ok: Comment created with id: " + commentIdString + "}"))
			return

		} else {
			w.WriteHeader(http.StatusForbidden)
			w.Write([]byte(`{"403":"Forbbiden"}`))
			return
		}
	}
}

func EditComment(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	} else {
		w.Header().Set("Content-Type", "application/json")
		exists := authenticate(r)
		if exists == true {
			decoder := json.NewDecoder(r.Body)
			var newComment models.Comment
			err := decoder.Decode(&newComment)

			if err != nil {
				panic("error")
			}

			vars := mux.Vars(r)
			id, _ := strconv.Atoi(vars["commentId"])

			user, _ := models.FindUserByID(r.Header.Get("Authorization"))
			owner, err := models.GetCommentOwnerById(uint(id))

			if err != nil {
				w.WriteHeader(http.StatusNotFound)
				w.Write([]byte(`{"Error":"The comment does not exist"}`))
				return
			}

			if owner != user.FirebaseID {
				w.WriteHeader(http.StatusForbidden)
				w.Write([]byte(`{"Error":"This comment is not yours"}`))
				return
			}

			models.UpdateCommentById(uint(id), newComment.Content)

			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"OK":"Comment updated succesfully}"`))
			return

		} else {
			w.WriteHeader(http.StatusForbidden)
			w.Write([]byte(`{"403":"Forbbiden"}`))
			return
		}
	}
}

func DeleteComment(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	} else {
		w.Header().Set("Content-Type", "application/json")
		exists := authenticate(r)
		if exists == true {
			vars := mux.Vars(r)
			//id, _ := strconv.Atoi(vars["id"])
			commentId, _ := strconv.Atoi(vars["commentId"])
			user, _ := models.FindUserByID(r.Header.Get("Authorization"))
			owner, err := models.GetCommentOwnerById(uint(commentId))

			if err != nil {
				w.WriteHeader(http.StatusNotFound)
				w.Write([]byte(`{"Error":"The comment does not exists"}`))
				return
			}

			if owner != user.FirebaseID {
				w.WriteHeader(http.StatusForbidden)
				w.Write([]byte(`{"Error":"You are not the owner of the issue so you can not delete it"}`))
				return
			}

			models.DeleteCommentById(uint(commentId))
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"200":"Comment deleted succesfully"}`))

		} else {
			w.WriteHeader(http.StatusForbidden)
			w.Write([]byte(`{"403":"Forbbiden"}`))
			return
		}
	}
}
