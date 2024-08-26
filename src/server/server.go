package server

import (
	"encoding/json"
	"log"
	"net/http"
	"notes_service/database"
	"os"
	"path/filepath"
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

func initLogger() {
	logDir := "logs"
	if err := os.MkdirAll(logDir, os.ModePerm); err != nil {
		log.Fatalf("ERROR: Can't create logs dir: %v", err)
	}

	logFile, err := os.OpenFile(filepath.Join(logDir, "app.log"), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("ERROR: Can't open log file: %v", err)
	}

	log.SetOutput(logFile)
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
}

func Run() error {
	plugInHandlers()
	initLogger()
	log.Println("INFO: Running server on :8080")
	return http.ListenAndServe(":8080", nil)
}

func plugInHandlers() {
	http.HandleFunc("/push_note", pushNote)
	http.HandleFunc("/get_notes", getNotes)
	http.HandleFunc("/register", register)
}

func register(w http.ResponseWriter, r *http.Request) {
	log.Println("INFO: Starting to process registration request")
	if r.Method != http.MethodPost {
		log.Printf("ERROR: Invalid request method: %s", r.Method)
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var request registerRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		log.Printf("ERROR: Error decoding request body: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(failureResponse{Error: err.Error()})
		return
	}

	log.Printf("INFO: Attempting to register user: %s", request.Credentials.Username)
	err = database.RegisterUser(request.Credentials.Username, request.Credentials.Password)
	if err != nil {
		log.Printf("ERROR: Error registering user: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(failureResponse{Error: err.Error()})
		return
	}

	log.Printf("INFO: User %s successfully registered", request.Credentials.Username)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(registerResponse{Message: "User successfully regustered"})
}

func pushNote(w http.ResponseWriter, r *http.Request) {
	log.Println("INFO: Starting to process pushNote request")
	if r.Method != http.MethodPost {
		log.Printf("ERROR: Invalid request method: %s", r.Method)
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var request pushRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		log.Printf("ERROR: Error decoding request body: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(failureResponse{Error: err.Error()})
		return
	}

	log.Printf("INFO: Attempting to authenticate user: %s", request.Credentials.Username)
	userID, err := database.Authenticate(request.Credentials.Username, request.Credentials.Password)
	if err != nil {
		log.Printf("ERROR: Authentication failed for user: %s, error: %v", request.Credentials.Username, err)
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(failureResponse{Error: err.Error()})
		return
	}

	log.Println("INFO: Sending note to yandex speller the note")
	validatedNote, err := validateAndFixErrors(request.Note)
	if err != nil {
		log.Printf("ERROR: Error validating note: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(failureResponse{Error: err.Error()})
		return
	}

	log.Printf("INFO: Inserting note for user ID: %d", userID)
	err = database.InsertUserNote(userID, validatedNote)
	if err != nil {
		log.Printf("ERROR: Error inserting note for user ID %d: %v", userID, err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(failureResponse{Error: err.Error()})
		return
	}

	log.Printf("INFO: Note successfully created for user ID: %d", userID)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(pushResponse{Message: "Note successfully created"})
}

func getNotes(w http.ResponseWriter, r *http.Request) {
	log.Println("INFO: Starting to process getNotes request")
	if r.Method != http.MethodPost {
		log.Printf("ERROR: Invalid request method: %s", r.Method)
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var request getRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		log.Printf("ERROR: Error decoding request body: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(failureResponse{Error: err.Error()})
		return
	}

	log.Printf("INFO: Attempting to authenticate user: %s", request.Credentials.Username)
	userID, err := database.Authenticate(request.Credentials.Username, request.Credentials.Password)
	if err != nil {
		log.Printf("ERROR: Authentication failed for user: %s, error: %v", request.Credentials.Username, err)
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(failureResponse{Error: err.Error()})
		return
	}

	log.Printf("INFO: Retrieving notes for user ID: %d", userID)
	notes, err := database.GetUserNotes(userID)
	if err != nil {
		log.Printf("ERROR: Error retrieving notes for user ID %d: %v", userID, err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(failureResponse{Error: err.Error()})
		return
	}

	log.Printf("INFO: Successfully retrieved notes for user ID: %d", userID)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(getResponse{Notes: notes})
}
