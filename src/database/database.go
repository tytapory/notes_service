package database

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

type user struct {
	UserID       int
	Username     string
	PasswordHash string
}

type userNote struct {
	NoteID   int
	UserID   int
	NoteText string
}

var (
	dbInstance *sql.DB
)

func ConnectToDatabase(dbUser, dbName, dbPassword string) error {
	connectionString := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable", dbUser, dbPassword, dbName)
	var err error
	dbInstance, err = sql.Open("postgres", connectionString)
	if err != nil {
		return err
	}

	if err = dbInstance.Ping(); err != nil {
		return err
	}
	return nil
}

func GetDatabase() *sql.DB {
	return dbInstance
}

func RegisterUser(username, password string) error {
	var exists bool
	err := GetDatabase().QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE username = $1)", username).Scan(&exists)
	if err != nil {
		return err
	}
	if exists {
		return err
	}

	hashedPassword, err := hashPassword(password)
	if err != nil {
		return fmt.Errorf("failed to hash password: %v", err)
	}

	query := `INSERT INTO users (username, password_hash) VALUES ($1, $2)`

	_, err = GetDatabase().Exec(query, username, hashedPassword)
	if err != nil {
		return fmt.Errorf("error registering user: %v", err)
	}

	return nil
}

func GetUserNotes(UserID int) ([]string, error) {
	query := `SELECT note_text FROM user_notes WHERE user_id = $1`
	rows, err := GetDatabase().Query(query, UserID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var notes []string
	for rows.Next() {
		var noteText string
		if err = rows.Scan(&noteText); err != nil {
			return nil, err
		}
		notes = append(notes, noteText)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}
	return notes, nil
}

func Authenticate(username, password string) (int, error) {
	hashedPassword, err := getHashedPassword(username)
	if err != nil {
		return -1, err
	}
	if checkPasswordHash(password, hashedPassword) {
		userID, err := getUserIDByUsername(username)
		return userID, err
	} else {
		return -1, fmt.Errorf("incorrect credentials")
	}
}

func getHashedPassword(username string) (string, error) {
	query := `SELECT password_hash FROM users WHERE username = $1;`
	var passwordHash string
	err := GetDatabase().QueryRow(query, username).Scan(&passwordHash)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", fmt.Errorf("user does not exists: %s", username)
		}
		return "", err
	}

	return passwordHash, nil

}

func InsertUserNote(userID int, noteText string) error {
	query := `INSERT INTO user_notes (user_id, note_text) VALUES ($1, $2)`

	_, err := GetDatabase().Exec(query, userID, noteText)
	if err != nil {
		return fmt.Errorf("error inserting user note: %v", err)
	}

	return nil
}

func getUserIDByUsername(username string) (int, error) {
	query := `SELECT user_id FROM users WHERE username = $1`

	var userID int
	err := GetDatabase().QueryRow(query, username).Scan(&userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, fmt.Errorf("no user found with username: %s", username)
		}
		return 0, fmt.Errorf("error querying user_id: %v", err)
	}

	return userID, nil
}

func hashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

func checkPasswordHash(password, hashedPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}
