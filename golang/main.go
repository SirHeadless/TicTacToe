package main

import (
	"log"
	// "math/rand"
	"Sirheadless/TicTacToe/golang/game"
	"Sirheadless/TicTacToe/golang/msgInfo"
	"Sirheadless/TicTacToe/golang/msgInfo/message"
	"Sirheadless/TicTacToe/golang/user"
	"github.com/gorilla/websocket"
	"net/http"
	"strconv"
)

const symbolX = "x"
const symbolO = "o"
const opponentConnected = "oponentConnected"
const won = "won"
const full = "full"
const winRow = "winRow"
const newGameStarted = "newGameStarted"

var clients = make(map[*websocket.Conn]*bool)          // connected clients
var broadcast = make(chan *msgInfo.MessageInformation) // broadcast channel

var openGame *game.Game

// Configure the upgrader
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		log.Printf("Websocket upgrader was called")
		return true
	},
}

func main() {
	log.Printf("Start handleConnectsions Clients: %v", clients)

	// Create a simple file server
	fs := http.FileServer(http.Dir("../public"))
	http.Handle("/", fs)

	// Configure websocket route
	http.HandleFunc("/ws", handleConnections)

	// Start listening for incoming chat messages
	go handleMessages()

	// Start the server on localhost port 8000 and log any errors
	log.Println("http server started on :8000")
	err := http.ListenAndServe(":8000", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func handleConnections(w http.ResponseWriter, r *http.Request) {
	var playedGame *game.Game
	log.Printf("Start handleConnections Clients: %v", clients)

	// Upgrade initial GET request to a websocket
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}

	// Make sure we close the connection when the function returns
	defer ws.Close()

	userNew := user.NewUser(ws)
	log.Printf("Open Game: %v", openGame)
	if openGame == nil {
		playedGame = game.NewGame(userNew)
		openGame = playedGame
	} else {
		openGame.SetUserO(userNew)
		playedGame = openGame
		openGame = nil
		sendOpponentsConnected(playedGame)
	}

	log.Printf("New user: %v", userNew)
	// Register our new client
	log.Printf("Client of ws: %v", clients[ws])

	for {
		var msg message.Message
		// Read in a new message as JSON and map it to a Message object
		err := ws.ReadJSON(&msg)
		log.Printf("GOT MESSAGE: %v FROM: %v", msg, userNew.UserId())
		if err != nil {
			log.Printf("error: %v", err)
			delete(clients, ws)
			break
		}
		if msg.MsgType == "field" || !playedGame.Finished() {
			if userNew == playedGame.UserX() {
				msg.Player = symbolX
			} else {
				msg.Player = symbolO
			}
			fieldNr, err := strconv.Atoi(msg.Message[5:len(msg.Message)])
			if err != nil {
				log.Printf("error: %v", err)
			}
			if playedGame.SetOneField(fieldNr, msg.Player) {
				msgInfo := msgInfo.NewMessageInformationByMessage(&msg, playedGame, nil)
				// Send the newly received message to the broadcast channel
				broadcast <- msgInfo
				if playedGame.IsGameOver() {
					if playedGame.Winner() != "" {
						sendWonAndWinRow(playedGame)
					} else {
						sendGameFull(playedGame)
					}
				}
			}
		} else if msg.MsgType == "newGame" || playedGame.Finished() {
			log.Printf("START NEW GAME")
			playedGame.StartNewGame()
			sendNewGame(playedGame)
		}
	}
}

func sendGameFull(playedGame *game.Game) {
	msgInfoNewGame := msgInfo.NewMessageInformation(full, "", playedGame, nil)
	broadcast <- msgInfoNewGame
}

func sendNewGame(playedGame *game.Game) {
	msgInfoNewGame := msgInfo.NewMessageInformation(newGameStarted, "", playedGame, nil)
	broadcast <- msgInfoNewGame
}

func sendWonAndWinRow(playedGame *game.Game) {
	msgInfoWon := msgInfo.NewMessageInformation(won, playedGame.Winner(), playedGame, nil)
	broadcast <- msgInfoWon
	//TODO: Improve
	winningRow := playedGame.WinRow()
	messageRow := "{\"f1\":" + strconv.Itoa(winningRow[0]) + ",\"f2\":" + strconv.Itoa(winningRow[1]) + ",\"f3\":" + strconv.Itoa(winningRow[2]) + "}"
	msgInfoRow := msgInfo.NewMessageInformation("winRow", messageRow, playedGame, nil)
	broadcast <- msgInfoRow
}

func sendOpponentsConnected(playedGame *game.Game) {
	msgInfoX := msgInfo.NewMessageInformation(opponentConnected, symbolX, playedGame, playedGame.UserX())
	msgInfoO := msgInfo.NewMessageInformation(opponentConnected, symbolO, playedGame, playedGame.UserO())

	broadcast <- msgInfoX
	broadcast <- msgInfoO
}

func handleMessages() {
	log.Printf("HandleMessages Clients: %v", clients)

	for {
		// Grab the next message from the broadcast channel
		msgInfoToSend := <-broadcast
		msgToSend := msgInfoToSend.Message()
		gameToSend := msgInfoToSend.Game()
		if msgInfoToSend.User() == nil {
			if gameToSend.UserO() != nil {
				sendMessage(gameToSend.UserX().Ws(), msgToSend)
				sendMessage(gameToSend.UserO().Ws(), msgToSend)
			}
		} else {
			sendMessage(msgInfoToSend.User().Ws(), msgToSend)
		}
	}
}

func sendMessage(ws *websocket.Conn, msg *message.Message) {
	log.Printf("MESSAGE TO SEND: %v", msg)
	err := ws.WriteJSON(msg)
	if err != nil {
		log.Printf("error: %v", err)
		ws.Close()
		delete(clients, ws)
	}
}
