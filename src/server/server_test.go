package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"notes_service/database"
	"testing"

	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

func setupDatabase(t *testing.T) {
	t.Helper()
	db := database.GetDatabase()
	tables := []string{"user_notes", "users"}
	for _, table := range tables {
		_, err := db.Exec("TRUNCATE TABLE " + table + " RESTART IDENTITY CASCADE;")
		if err != nil {
			t.Fatalf("Failed to truncate table %s: %v", table, err)
		}
	}
	_, err := db.Exec(`
INSERT INTO users (username, password_hash) VALUES
('testuser1', '$2a$10$ZmpzL3j.fEOgsfno.MHzNuYdkfQr5PRoUTWUbkJHhVvF6HMcwcwSW'), --password is "password"
('testuser2', '$2a$10$ZmpzL3j.fEOgsfno.MHzNuYdkfQr5PRoUTWUbkJHhVvF6HMcwcwSW'),
('testuser3', '$2a$10$ZmpzL3j.fEOgsfno.MHzNuYdkfQr5PRoUTWUbkJHhVvF6HMcwcwSW');

INSERT INTO user_notes (user_id, note_text) VALUES
(1, 'Test note 1 for user 1'),
(1, 'Test note 2 for user 1'),
(2, 'Test note 1 for user 2'),
(3, 'Test note 1 for user 3'),
(3, 'Test note 2 for user 3'),
(3, 'Test note 3 for user 3');`)
	if err != nil {
		t.Fatalf("Failed to execute SQL script: %v", err)
	}
}

func TestMain(m *testing.M) {
	err := database.ConnectToDatabase("mockdb", "mockdb", "password")
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	m.Run()
}

func TestRegisterHandler_Success(t *testing.T) {
	setupDatabase(t)
	payload := registerRequest{
		Credentials: credentials{
			Username: "newuser",
			Password: "newpassword",
		},
	}
	body, _ := json.Marshal(payload)
	req, err := http.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(body))
	assert.NoError(t, err)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(register)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusCreated, rr.Code)
	expected := `{"message":"User successfully registered"}`
	assert.JSONEq(t, expected, rr.Body.String())
}

func TestPushNoteHandler_Success(t *testing.T) {
	setupDatabase(t)
	payload := pushRequest{
		Credentials: credentials{
			Username: "testuser1",
			Password: "password",
		},
		Note: "This is a test note",
	}
	body, _ := json.Marshal(payload)
	req, err := http.NewRequest(http.MethodPost, "/push_note", bytes.NewBuffer(body))
	assert.NoError(t, err)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(pushNote)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusCreated, rr.Code)
	expected := `{"message":"Note successfully created"}`
	assert.JSONEq(t, expected, rr.Body.String())

	db := database.GetDatabase()

	var noteText string
	err = db.QueryRow(`
		SELECT note_text 
		FROM user_notes 
		WHERE user_id = (SELECT user_id FROM users WHERE username = 'testuser1')
		ORDER BY note_id DESC 
		LIMIT 1
	`).Scan(&noteText)

	assert.NoError(t, err)
	assert.Equal(t, "This is a test note", noteText)
}

func TestGetNotesHandler_Success(t *testing.T) {
	setupDatabase(t)
	payload := getRequest{
		Credentials: credentials{
			Username: "testuser1",
			Password: "password",
		},
	}
	body, _ := json.Marshal(payload)
	req, err := http.NewRequest(http.MethodPost, "/get_notes", bytes.NewBuffer(body))
	assert.NoError(t, err)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(getNotes)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	expected := `{"notes":["Test note 1 for user 1", "Test note 2 for user 1"]}`
	assert.JSONEq(t, expected, rr.Body.String())
}

func TestPushNoteHandler_AuthenticationError(t *testing.T) {
	setupDatabase(t)
	payload := pushRequest{
		Credentials: credentials{
			Username: "testuser1",
			Password: "incorrect",
		},
		Note: "This is a test note",
	}
	body, _ := json.Marshal(payload)
	req, err := http.NewRequest(http.MethodPost, "/push_note", bytes.NewBuffer(body))
	assert.NoError(t, err)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(pushNote)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusUnauthorized, rr.Code)
	expected := `{"error":"invalid credentials"}`
	assert.JSONEq(t, expected, rr.Body.String())
}

func TestGetNotesHandler_AuthenticationError(t *testing.T) {
	setupDatabase(t)
	payload := getRequest{
		Credentials: credentials{
			Username: "testuser1",
			Password: "incorrect",
		},
	}
	body, _ := json.Marshal(payload)
	req, err := http.NewRequest(http.MethodPost, "/get_notes", bytes.NewBuffer(body))
	assert.NoError(t, err)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(getNotes)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusUnauthorized, rr.Code)
	expected := `{"error":"invalid credentials"}`
	assert.JSONEq(t, expected, rr.Body.String())
}
