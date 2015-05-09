package slack_rtm

type Error struct {
	Code string `json:"code"`
	Msg  string `json:"msg"`
}

type ResponseMessage struct {
	Type    string `json:"type"`
	Text    string `json:"text"`
	Id      string `json:"id"`
	Channel string `json:"channel"`
}

type AbstractEvent struct {
	Type    string `json="type"`
	SubType string `json="subtype"`
}

type HelloType struct {
	Type string `json="type"`
}

type MessageType struct {
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

type UserTypingType struct {
	Type    string `json="type"`
	Channel string `json="channel"`
	User    string `json="user"`
}

type HelloHandler interface {
	OnHello(c *SlackContext) error
}

type MessageHandler interface {
	OnMessage(c *SlackContext, m *MessageType) error
}

type TextMessageHandler func(m *MessageType) error

type EventType int

const (
	NoEvent EventType = iota
	HelloEvent
	MessageEvent
	TextMessageEvent
	UserTypingEvent
	ChannelMarkedEvent
	ChannelCreatedEvent
	ChannelJoinedEvent
	ChannelLeftEvent
	ChannelDeletedEvent
	ChannelRenameEvent
	ChannelArchiveEvent
	ChannelUnarchiveEvent
	ChannelHistoryChangedEvent
	ImCreatedEvent
	ImOpenEvent
	ImCloseEvent
	ImMarkedEvent
	ImHistoryChangedEvent
	GroupJoinedEvent
	GroupLeftEvent
	GroupOpenEvent
	GroupCloseEvent
	GroupArchiveEvent
	GroupUnarchiveEvent
	GroupRenameEvent
	GroupMarkedEvent
	GroupHistoryChangedEvent
	FileCreatedEvent
	FileSharedEvent
	FileUnsharedEvent
	FilePublicEvent
	FilePrivateEvent
	FileChangeEvent
	FileDeletedEvent
	FileCommentAddedEvent
	FileCommentEditedEvent
	FileCommentDeletedEvent
	PinAddedEvent
	PinRemovedEvent
	PresenceChangedEvent
	ManualPresenceChangeEvent
	PrefChangeEvent
	UserChangeEvent
	TeamJoinEvent
	StarAddedEvent
	StarRemovedEvent
	EmojiChangedEvent
	CommandsChangedEvent
	TeamPlanChangeEvent
	TeamPrefChangeEvent
	TeamRenameEvent
	TeamDomainChangeEvent
	EmailDomainChangedEvent
	BotAddedEvent
	BotChangedEvent
	AccountsChangedEvent
	TeamMigrationStartedEvent
)
