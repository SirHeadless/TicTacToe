package game

import (
	"Sirheadless/TicTacToe/golang/user"
	"fmt"
	"github.com/gorilla/websocket"
	"net/http"
	"net/http/httptest"

	"strconv"
	"testing"
)

var game *Game

var ttest *testing.T

func TestNewGame(t *testing.T) {

	// Initiate game

	userX := user.NewUser(ws)

	X := "x"
	EmptyField := ""
	game = NewGame(userX)
	if game.Turn() != "x" || game.UserO() != nil || game.UserX() == nil {
		t.Errorf("Expected [game.Turn, game.UserY, game.UserY] is [x, nil, <> nill] but got [%v,%v,%v]", game.Turn(), game.UserO(), game.UserX())
	}
	expField := startField
	if game.Field() != expField {
		t.Errorf("Expected empty field but got %v", startField)
	}

	// Test before game Started

	t.Run("SetFieldBeforeStarted", func(t *testing.T) { testUpdateField(1, "x", false, startField, t) })

	t.Run("CheckIfGameStartedBefireStart", func(t *testing.T) { testGameStarted(false, t) })

	t.Run("CheckIfGameIsWoneBeforeStart", func(t *testing.T) { testIsWon(false, nil, nil, t) })

	//Test while game is running

	userO := user.NewUser(ws2)
	game.SetUserO(userO)

	t.Run("CheckIfGameStarted 2", func(t *testing.T) { testGameStarted(true, t) })

	t.Run("UpdateField1x", func(t *testing.T) { testUpdateField(1, "x", true, [9]string{"", "x", "", "", "", "", "", "", ""}, t) })

	t.Run("CheckIfWon 2", func(t *testing.T) { testIsWon(false, nil, nil, t) })

	t.Run("UpdateField3x", func(t *testing.T) { testUpdateField(3, "x", false, [9]string{"", "x", "", "", "", "", "", "", ""}, t) })

	t.Run("UpdateField3o", func(t *testing.T) { testUpdateField(3, "o", true, [9]string{"", "x", "", "o", "", "", "", "", ""}, t) })

	t.Run("CheckIfWon 3", func(t *testing.T) { testIsWon(false, nil, nil, t) })

	t.Run("UpdateField4x", func(t *testing.T) { testUpdateField(4, "x", true, [9]string{"", "x", "", "o", "x", "", "", "", ""}, t) })

	t.Run("CheckIfWon 4", func(t *testing.T) { testIsWon(false, nil, nil, t) })

	//Field is already taken
	t.Run("UpdateTakenField3o", func(t *testing.T) {
		testUpdateField(3, "o", false, [9]string{"", "x", "", "o", "x", "", "", "", ""}, t)
	})

	t.Run("CheckIfWon 3", func(t *testing.T) {
		testIsWon(false, nil, nil, t)
	})

	t.Run("UpdateField5o", func(t *testing.T) {
		testUpdateField(5, "o", true, [9]string{"", "x", "", "o", "x", "o", "", "", ""}, t)
	})

	t.Run("CheckIfWon 4", func(t *testing.T) {
		testIsWon(false, nil, nil, t)
	})

	t.Run("UpdateField7x", func(t *testing.T) {
		testUpdateField(7, "x", true, [9]string{"", "x", "", "o", "x", "o", "", "x", ""}, t)
	})

	t.Run("CheckIfWon 5 with Win", func(t *testing.T) {
		testIsWon(true, &[3]int{1, 4, 7}, &X, t)
	})

	t.Run("StartNewGame", func(t *testing.T) {
		testStartNewGame(startField, false, "x", [3]int{0, 0, 0}, "", t)
	})

	t.Run("UpdateField0x-2", func(t *testing.T) {
		testUpdateField(0, "x", true, [9]string{"x", "", "", "", "", "", "", "", ""}, t)
	})

	t.Run("UpdateField1o2", func(t *testing.T) {
		testUpdateField(1, "o", true, [9]string{"x", "o", "", "", "", "", "", "", ""}, t)
	})

	t.Run("UpdateField3x2", func(t *testing.T) {
		testUpdateField(3, "x", true, [9]string{"x", "o", "", "x", "", "", "", "", ""}, t)
	})

	t.Run("UpdateField4o2", func(t *testing.T) {
		testUpdateField(4, "o", true, [9]string{"x", "o", "", "x", "o", "", "", "", ""}, t)
	})

	t.Run("UpdateField7x2", func(t *testing.T) {
		testUpdateField(7, "x", true, [9]string{"x", "o", "", "x", "o", "", "", "x", ""}, t)
	})

	t.Run("UpdateField6o2", func(t *testing.T) {
		testUpdateField(6, "o", true, [9]string{"x", "o", "", "x", "o", "", "o", "x", ""}, t)
	})

	t.Run("UpdateField2x2", func(t *testing.T) {
		testUpdateField(2, "x", true, [9]string{"x", "o", "x", "x", "o", "", "o", "x", ""}, t)
	})

	t.Run("UpdateField5o2", func(t *testing.T) {
		testUpdateField(5, "o", true, [9]string{"x", "o", "x", "x", "o", "o", "o", "x", ""}, t)
	})

	t.Run("UpdateField8x2", func(t *testing.T) {
		testUpdateField(8, "x", true, [9]string{"x", "o", "x", "x", "o", "o", "o", "x", "x"}, t)
	})

	t.Run("TestIfGameIsOver", func(t *testing.T) {
		testIsGameOver(true, &zeroRow, &EmptyField, t)
	})

	t.Run("TestStartNewGameAfterGameIsFull", func(t *testing.T) {
		testStartNewGame(startField, false, X, zeroRow, EmptyField, t)
	})

}

