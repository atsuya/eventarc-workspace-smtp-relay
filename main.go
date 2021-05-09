package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"net/smtp"
	"os"

	"github.com/atsuya/eventarc-workspace-smtp-relay/auth"
)

func SendEmail(subject string, message string) {
	smtphost := "smtp-relay.gmail.com"
	smtpport := 587
	hellohost := ""
	username := ""
	password := os.Args[1]
	from := ""
	to := ""

	smtpserver := fmt.Sprintf("%s:%d", smtphost, smtpport)
	client, err := smtp.Dial(smtpserver)
	if err != nil {
		log.Panic(err)
	}

	if err = client.Hello(hellohost); err != nil {
		log.Panic(err)
	}

	tlsconfig := &tls.Config{
		InsecureSkipVerify: false,
		ServerName:         smtphost,
	}
	client.StartTLS(tlsconfig)

	authlogin := auth.AuthLogin(username, password)
	if err = client.Auth(authlogin); err != nil {
		log.Panic(err)
	}

	if err = client.Mail(from); err != nil {
		log.Panic(err)
	}

	if err = client.Rcpt(to); err != nil {
		log.Panic(err)
	}

	writer, err := client.Data()
	if err != nil {
		log.Panic(err)
	}

	data := fmt.Sprintf("To: %s\r\nSubject: %s\r\n\r\n%s", to, subject, message)
	_, err = writer.Write([]byte(data))
	if err != nil {
		log.Panic(err)
	}

	err = writer.Close()
	if err != nil {
		log.Panic(err)
	}

	client.Quit()
}

func HandleEventarc(w http.ResponseWriter, r *http.Request) {
	s := fmt.Sprintf("Detected change in Cloud Storage bucket: %s", string(r.Header.Get("Ce-Subject")))
	log.Printf(s)

	SendEmail("yo watch out", "read this shit")

	fmt.Fprintln(w, s)
}

func main() {

	http.HandleFunc("/", HandleEventarc)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Listening on port %s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}
