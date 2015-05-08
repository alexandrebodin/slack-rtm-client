package slack_rtm

type slackDispatcher struct {
	helloListeners   []HelloHandler
	messageListeners []MessageHandler
}

func (d *slackDispatcher) addMessageListener(handler MessageHandler) {

	d.messageListeners = append(d.messageListeners, handler)
}

func (d *slackDispatcher) dispatchMessage(data *MessageType) {

	for _, listener := range d.messageListeners {
		listener.OnMessage(data)
	}
}

func (d *slackDispatcher) addHelloListener(handler HelloHandler) {

	d.helloListeners = append(d.helloListeners, handler)
}

func (d *slackDispatcher) dispatchHello() {

	for _, listener := range d.helloListeners {
		listener.OnHello()
	}
}
