package main

import (
	"database/sql"
	"fmt"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"os"
	"strconv"
	"wdst_be/rsvp"
)

func main() {
	logFile := openLogFile("wdst_be.log")
	log.SetOutput(logFile)
	defer logFile.Close()
	db := connectToDb()
	defer db.Close()

	router := mux.NewRouter()
	router.Use(enableCORS)
	router.HandleFunc("/rsvp", rsvp.HandleRSVP(db)).Methods("POST")

	log.Println("App launched and listening on port 8000")
	fmt.Println("App launched and listening on port 8000")
	log.Fatal(http.ListenAndServe(":8000", router))
}

func connectToDb() *sql.DB {
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

func openLogFile(logFilename string) *os.File {
	logFile, err := os.OpenFile(logFilename, os.O_APPEND|os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		log.Panic(err)
	}
	return logFile
}

func enableCORS(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// Allow requests from any origin

		w.Header().Set("Access-Control-Allow-Origin", "<origin> | homeDomain")
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		// Allow specified HTTP methods

		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Private-Network", "true")

		// Allow specified headers

		w.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Accept")

		// Continue with the next handler

		next.ServeHTTP(w, r)
	})
}
