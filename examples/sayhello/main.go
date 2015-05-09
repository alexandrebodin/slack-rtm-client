package main

import (
	slack "github.com/alexandrebodin/slack_rtm"
	"log"
	"os"
)

type handler struct{}

func (h *handler) OnHello(c *slack.SlackContext) error {

	log.Println("connection established")
	return nil
}

func (h *handler) OnMessage(c *slack.SlackContext, m *slack.MessageType) error {

	r := slack.ResponseMessage{
		Id:      "1",
		Type:    "message",
		Text:    "Coucou les copains",
		Channel: m.Channel,
	}
	c.Client.WriteMessage(r)
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
