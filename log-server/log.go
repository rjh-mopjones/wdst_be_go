package log_server

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
	"wdst_be/email"
)

func OpenLogFile(logFilename string) *os.File {
	logFile, err := os.OpenFile(logFilename, os.O_APPEND|os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		log.Panic(err)
	}
	return logFile
}

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

func RefreshLog() func(http.ResponseWriter, *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		logfile := os.Getenv("WDST_LOG_FILE")
		currentTime := strings.ToLower(time.Now().Format("Mon-Jan-2-2006"))
		subject := "Logs for " + currentTime
		log.Println("Sending " + subject)
		newLogFile := logfile + "-" + currentTime + ".log"

		e := os.Rename(logfile+".log", newLogFile)
		if e != nil {
			log.Fatal(e)
		}

		applicationLogFile := OpenLogFile(logfile + ".log")
		mw := io.MultiWriter(os.Stdout, applicationLogFile)
		log.SetOutput(mw)
		email.SendEmail("roryhedderman@gmail.com", "logs for "+currentTime+" attached",
			"Logs for "+currentTime, newLogFile)

		writer.WriteHeader(http.StatusOK)
		writer.Header().Set("Content-Type", "application/json")
		json.NewEncoder(writer).Encode("refreshed logs")
	}
}
