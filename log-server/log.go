package log_server

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
)

func HandleLog() func(writer http.ResponseWriter, request *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		buf := new(bytes.Buffer)
		buf.ReadFrom(request.Body)
		log.Println("Incoming Log: " + request.RemoteAddr + " :- " + buf.String())
		writer.WriteHeader(http.StatusOK)
		writer.Header().Set("Content-Type", "application/json")
		json.NewEncoder(writer).Encode("logged")
	}
}
