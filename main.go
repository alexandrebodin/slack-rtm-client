package main

import (
	"encoding/json"
	"fmt"
	ws "github.com/gorilla/websocket"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
	//	"strings"
)

type SlackBot struct {
	Ok    bool   `json="ok"`
	Error string `json="error"`
	Url   string `json="url"`
	Self  struct {
		ID   string `json="id"`
		Name string `json="name"`
	} `json="self"`
	Users    []User    `json="users"`
	Team     Team      `json="team"`
	Ims      []Im      `json="ims"`
	Groups   []Group   `json="groups"`
	Channels []Channel `json="channels"`
	Bots     []Bot     `json="bots"`
}

type User struct {
	Id                string `json="id"`
	Name              string `json="name"`
	Deleted           bool   `json="deleted"`
	Status            string `json="status"`
	Color             string `json="color"`
	RealName          string `json="real_name"`
	IsAdmin           bool   `json="is_admin"`
	IsOwner           bool   `json="is_owner"`
	IsPrimaryOwner    bool   `json="is_primary_owner"`
	IsRestricted      bool   `json="is_restricted"`
	IsUltraRestricted bool   `json="is_ultra_restricted"`
	IsBot             bool   `json="is_bot"`
	HasFiles          bool   `json="has_files"`
	Presence          string `json="presence"`
	TimeZone          string `json="tz"`
	TimeZoneLabel     string `json="tz_label"`
	TimeZoneOffset    string `json="tz_offset"`
	Profile           struct {
		Firstname          string `json="first_name"`
		Lastname           string `json="last_name"`
		RealName           string `json="real_name"`
		RealNameNormalized string `json="real_name_normalized"`
		Email              string `json="email"`
		Phone              string `json="phone"`
		Skype              string `json="skype"`
		Title              string `json="title"`
		Image24            string `json="image_24"`
		Image32            string `json="image_32"`
		Image48            string `json="image_48"`
		Image72            string `json="image_72"`
		Image192           string `json="image_192"`
		ImageOriginal      string `json="image_original"`
	} `json="profile"`
}

type Team struct {
	Id               string `json="id"`
	Name             string `json="name"`
	EmailDomain      string `json="email_domain"`
	Domain           string `json="domain"`
	OverStorageLimit bool   `json="over_storage_limit"`
	Plan             string `json="plan"`
}

type Group struct {
	Id                 string        `json="id"`
	Name               string        `json="name"`
	IsGroup            bool          `json="is_group"`
	Creator            string        `json="creator"`
	IsArchived         bool          `json="is_archived"`
	HasPins            bool          `json="has_pins"`
	IsOpen             bool          `json="is_open"`
	LastRead           string        `json="last_read"`
	Latest             LatestMessage `json="latest"`
	UnReadCount        int           `json="unread_count"`
	UnreadCountDisplay int           `json="unread_count_display"`
	MemberIds          []string      `json="members"`
	Purpose            Purpose       `json="purpose"`
	Topic              Topic         `json="topic"`
}

type Purpose struct {
	Value   string `json="value"`
	Creator string `json="creator"`
}

type Topic struct {
	Value   string `json="value"`
	Creator string `json="creator"`
}

type LatestMessage struct {
	Type      string `json="type"`
	User      string `json="user"`
	Text      string `json="text"`
	TimeStamp string `json="ts"`
}

type Channel struct {
	Id                 string        `json="id"`
	Name               string        `json="name"`
	Created            int           `json="created"`
	Creator            string        `json="creator"`
	IsArchived         bool          `json="is_archived"`
	IsGeneral          bool          `json="is_general"`
	HasPins            bool          `json="has_pins"`
	IsMember           bool          `json="is_member"`
	Latest             LatestMessage `json="latest"`
	MemberIds          []string      `json="members"`
	UnReadCount        int           `json="unread_count"`
	UnreadCountDisplay int           `json="unread_count_display"`
	Purpose            Purpose       `json="purpose"`
	Topic              Topic         `json="topic"`
}

type Im struct {
	Id                 string        `json="id"`
	IsIm               bool          `json="is_im"`
	User               string        `json="user"`
	LastRead           string        `json="last_read"`
	UnreadCount        int           `json="unread_count"`
	UnreadCountDisplay int           `json="unread_count_display"`
	IsOpen             bool          `json="is_open"`
	IsUserDeleted      bool          `json="is_user_deleted"`
	Latest             LatestMessage `json="latest"`
}

type Bot struct {
	Id      string `json="id"`
	Name    string `json="name"`
	Deleted bool   `json="deleted"`
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
var slackAddr = "https://slack.com/api/rtm.start"

func NewSlackBot(token string) *SlackBot {

	resp, err := http.Get(slackAddr + "?token=" + token)
	defer resp.Body.Close()
	if err != nil {
		log.Fatal(err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	slackBot := &SlackBot{}
	err = json.Unmarshal(body, &slackBot)
	if err != nil {
		log.Fatal(err)
	}

	return slackBot
}

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

func main() {

	token := os.Getenv("GILIBOT_SLACK_TOKEN")
	if token == "" {
		log.Fatal("slack token is missing")
	}

	slackBot := NewSlackBot(token)

	socketUrl := slackBot.Url

	d := ws.DefaultDialer
	con, _, err := d.Dial(socketUrl, nil)
	if err != nil {
		log.Fatal(err)
	}

	//defer con.Close()

	con.SetReadLimit(maxMessageSize)
	con.SetReadDeadline(time.Now().Add(pongWait))
	var message map[string]interface{}

	i := 1
	for {

		err := con.ReadJSON(&message)
		if err != nil {
			log.Fatal(err)
		}
		log.Println("message received")
		log.Println(message)
		//if no type set ignore
		if _, ok := message["type"]; !ok {
			continue
		}

		if mType, ok := message["type"].(string); ok {
			switch mType {
			case "message":
				con.WriteJSON(map[string]interface{}{
					"id":      fmt.Sprintf("%v", i),
					"type":    "message",
					"channel": message["channel"].(string),
					"text":    "hello to you !",
				})
				i++
			}
		}
	}
}
