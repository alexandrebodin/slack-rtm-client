package slack_rtm

type slackDispatcher struct {
	helloListeners   []HelloHandler
	messageListeners []MessageHandler
}

type SlackContext struct {
	Client *slackClient
}

func (d *slackDispatcher) addMessageListener(handler MessageHandler) {

	d.messageListeners = append(d.messageListeners, handler)
}

func (d *slackDispatcher) dispatchMessage(c *SlackContext, m *MessageType) {

	for _, listener := range d.messageListeners {
		listener.OnMessage(c, m)
	}
}

func (d *slackDispatcher) addHelloListener(handler HelloHandler) {

	d.helloListeners = append(d.helloListeners, handler)
}

func (d *slackDispatcher) dispatchHello(c *SlackContext) {

	for _, listener := range d.helloListeners {
		listener.OnHello(c)
	}
}
