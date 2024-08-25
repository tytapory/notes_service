package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"notes_service/database"
)

type pushRequest struct {
	Credentials credentials `json:"credentials"`
	Note        string      `json:"note"`
}

type registerRequest struct {
	Credentials credentials `json:"credentials"`
}

type getRequest struct {
	Credentials credentials `json:"credentials"`
}

type getResponse struct {
	Notes []string `json:"notes"`
}

type pushResponse struct {
	Message string `json:"message"`
}

type registerResponse struct {
	Message string `json:"message"`
}

type credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type failureResponse struct {
	Error string `json:"error"`
}

func Run() error {
	plugInHandlers()
	fmt.Printf("Running server on :8080")
	return http.ListenAndServe(":8080", nil)
}

func plugInHandlers() {
	http.HandleFunc("/push_note", pushNote)
	http.HandleFunc("/get_notes", getNotes)
	http.HandleFunc("/register", register)
}

func register(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var request registerRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(failureResponse{Error: err.Error()})
		return
	}

	err = database.RegisterUser(request.Credentials.Username, request.Credentials.Password)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(failureResponse{Error: err.Error()})
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(registerResponse{Message: "User successfully regustered"})
}

func pushNote(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var request pushRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(failureResponse{Error: err.Error()})
		return
	}

	userID, err := database.Authenticate(request.Credentials.Username, request.Credentials.Password)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(failureResponse{Error: err.Error()})
		return
	}

	validatedNote, err := verify(request.Note)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(failureResponse{Error: err.Error()})
		return
	}
	err = database.InsertUserNote(userID, validatedNote)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(failureResponse{Error: err.Error()})
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(pushResponse{Message: "Note successfully created"})
}

func getNotes(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var request getRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(failureResponse{Error: err.Error()})
		return
	}

	userID, err := database.Authenticate(request.Credentials.Username, request.Credentials.Password)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(failureResponse{Error: err.Error()})
		return
	}

	notes, err := database.GetUserNotes(userID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(failureResponse{Error: err.Error()})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(getResponse{Notes: notes})
}