func testStartNewGame(expField [9]string, expFinished bool, expTurn string, expWinRow [3]int, expWinner string, t *testing.T) {
	game.StartNewGame()
	if game.field != expField {
		t.Errorf("expected[%v] result [%v]", expField, game.field)
	}
	if game.finished != expFinished {
		t.Errorf("expected[%v] result [%v]", expFinished, game.finished)
	}
	if game.turn != expTurn {
		t.Errorf("expected[%v] result [%v]", expTurn, game.turn)
	}
	if game.winRow != expWinRow {
		t.Errorf("expected[%v] result [%v]", expWinRow, game.winRow)
	}
	if game.winner != expWinner {
		t.Errorf("expected[%v] result [%v]", expWinner, game.winner)
	}

}

func testUpdateField(field int, symbol string, expIsUpdated bool, expField [9]string, t *testing.T) {
	isUpdated := game.SetOneField(field, symbol)
	if isUpdated != expIsUpdated {
		t.Errorf("expected[%v] result [%v]", expIsUpdated, isUpdated)
	}

	if game.Field() != expField {
		t.Errorf("game.Field: expected[%v] result [%v]", expField, game.Field())
	}
}

func testGameStarted(expStarted bool, t *testing.T) {
	started := game.GameStarted()
	if started != expStarted {
		t.Errorf("game.GameStarted(): expected[%v] result [%v]", strconv.FormatBool(started), strconv.FormatBool(expStarted))
	}
}

func testIsGameOver(expIsOver bool, expRow *[3]int, expSymbol *string, t *testing.T) {
	isOver := game.IsGameOver()
	if isOver != expIsOver {
		t.Errorf("expected[%v] result [%v]", expIsOver, isOver)
	}
	if isOver {
		if game.WinRow() != *expRow {
			t.Errorf("expected[%v] result [%v]", *expRow, game.WinRow())
		}
		if game.Winner() != *expSymbol {
			t.Errorf("expected[%v] result [%v]", *expSymbol, game.Winner())
		}
		expIsFinished := true
		if game.Finished() != expIsFinished {
			t.Errorf("expected[%v] result [%v]", expIsFinished, game.Finished())
		}
	}
}

func testIsWon(expIsWon bool, expRow *[3]int, expSymbol *string, t *testing.T) {
	isWon := game.CheckIfWon()
	if isWon != expIsWon {
		t.Errorf("expected[%v] result [%v]", expIsWon, isWon)
	}
	if isWon {
		if game.WinRow() != *expRow {
			t.Errorf("expected[%v] result [%v]", *expRow, game.WinRow())
		}
		if game.Winner() != *expSymbol {
			t.Errorf("expected[%v] result [%v]", *expSymbol, game.Winner())
		}
		expIsFinished := true
		if game.Finished() != expIsFinished {
			t.Errorf("expected[%v] result [%v]", expIsFinished, game.Finished())
		}
	}
}

var ws *websocket.Conn
var ws2 *websocket.Conn
var srv = httptest.NewServer(http.HandlerFunc(handlerToBeTested))

func handlerToBeTested(w http.ResponseWriter, req *http.Request) {
	conn, err := upgrader.Upgrade(w, req, nil)
	conn2, err := upgrader.Upgrade(w, req, nil)
	ws = conn
	ws2 = conn2
	if err != nil {
		fmt.Println("Error in creating the websocket connection")
	}
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}
