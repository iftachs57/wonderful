package consts

import (
	"time"

)

const (
	MaxMessageSize    = 8192
	ReadBufferSize    = 1024
	WriteBufferSize   = 1024
	Conf_Location 	  = "{configuration file location}"
	URL_RawQuery 	  = "model="
	Respons_done      = "response.done"
	Conversation_item_create = "conversation.item.create"
	Respons_create	  = "response.create"
)

var (
	HandshakeTimeout = 5 * time.Second
	PongWait   = 60 * time.Second
	WriteWait  = 10 * time.Second
	PingPeriod = 50 * time.Second
)

