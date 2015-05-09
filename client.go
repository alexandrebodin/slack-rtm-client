package slack_rtm

import (
	"encoding/json"
	"errors"
	ws "github.com/gorilla/websocket"
	"io/ioutil"
	"net/http"
)

var slackAddr = "https://slack.com/api/rtm.start"

var (
	errInvalidEvent = errors.New("slackClient: message received but no type specified")
	errTypeNotFound = errors.New("slackClient: message received but type unrecognized")
)

type slackClient struct {
	slackData  SlackData
	dispatcher *slackDispatcher
	conn       *ws.Conn
}

func New(token string) (*slackClient, error) {

	resp, err := http.Get(slackAddr + "?token=" + token)
	defer resp.Body.Close()
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	s := &slackClient{
		dispatcher: &slackDispatcher{},
	}
	err = json.Unmarshal(body, &s.slackData)
	if err != nil {
		return nil, err
	}

	return s, nil
}

func (s *slackClient) Run(h HelloHandler) error {

	conn, _, err := ws.DefaultDialer.Dial(s.slackData.Url, nil)
	if err != nil {
		return err
	}

	s.conn = conn

	if h != nil {
		s.dispatcher.addHelloListener(h)
	}

	s.startReader()
	return nil
}

func (s *slackClient) AddListener(eType EventType, v interface{}) {

	switch eType {
	case HelloEvent:
		s.dispatcher.addHelloListener(v.(HelloHandler))
	case MessageEvent:
		s.dispatcher.addMessageListener(v.(MessageHandler))
	}

}

type Event struct {
	t    EventType
	data interface{}
}

func (s *slackClient) startReader() {

	for {
		_, data, err := s.conn.ReadMessage()
		if err != nil {
			continue
		}

		var event AbstractEvent
		err = json.Unmarshal(data, &event)
		if err != nil {
			continue
		}

		ctx := &SlackContext{s}

		d := s.dispatcher
		switch event.Type {
		case "hello":
			d.dispatchHello(ctx)
		case "message":
			m := &MessageType{}
			json.Unmarshal(data, &m)
			d.dispatchMessage(ctx, m)

		default:
			continue
		}
	}
}

func (s *slackClient) WriteMessage(v interface{}) error {
	return s.conn.WriteJSON(v)
}
