package user

import (
	"math/rand"

	"github.com/gorilla/websocket"
)

type User struct {
	userId int
	name   string
	ws     *websocket.Conn
}

func NewUser(ws *websocket.Conn) *User {
	user := new(User)
	user.userId = rand.Int()
	user.ws = ws
	return user
}

//UserId return the current userId
func (user User) UserId() int {
	return user.userId
}

//Name return the current name
func (user User) Name() string {
	return user.name
}

//Ws return the current websocket connection
func (user User) Ws() *websocket.Conn {
	return user.ws
}
