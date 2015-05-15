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

type SlackClient struct {
	slackData  SlackData
	dispatcher *slackDispatcher
	conn       *ws.Conn
}

func New(token string) (*SlackClient, error) {

	resp, err := http.Get(slackAddr + "?token=" + token)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	s := &SlackClient{
		dispatcher: &slackDispatcher{},
	}
	err = json.Unmarshal(body, &s.slackData)
	if err != nil {
		return nil, err
	}

	return s, nil
}

func (s *SlackClient) Run(h ...HelloHandler) error {

	conn, _, err := ws.DefaultDialer.Dial(s.slackData.Url, nil)
	if err != nil {
		return err
	}

	s.conn = conn

	if len(h) > 0 {
		s.dispatcher.addHelloListener(h[0])
	}

	s.startReader()
	return nil
}

func (s *SlackClient) AddListener(eType EventType, v interface{}) {

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

func (s *SlackClient) startReader() {

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

func (s *SlackClient) WriteMessage(v interface{}) error {
	return s.conn.WriteJSON(v)
}
