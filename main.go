package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/smtp"
	"os"

	"github.com/atsuya/google-workspace-smtp-relay/auth"
)

func main() {
	smtphost := "smtp-relay.gmail.com"
	smtpport := 587
	hellohost := ""
	username := ""
	password := os.Args[1]
	from := ""
	to := ""
	subject := ""
	body := ""

	smtpserver := fmt.Sprintf("%s:%d", smtphost, smtpport)
	client, err := smtp.Dial(smtpserver)
	if err != nil {
		log.Panic(err)
	}

	if err = client.Hello(hellohost); err != nil {
		log.Panic(err)
	}

	tlsconfig := &tls.Config{
		ServerName:         smtphost,
		InsecureSkipVerify: false,
		MinVersion:         tls.VersionTLS13,
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

	data := fmt.Sprintf("To: %s\r\nSubject: %s\r\n\r\n%s", to, subject, body)
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
