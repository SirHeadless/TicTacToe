package game

import (
	"Sirheadless/TicTacToe/golang/user"
	"github.com/golang/glog"
)

var zeroRow = [3]int{0, 0, 0}

var startField = [9]string{"", "", "", "", "", "", "", "", ""}

var winRows = [8][3]int{{0, 1, 2}, {3, 4, 5}, {6, 7, 8}, {0, 3, 6}, {1, 4, 7}, {2, 5, 8}, {0, 4, 8}, {4, 2, 6}}

// Game has all components that are required for a game
type Game struct {
	field    [9]string
	userX    *user.User
	userO    *user.User
	turn     string
	winRow   [3]int
	winner   string
	finished bool
}

//NewGame creates a new game and reutrns it
func NewGame(userX *user.User) *Game {
	game := new(Game)
	game.turn = "x"
	game.userX = userX
	game.field = startField
	game.finished = false
	game.winRow = zeroRow

	return game

}

//SetOneField updates the field of the game
func (game *Game) SetOneField(fieldNr int, symbol string) bool {
	if !game.GameStarted() {
		glog.Warningf("Executed UpdateField before game started fieldNr %v, symbol %v", fieldNr, symbol)
		return false
	}
	if game.turn != symbol {
		glog.Warningf("Exectued UpdateField with symbol where it is not its turn(%v) %v", game.turn, symbol)
		return false
	}
	if game.field[fieldNr] != "" {
		glog.Warningf("Exectued UpdateField with already used field(%v) %v", game.turn, symbol)
		return false
	}
	game.ToogleTurn()
	game.field[fieldNr] = symbol
	return true
}

//GameStarted returns if game started
func (game Game) GameStarted() bool {
	if game.userO == nil || game.userX == nil || game.finished {
		return false
	}
	return true
}

//Field return the current game field
func (game Game) Field() [9]string {
	return game.field
}

//UserX return the x player
func (game Game) UserX() *user.User {
	return game.userX
}

//UserO return the O player
func (game Game) UserO() *user.User {
	return game.userO
}

//Turn return which players turn it is
func (game Game) Turn() string {
	return game.turn
}

//WinRow returns the wining row
func (game Game) WinRow() [3]int {
	return game.winRow
}

//Winner returns the wining row
func (game Game) Winner() string {
	return game.winner
}

//Finished returns if the game is finished
func (game Game) Finished() bool {
	return game.finished
}

//SetUserO sets UserO
func (game *Game) SetUserO(user *user.User) {
	game.userO = user
}

//CheckIfWon returns if the game is has a Winner
func (game *Game) CheckIfWon() bool {
	for _, row := range winRows {
		isWinning, symbol := game.areFieldsSame(row)
		if isWinning {
			game.setFinishedVariables(true, row, symbol)
			return true
		}
	}
	return false

}

func (game Game) areFieldsSame(fields [3]int) (bool, string) {
	if (game.field[fields[0]] == "o" || game.field[fields[0]] == "x") &&
		game.field[fields[0]] == game.field[fields[1]] && game.field[fields[1]] == game.field[fields[2]] {
		return true, game.field[fields[0]]
	}
	return false, ""
}

//ToogleTurn is setting the next player to turn
func (game *Game) ToogleTurn() {
	if game.turn == "x" {
		game.turn = "o"
	} else {
		game.turn = "x"
	}
}

//StartNewGame starts a new game
func (game *Game) StartNewGame() {
	game.field = startField
	game.finished = false
	game.turn = "x"
	game.winRow = zeroRow
	game.winner = ""
}

func (game *Game) checkIfFull() bool {
	for i := 0; i < len(game.field); i++ {
		if game.field[i] == "" {
			return false
		}
	}
	game.setFinishedVariables(true, zeroRow, "")
	return true
}

//IsGameOver set returns if the game is over and sets the finish variables if so
func (game *Game) IsGameOver() bool {
	if game.CheckIfWon() {
		return true
	}
	if game.checkIfFull() {
		return true
	}
	return false
}

func (game *Game) setFinishedVariables(isFinished bool, winRow [3]int, symbol string) {
	game.finished = isFinished
	game.winRow = winRow
	game.winner = symbol
}
