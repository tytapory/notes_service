package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"notes_service/database"
	"os"
	"testing"
	"time"

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
	scriptPath := "../../sql/insert_mock_data.sql"
	script, err := os.ReadFile(scriptPath)
	if err != nil {
		t.Fatalf("Failed to read SQL script: %v", err)
	}
	_, err = db.Exec(string(script))
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

	go Run()

	time.Sleep(5 * time.Second)

	m.Run()
}

func TestPostNoteHandler_Success(t *testing.T) {
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
