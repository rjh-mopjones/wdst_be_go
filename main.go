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
	logserver "wdst_be/log-server"
	"wdst_be/rsvp"
)

func main() {
	applicationLogFile := logserver.OpenLogFile(os.Getenv("WDST_LOG_FILE") + ".log")
	mw := io.MultiWriter(os.Stdout, applicationLogFile)
	log.SetOutput(mw)
	defer applicationLogFile.Close()
	db := db.ConnectToDb()
	defer db.Close()
	port := ":8000"

	router := mux.NewRouter()
	router.HandleFunc("/rsvp", rsvp.HandleRSVP(db)).Methods("POST")
	router.HandleFunc("/log-server", logserver.HandleLog()).Methods("POST")
	router.HandleFunc("/log-server", logserver.RefreshLog()).Methods("GET")

	handler := cors.AllowAll().Handler(router)
	log.Println("App launched and listening on port " + port)
	log.Fatal(http.ListenAndServe(port, handler))
}
