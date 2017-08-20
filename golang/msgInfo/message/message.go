package message

import (
// "github.com/golang/glog"
)

//Message is a type that will be send as JSON to the client
type Message struct {
	Message string `json:"message"`
	MsgType string `json:"type"`
	Player  string `json:"player"`
}

//NewMessage creates a new Message object
func NewMessage(msgType string, message string) *Message {
	msg := new(Message)
	msg.MsgType = msgType
	msg.Message = message
	return msg
}
