package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/smtp"
	"os"

	"github.com/atsuya/eventarc-workspace-smtp-relay/auth"
)

func main() {
	smtphost := "smtp-relay.gmail.com"
	smtpport := 587
	hellohost := ""
	username := ""
	password := os.Args[1]
	from := ""
	to := ""

	tlsconfig := &tls.Config{
		InsecureSkipVerify: false,
		ServerName:         smtphost,
	}

	authlogin := auth.AuthLogin(username, password)

	smtpserver := fmt.Sprintf("%s:%d", smtphost, smtpport)

	client, err := smtp.Dial(smtpserver)
	if err != nil {
		log.Panic(err)
	}

	if err = client.Hello(hellohost); err != nil {
		log.Panic(err)
	}

	client.StartTLS(tlsconfig)
	//connection, err := tls.Dial("tcp", smtpserver, tlsconfig)
	//if err != nil {
	//	log.Panic(err)
	//}

	//client, err := smtp.NewClient(connection, smtphost)
	//if err != nil {
	//	log.Panic(err)
	//}

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

	_, err = writer.Write([]byte("yo this is a test message"))
	if err != nil {
		log.Panic(err)
	}

	err = writer.Close()
	if err != nil {
		log.Panic(err)
	}

	client.Quit()
}
