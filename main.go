package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"log"
	"net/http"
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
	// Initialize router
	const (
		host     = "localhost"
		port     = 5432
		user     = ""
		password = ""
		dbname   = "wedding"
	)
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	// Connect to the PostgreSQL database
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Verify the connection
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Successfully connected to PostgreSQL!")

	router := mux.NewRouter()

	router.HandleFunc("/rsvp", createRSVP(db)).Methods("POST")

	// Start server
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
		err := db.QueryRow(sqlStatement, rsvp.FullName).Scan(&id)
		if err != nil {
			log.Fatal(err)
		}
		_ = json.NewEncoder(writer).Encode(id)
	}
}
