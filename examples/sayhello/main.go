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
		eType, event, err := slackClient.NextEvent()
		if err != nil {
			log.Fatal(err)
		}

		if eType == slack.MessageEvent {
			log.Println("Message received")
			log.Println(event)
		}
	}
}
