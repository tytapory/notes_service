package main

import (
	"fmt"

	"notes_service/server"
)

func main() {
	err := server.Run()
	if err != nil {
		fmt.Print(err.Error())
	}
}
