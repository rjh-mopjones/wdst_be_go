package main

import (
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"github.com/rs/cors"
	"io"
	"log"
	"net/http"
	"os"
	"wdst_be/db"
	log_server "wdst_be/log-server"
	"wdst_be/rsvp"
)

func main() {
	applicationLogFile := openLogFile(os.Getenv("WDST_LOG_FILE"))
	mw := io.MultiWriter(os.Stdout, applicationLogFile)
	log.SetOutput(mw)
	defer applicationLogFile.Close()
	db := db.ConnectToDb()
	defer db.Close()
	port := ":8000"

	router := mux.NewRouter()
	router.HandleFunc("/rsvp", rsvp.HandleRSVP(db)).Methods("POST")
	router.HandleFunc("/log-server", log_server.HandleLog()).Methods("POST")

	handler := cors.AllowAll().Handler(router)
	log.Println("App launched and listening on port " + port)
	log.Fatal(http.ListenAndServe(port, handler))
}

func openLogFile(logFilename string) *os.File {
	logFile, err := os.OpenFile(logFilename, os.O_APPEND|os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		log.Panic(err)
	}
	return logFile
}
