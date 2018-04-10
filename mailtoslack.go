package main

import (
	"bytes"
	"fmt"
	"log"
	"net"
	"net/mail"
	"os"
	"strings"

	"github.com/mhale/smtpd"
	"github.com/nlopes/slack"
	"github.com/veqryn/go-email/email"
)

var (
	port         = os.Getenv("PORT")
	slackToken   = os.Getenv("SLACK_TOKEN")
	slackChannel = os.Getenv("SLACK_CHANNEL")
	domainList   = os.Getenv("DOMAIN_LIST")
	filetype     string
)

func mailHandler(origin net.Addr, from string, to []string, data []byte) {
	msg, _ := email.ParseMessage(bytes.NewReader(data))
	subject := msg.Header.Get("Subject")
	sender := msg.Header.Get("From")
	recipient := msg.Header.Get("To")

	// If we have been given a list of recipient domains, filter on these
	if len(domainList) > 0 {
		domains := strings.Split(domainList, ",")
		rcpt, _ := mail.ParseAddress(recipient)
		recipientDomain := strings.Split(rcpt.Address, "@")[1]
		ok := false
		for i := 0; i < len(domains); i++ {
			if recipientDomain == domains[i] {
				ok = true
				break
			}
		}
		if !ok {
			return
		}
	}
	api := slack.New(slackToken)
	for _, part := range msg.MessagesContentTypePrefix("text/plain") {
		uploadparams := slack.FileUploadParameters{
			Channels:       []string{slackChannel},
			Title:          fmt.Sprintf("Subject: %s", subject),
			Filetype:       "text",
			Content:        string(part.Body),
			InitialComment: fmt.Sprintf("To: %s\nFrom: %s", recipient, sender),
		}
		file, err := api.UploadFile(uploadparams)
		if err != nil {
			log.Fatal(err)
			return
		}
		log.Printf("Message successfully sent to channel %s as text file %s", slackChannel, file.Name)
	}
	for _, part := range msg.MessagesContentTypePrefix("text/html") {
		uploadparams := slack.FileUploadParameters{
			Channels:       []string{slackChannel},
			Title:          fmt.Sprintf("Subject: %s", subject),
			Filetype:       "html",
			Content:        string(part.Body),
			InitialComment: fmt.Sprintf("To: %s\nFrom: %s", recipient, sender),
		}
		file, err := api.UploadFile(uploadparams)
		if err != nil {
			log.Fatal(err)
			return
		}
		log.Printf("Message successfully sent to channel %s as HTML file %s", slackChannel, file.Name)
	}
}

func main() {
	if port == "" {
		port = "2525"
		log.Printf("No PORT found, defaulting to %s", port)
	}
	if slackToken == "" {
		log.Fatal("No SLACK_TOKEN found")
		return
	}
	if slackChannel == "" {
		log.Fatal("No SLACK_CHANNEL found")
		return
	}
	log.Printf("Listening for mail on port %s", port)
	smtpd.ListenAndServe(fmt.Sprintf("0.0.0.0:%s", port), mailHandler, "Sendmail 8.11.3", "")
}
