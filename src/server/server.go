package main

import (
	"fmt"
	"net/http"
)

type pushRequest struct {
	Key  string `json:"key"`
	Note string `json:"note"`
}

type getRequest struct {
	Key string `json:"key"`
}

type getResponce struct {
	Notes []string `json:"notes"`
}

type failureResponce struct {
	Err string `json:"error"`
}

func main() {
	err := runServer()
	if err != nil {
		fmt.Print(err.Error())
	}
}

func runServer() error {
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
