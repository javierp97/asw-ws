package handlers

import (
	"asw-project/models"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
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

const maxUploadSize = 8 * 1024 * 1024 // 2 mb
const uploadPath = "./tmp"

func PostAttachment(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	} else {
		w.Header().Set("Content-Type", "application/json")
		exists := authenticate(r)
		if exists == true {
			//key := r.Header.Get("Authorization")
			vars := mux.Vars(r)
			id, _ := strconv.Atoi(vars["id"])

			//Check if the issue exists
			iss, erriss := models.FindIssueByID(uint(id))
			if erriss != nil || iss.Title == "" {
				w.WriteHeader(http.StatusNotFound)
				w.Write([]byte(`{"Error":"The issue does not exist"}`))
				return
			}
			r.Body = http.MaxBytesReader(w, r.Body, maxUploadSize)
			if err := r.ParseForm(); err != nil {
				fmt.Println(err)
				renderError(w, "ERROR_PROCESING_FILE", http.StatusBadRequest)
				return
			}

			// parse and validate file and post parameters
			fileType := r.PostFormValue("type")
			file, _, err := r.FormFile("uploadFile")
			if err != nil {
				fmt.Println("printeo invalid_file")
				fmt.Println(err)
				renderError(w, "INVALID_FILE", http.StatusBadRequest)
				return
			}
			defer file.Close()
			fileBytes, err := ioutil.ReadAll(file)
			if err != nil {
				fmt.Println("printeo invalid_file")
				fmt.Println(err)

				renderError(w, "INVALID_FILE", http.StatusBadRequest)
				return
			}

			// check file type, detectcontenttype only needs the first 512 bytes
			filetype := http.DetectContentType(fileBytes)
			var typefile string
			switch filetype {
			case "image/jpeg":
				typefile = ".jpeg"
			case "image/jpg":
				typefile = ".jpg"
			case "image/gif":
				typefile = ".gif"
			case "image/png":
				typefile = ".png"
			case "application/pdf":
				typefile = ".pdf"
				break
			default:
				renderError(w, "INVALID_FILE_TYPE", http.StatusBadRequest)
				return
			}

			fileName := randToken(12)
			/*
				fileEndings, err := mime.ExtensionsByType(fileType)
				if err != nil {
					fmt.Println(filetype)
					renderError(w, "CANT_READ_FILE_TYPE", http.StatusInternalServerError)
					return
				}
			*/
			newPath := filepath.Join(uploadPath, fileName+typefile) //fileEndings[0])
			fmt.Printf("FileType: %s, File: %s\n", fileType, newPath)

			// write file
			newFile, err := os.Create(newPath)
			if err != nil {
				renderError(w, "CANT_WRITE_FILE", http.StatusInternalServerError)
				return
			}
			defer newFile.Close() // idempotent, okay to call twice
			if _, err := newFile.Write(fileBytes); err != nil || newFile.Close() != nil {
				renderError(w, "CANT_WRITE_FILE", http.StatusInternalServerError)
				return
			}
			var at models.Attachment
			at.FilePath = newPath
			at.IssueID = iss.ID
			errs := models.CreateAttachment(at)
			if errs != nil {
				renderError(w, "CANT_SAVE_FILE", http.StatusInternalServerError)
				return
			}
			w.WriteHeader(http.StatusOK)
			resp, _ := json.Marshal(at)
			w.Write([]byte(resp))
		} else {
			w.Write([]byte(`{"403":"Forbbiden"}`))
		}
	}
}

func renderError(w http.ResponseWriter, message string, statusCode int) {
	w.WriteHeader(http.StatusBadRequest)
	w.Write([]byte(message))
}

func randToken(len int) string {
	b := make([]byte, len)
	rand.Read(b)
	return fmt.Sprintf("%x", b)
}

func PutAttachment(w http.ResponseWriter, r *http.Request) {
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

func DeleteAttachment(w http.ResponseWriter, r *http.Request) {
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
