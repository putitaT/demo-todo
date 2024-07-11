package database

import (
	"fmt"
	"log"
)

func CreateTable() {
	db := ConnectDB()
	createTb := `
	CREATE TABLE IF NOT EXISTS todos (
			id SERIAL PRIMARY KEY,
			title TEXT,
			status TEXT
	);
	`
	_, err := db.Exec(createTb)

	if err != nil {
		log.Fatal("can't create table", err)
	}

	fmt.Println("create table success")
}
