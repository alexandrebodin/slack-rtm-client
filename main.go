package main

import (
	"encoding/json"
	"fmt"
	ws "github.com/gorilla/websocket"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

type rtmstart struct {
	Url   string
	Users []struct {
		Id   string
		Name string
	}
}

type Message struct {
	Type    string
	Text    string
	User    string
	Channel string
	OK      string
	ReplyTo string `json:"reply_to"`
	Error   struct {
		Code string `json:"code"`
		Msg  string `json:"msg"`
	}
}

type Response struct {
	Type    string `json:"type"`
	Text    string `json:"text"`
	Id      string `json:"id"`
	Channel string `json:"channel"`
}

var messagesSent int64 = 1

func main() {
	slackAddr := "https://slack.com/api/rtm.start"
	token := os.Getenv("GILIBOT_SLACK_TOKEN")
	if token == "" {
		log.Fatal("slack token is missing")
	}

	resp, _ := http.Get(slackAddr + "?token=" + token + "&pretty=1")
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)

	responseData := &rtmstart{}
	json.Unmarshal(body, &responseData)

	socketUrl := responseData.Url

	var d *ws.Dialer
	c, _, err := d.Dial(socketUrl, *new(http.Header))
	if err != nil {
		log.Fatal(err)
	}

	for {
		_, p, err := c.ReadMessage()
		if err != nil {
			log.Fatal(err)
		}

		var y map[string]interface{}
		json.Unmarshal(p, &y)

		for k, v := range y {
			fmt.Println(k, " => ", v.(string))
		}
		fmt.Println("-------------------")

		message := &Message{}
		json.Unmarshal(p, &message)

		if message.Type == "message" && strings.ToLower(message.Text) == "hello kilibot" {

			response := &Response{
				Id:      fmt.Sprintf("%v", messagesSent),
				Type:    "message",
				Channel: message.Channel,
				Text:    "hello to you !",
			}
			json, err := json.Marshal(response)
			if err != nil {
				log.Fatal(err)
			}

			err = c.WriteMessage(ws.TextMessage, json)
			if err != nil {
				log.Fatal(err)
			}
			messagesSent++
		}
	}
}
