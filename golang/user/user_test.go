package user

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/gorilla/websocket"
)

var ws *websocket.Conn
var srv = httptest.NewServer(http.HandlerFunc(handlerToBeTested))

func TestNewUser(t *testing.T) {
	user := NewUser(ws)
	if reflect.TypeOf(user.UserId).Kind() != reflect.Int {
		t.Error("Expected type of user.UserId as Int but got %v", reflect.TypeOf(user.UserId).Kind())
	}

}

func handlerToBeTested(w http.ResponseWriter, req *http.Request) {
	conn, err := upgrader.Upgrade(w, req, nil)
	ws = conn
	if err != nil {
		fmt.Println("Error in creating the websocket connection")
	}
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}
