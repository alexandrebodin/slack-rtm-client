package slack_rtm

type Event interface{}

type AbstractEvent struct {
	Type    string `json="type"`
	SubType string `json="subtype"`
}

type HelloEvent struct {
	Type string `json="type"`
}

type MessageEvent struct {
	Type    string `json="type"`
	SubType string `json="subtype"`
	Channel string `json="channel"`
	User    string `json="user"`
	Text    string `json="text"`
	Ts      string `json="ts"`
	Edited  struct {
		User string `json="user"`
		Ts   string `json="ts"`
	} `json="edited"`
}

type UserTypingEvent struct {
	Type    string `json="type"`
	Channel string `json="channel"`
	User    string `json="user"`
}
