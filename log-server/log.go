package log_server

import (
	"bytes"
	"log"
	"net/http"
)

func HandleLog() func(writer http.ResponseWriter, request *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		buf := new(bytes.Buffer)
		buf.ReadFrom(request.Body)
		log.Println("Incoming Log: " + request.RemoteAddr + " :- " + buf.String())
	}
}
