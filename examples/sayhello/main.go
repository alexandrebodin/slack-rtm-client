package main

import (
	slack "github.com/alexandrebodin/slack_rtm"
	"log"
	"os"
)

type handler struct{}

func (h *handler) OnHello() error {

	log.Println("connection established")
	return nil
}

func (h *handler) OnMessage(m *slack.MessageType) error {

	log.Println("message received")
	log.Println(m)
	return nil
}

func main() {

	token := os.Getenv("GILIBOT_SLACK_TOKEN")
	if token == "" {
		log.Fatal("slack token is missing")
	}

	h := &handler{}

	slackClient, err := slack.New(token)
	if err != nil {
		log.Fatal(err)
	}

	slackClient.AddListener(slack.MessageEvent, h)

	err = slackClient.Run(h)
	if err != nil {
		log.Fatal(err)
	}
}
