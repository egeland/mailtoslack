package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/mail"
	"os"
	"strings"

	"github.com/mhale/smtpd"
	"github.com/nlopes/slack"
)

var (
	port         = os.Getenv("PORT")
	slackToken   = os.Getenv("SLACK_TOKEN")
	slackChannel = os.Getenv("SLACK_CHANNEL")
	domainList   = os.Getenv("DOMAIN_LIST")
)

func mailHandler(origin net.Addr, from string, to []string, data []byte) {
	msg, _ := mail.ReadMessage(bytes.NewReader(data))
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

	body, err := ioutil.ReadAll(msg.Body)
	if err != nil {
		log.Fatal(err)
		return
	}
	api := slack.New(slackToken)
	postData := slack.PostMessageParameters{}
	attachment := slack.Attachment{
		Text:       string(body),
		Title:      fmt.Sprintf("Subject: %s", subject),
		Fallback:   subject,
		AuthorName: fmt.Sprintf("To: %s", recipient),
		Footer:     fmt.Sprintf("From: %s", sender),
	}
	postData.Username = "Mailbot"
	postData.IconEmoji = ":email:"
	postData.Attachments = []slack.Attachment{attachment}
	channelID, timestamp, err := api.PostMessage(slackChannel, "", postData)
	if err != nil {
		log.Fatal(err)
		return
	}
	log.Printf("Message successfully sent to channel %s at %s", channelID, timestamp)
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
