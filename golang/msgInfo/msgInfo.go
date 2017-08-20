package msgInfo

import (
	"Sirheadless/TicTacToe/golang/game"
	"Sirheadless/TicTacToe/golang/msgInfo/message"
	"Sirheadless/TicTacToe/golang/user"
)

//MessageInformation contains a message and additional necessary information like receiver
type MessageInformation struct {
	message *message.Message
	game    *game.Game
	user    *user.User
}

//NewMessageInformation creates a new MessageInformation
func NewMessageInformation(msgType string, msg string, game *game.Game, user *user.User) *MessageInformation {
	msgInfo := new(MessageInformation)
	msgInfo.message = message.NewMessage(msgType, msg)
	msgInfo.game = game
	msgInfo.user = user
	return msgInfo
}

//NewMessageInformationByMessage creates a new MessageInformation
func NewMessageInformationByMessage(msg *message.Message, game *game.Game, user *user.User) *MessageInformation {
	msgInfo := new(MessageInformation)
	msgInfo.message = msg
	msgInfo.game = game
	msgInfo.user = user
	return msgInfo
}

//Message return message
func (msgInfo MessageInformation) Message() *message.Message {
	return msgInfo.message
}

//Game return game
func (msgInfo MessageInformation) Game() *game.Game {
	return msgInfo.game
}

//User return user
func (msgInfo MessageInformation) User() *user.User {
	return msgInfo.user
}
