package server

import (
	"net/http"
)

type pushRequest struct {
	Credentials credentials `json:"credentials"`
	Note        string      `json:"note"`
}

type getRequest struct {
	Credentials credentials `json:"credentials"`
}

type getResponce struct {
	Notes []string `json:"notes"`
}

type credentials struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type failureResponce struct {
	Error string `json:"error"`
}

func Run() error {
	plugInHandlers()
	return http.ListenAndServe(":80", nil)
}

func plugInHandlers() {
	http.HandleFunc("/push_note", pushNote)
	http.HandleFunc("/get_notes", getNotes)
}

func pushNote(w http.ResponseWriter, r *http.Request) {

}

func getNotes(w http.ResponseWriter, r *http.Request) {

}
