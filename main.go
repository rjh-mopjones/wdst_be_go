package main

import (
	"fmt"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"os"
	"wdst_be/db"
	log_server "wdst_be/log-server"
	"wdst_be/rsvp"
)

func main() {
	applicationLogFile := openLogFile("wdst.log")
	log.SetOutput(applicationLogFile)
	defer applicationLogFile.Close()
	db := db.ConnectToDb()
	defer db.Close()

	router := mux.NewRouter()
	router.Use(enableCORS)
	router.HandleFunc("/rsvp", rsvp.HandleRSVP(db)).Methods("POST")
	router.HandleFunc("/log-server", log_server.HandleLog()).Methods("POST")

	log.Println("App launched and listening on port 8000")
	fmt.Println("App launched and listening on port 8000")
	log.Fatal(http.ListenAndServe(":8000", router))
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
