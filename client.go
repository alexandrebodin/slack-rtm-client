package slack_rtm

import (
	"encoding/json"
	"errors"
	ws "github.com/gorilla/websocket"
	"io/ioutil"
	"net/http"
	"reflect"
)

type SlackData struct {
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

type SlackMessage struct {
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

var slackAddr = "https://slack.com/api/rtm.start"

var (
	errInvalidEvent = errors.New("slackClient: message received but no type specified")
	errTypeNotFound = errors.New("slackClient: message received but type unrecognized")
)

type SlackClient struct {
	SlackData SlackData
	conn      *ws.Conn
}

func NewSlackClient(token string) (*SlackClient, error) {

	resp, err := http.Get(slackAddr + "?token=" + token)
	defer resp.Body.Close()
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	s := &SlackClient{}
	err = json.Unmarshal(body, &s.SlackData)
	if err != nil {
		return nil, err
	}

	d := ws.DefaultDialer
	conn, _, err := d.Dial(s.SlackData.Url, nil)
	if err != nil {
		return nil, err
	}

	s.conn = conn
	return s, nil
}

func (s *SlackClient) NextEvent() (Event, error) {

	_, data, err := s.conn.ReadMessage()
	if err != nil {
		return nil, err
	}

	//get type and subtype for primary detection
	var event AbstractEvent
	err = json.Unmarshal(data, &event)
	if err != nil {
		return nil, err
	}

	//once concrete event initialzed, fill it
	concreteEvent, err := retrieveEvent(event)
	if err != nil {
		return nil, err
	}

	t := reflect.TypeOf(concreteEvent)
	ref := reflect.New(t).Interface()

	err = json.Unmarshal(data, &ref)
	if err != nil {
		return nil, err
	}

	return ref, nil
}

func retrieveEvent(e AbstractEvent) (Event, error) {

	switch e.Type {
	case "hello":
		return HelloEvent{}, nil
	case "message":
		if e.SubType != "" {
			return detectSubType(e)
		}
		return MessageEvent{}, nil
	case "user_typing":
		return UserTypingEvent{}, nil
	default:
		return nil, errTypeNotFound
	}
}

func detectSubType(e AbstractEvent) (Event, error) {

	switch e.SubType {
	default:
		return nil, errTypeNotFound
	}
}
