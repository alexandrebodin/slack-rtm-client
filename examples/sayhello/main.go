package main

import (
	slack "github.com/alexandrebodin/slack_rtm"
	"log"
	"os"
)

func main() {

	token := os.Getenv("GILIBOT_SLACK_TOKEN")
	if token == "" {
		log.Fatal("slack token is missing")
	}

	slackClient, err := slack.NewSlackClient(token)
	if err != nil {
		log.Fatal(err)
	}

	for {
		event, err := slackClient.NextEvent()
		if err != nil {
			log.Println(err)
			continue
		}

		switch event.(type) {
		case *slack.HelloEvent:
			log.Println("connection established")
		case *slack.MessageEvent:
			a := event.(*slack.MessageEvent)
			log.Printf("message received => %v", a.Text)
		case *slack.UserTypingEvent:
			a := event.(*slack.UserTypingEvent)
			log.Printf("User %v is typing in channel %v", a.User, a.Channel)
		}
	}
}
