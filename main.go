package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"os"
	"strconv"
)

type dtoRsvp struct {
	FullName       string `json:"fullName"`
	Email          string `json:"email"`
	Starter        string `json:"starter"`
	Main           string `json:"main"`
	Dessert        string `json:"dessert"`
	Song           string `json:"song"`
	Message        string `json:"message"`
	Diet           string `json:"diet"`
	Attendance     bool   `json:"attendance"`
	AdditionalRSVP []struct {
		FullName   string `json:"fullName"`
		Attendance bool   `json:"attendance"`
		Diet       string `json:"diet"`
		Starter    string `json:"starter"`
		Main       string `json:"main"`
		Dessert    string `json:"dessert"`
	} `json:"additionalRSVP"`
}

func main() {
	logFile := openLogFile("wdst_be.log")
	log.SetOutput(logFile)
	defer logFile.Close()
	db := connectToDb()
	defer db.Close()

	router := mux.NewRouter()

	router.HandleFunc("/rsvp", createRSVP(db)).Methods("POST")

	log.Println("Listening on port 8000")
	log.Fatal(http.ListenAndServe(":8000", router))
}

func createRSVP(db *sql.DB) func(writer http.ResponseWriter, request *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("Content-Type", "application/json")
		sqlStatement := "INSERT INTO rsvp (first_name) VALUES ($1) RETURNING id"
		var rsvp dtoRsvp
		var id int
		_ = json.NewDecoder(request.Body).Decode(&rsvp)
		reqBody, _ := json.Marshal(rsvp)
		err := db.QueryRow(sqlStatement, rsvp.FullName).Scan(&id)
		log.Println("		ID: " + strconv.Itoa(id) + "         " + string(reqBody))

		if err != nil {
			log.Fatal(err)
		}
		json.NewEncoder(writer).Encode(id)
	}
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
	fmt.Println("Successfully connected to PostgreSQL!")
	return db
}

func openLogFile(logFilename string) *os.File {
	logFile, err := os.OpenFile(logFilename, os.O_APPEND|os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		log.Panic(err)
	}
	return logFile
}
