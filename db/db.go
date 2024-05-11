package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"
)

func ConnectToDb() *sql.DB {
	host := os.Getenv("POSTGRES_HOST")
	port, err := strconv.Atoi(os.Getenv("POSTGRES_PORT"))
	user := os.Getenv("POSTGRES_USER")
	password := os.Getenv("POSTGRES_SECRET")
	dbname := os.Getenv("POSTGRES_WDST_DB")

	if err != nil {
		log.Println("Error parsing POSTGRES_PORT")
		log.Fatal(err)
	}
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatal(err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Successfully connected to PostgreSQL on port " + os.Getenv("POSTGRES_PORT"))
	return db
}
