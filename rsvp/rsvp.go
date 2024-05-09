package rsvp

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"net/smtp"
	"os"
	"strconv"
	"strings"
)

type dtoAdditionalRSVP struct {
	FullName   string `json:"fullName"`
	Attendance bool   `json:"attendance"`
	Diet       string `json:"diet"`
	Starter    string `json:"starter"`
	Main       string `json:"main"`
	Dessert    string `json:"dessert"`
	Email      string `json:"email"`
}

type dtoRsvp struct {
	FullName       string              `json:"fullName"`
	Email          string              `json:"email"`
	Starter        string              `json:"starter"`
	Main           string              `json:"main"`
	Dessert        string              `json:"dessert"`
	Song           string              `json:"song"`
	Message        string              `json:"message"`
	Diet           string              `json:"diet"`
	Attendance     bool                `json:"attendance"`
	AdditionalRSVP []dtoAdditionalRSVP `json:"additionalRSVP"`
}

func HandleRSVP(db *sql.DB) func(writer http.ResponseWriter, request *http.Request) {
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

			if addRsvp.Email == "" {
				addRsvp.Email = rsvp.Email
			}

			err := db.QueryRow(sqlStatement, addRsvp.FullName, addRsvp.Email, addRsvp.Starter,
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

		go receiptGenerator.sendEmail(rsvp)
		returnIds = append(returnIds, id)

		if err != nil {
			log.Fatal(err)
		}

		writer.WriteHeader(http.StatusOK)
		json.NewEncoder(writer).Encode("OK")
	}
}

func (rsvp dtoRsvp) sendEmail() bool {
	from := os.Getenv("EMAIL_USER")
	pass := os.Getenv("EMAIL_APP_SECRET")
	to := rsvp.Email

	body := rsvp.generateBody()
	msg := "From: " + from + "\n" +
		"To: " + to + "\n" +
		"Subject: Thank you for RSVPing!\n\n" + body

	err := smtp.SendMail("smtp.gmail.com:587",
		smtp.PlainAuth("", from, pass, "smtp.gmail.com"),
		from, []string{to}, []byte(msg))

	if err != nil {
		log.Printf("smtp error: %s", err)
		return false
	}

	log.Print("sent email")
	return true
}

func (rsvp dtoRsvp) generateBody() string {
	var sb strings.Builder
	sb.WriteString("Hello " + strings.Fields(rsvp.FullName)[0] + ", \n\n")
	sb.WriteString("We have received your RSVP and are really looking forward to seeing you on the 21st June 2025! " +
		"This is our personal email address so you can reply to this email with any questions. \n")
	sb.WriteString("If you would like to change any details about the RSVP, you can either re-rsvp via the " +
		"website (we will take the most recent one) or you can just reply to this email with any changes." + "\n")
	sb.WriteString("Below are the details we have received:- " + "\n\n\n")
	sb.WriteString("    Fullname:-                                  " + rsvp.FullName + "\n")
	sb.WriteString("    Starter:-                                      " + rsvp.Starter + "\n")
	sb.WriteString("    Main:-                                         " + rsvp.Main + "\n")
	sb.WriteString("    Dessert:-                                    " + rsvp.Dessert + "\n")
	sb.WriteString("    Song Request:-                          " + rsvp.Song + "\n")
	sb.WriteString("    Dietary Requirements:-              " + rsvp.Diet + "\n")
	sb.WriteString("    Message:-                          " + rsvp.Message + "\n\n\n")
	flag := true
	for _, addRsvp := range rsvp.AdditionalRSVP {
		if addRsvp.Email == rsvp.Email {
			if flag {
				sb.WriteString("We can also see that you have RSVP'd for:- " + "\n\n")
				flag = false
			} else {
				sb.WriteString("And: " + "\n\n")
			}
			sb.WriteString("    Fullname:-                                  " + addRsvp.FullName + "\n")
			sb.WriteString("    Starter:-                                      " + addRsvp.Starter + "\n")
			sb.WriteString("    Main:-                                         " + addRsvp.Main + "\n")
			sb.WriteString("    Dessert:-                                    " + addRsvp.Dessert + "\n")
			sb.WriteString("    Dietary Requirements:-              " + addRsvp.Diet + "\n\n")
		}
	}
	sb.WriteString("We look forward to seeing you!" + "\n\n")
	sb.WriteString("Kind Regards," + "\n")
	sb.WriteString("Mop and Ellie" + "\n")
	return sb.String()
}

type receiptGenerator interface {
	sendEmail() bool
}
