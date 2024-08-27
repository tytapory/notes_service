package main

import (
	"flag"
	"fmt"

	"notes_service/database"
	"notes_service/server"
)

func main() {
	dbUser := flag.String("dbuser", "notes_service", "Username for the database")
	dbName := flag.String("dbname", "notes", "Name of the database")
	dbPassword := flag.String("dbpassword", "password", "Password for the database")
	flag.Parse()

	err := database.ConnectToDatabase(*dbUser, *dbName, *dbPassword)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	err = server.Run()
	if err != nil {
		fmt.Println(err.Error())
	}
}
