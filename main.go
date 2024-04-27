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

	log.Println("App launched and listening on port 8000")
	fmt.Println("App launched and listening on port 8000")
	log.Fatal(http.ListenAndServe(":8000", router))
}

func createRSVP(db *sql.DB) func(writer http.ResponseWriter, request *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("Content-Type", "application/json")

		var returnIds []int
		sqlStatement := "INSERT INTO rsvp (full_name, email, " +
			"dinner_starter, dinner_main, dinner_dessert, " +
			"song, message, dietary_requirements, attendance) " +
			"VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING id"

		var rsvp dtoRsvp
		_ = json.NewDecoder(request.Body).Decode(&rsvp)
		for _, addRsvp := range rsvp.AdditionalRSVP {
			var addId int
			err := db.QueryRow(sqlStatement, addRsvp.FullName, rsvp.Email, addRsvp.Starter,
				addRsvp.Main, addRsvp.Dessert, "", "",
				addRsvp.Diet, addRsvp.Attendance).Scan(&addId)
			if err != nil {
				log.Fatal(err)
			}
			log.Println(" ID: " + strconv.Itoa(addId) + ",  " + "Processed third party RSVP for " +
				addRsvp.FullName + " of " + strconv.FormatBool(addRsvp.Attendance) + " by " + rsvp.FullName)
			returnIds = append(returnIds, addId)
		}

		var id int
		err := db.QueryRow(sqlStatement, rsvp.FullName, rsvp.Email, rsvp.Starter,
			rsvp.Main, rsvp.Dessert, rsvp.Song, rsvp.Message,
			rsvp.Diet, rsvp.Attendance).Scan(&id)
		log.Println(" ID: " + strconv.Itoa(id) + ",  " + "Processed RSVP for " +
			rsvp.FullName + " of " + strconv.FormatBool(rsvp.Attendance))
		returnIds = append(returnIds, id)

		if err != nil {
			log.Fatal(err)
		}
		writer.WriteHeader(http.StatusOK)
		json.NewEncoder(writer).Encode("OK")
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
