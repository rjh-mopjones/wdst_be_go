package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

type RSVP struct {
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
	router := mux.NewRouter()

	// Route handlers
	router.HandleFunc("/rsvp", createRSVP).Methods("POST")

	// Start server
	log.Println("Listening on port 8000")
	log.Fatal(http.ListenAndServe(":8000", router))
}

func createRSVP(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "application/json")
	var rsvp RSVP
	_ = json.NewDecoder(request.Body).Decode(&rsvp)
	log.Println(rsvp)
	json.NewEncoder(writer).Encode(rsvp)
}
