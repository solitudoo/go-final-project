package main

import (
	"fmt"
	"os"

	"go-final-project/pkg/db"
	"go-final-project/pkg/server"
)

func main() {
	dbFile := os.Getenv("TODO_DBFILE")
	if dbFile == "" {
		dbFile = "scheduler.db"
	}
	if err := db.Init(dbFile); err != nil {
		fmt.Println(err)
		return
	}

	if err := server.Run(); err != nil {
		fmt.Println(err)
		return
	}

}
