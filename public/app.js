// new Vue({
//     el: '#app',

//     data: {
//         ws: null, // Our websocket
//         newMsg: '', // Holds new messages to be sent to the server
//         chatContent: '', // A running list of chat messages displayed on the screen
//         email: null, // Email address used for grabbing an avatar
//         username: null, // Our username
//         joined: false // True if email and username have been filled in
//     },

//     created: function () {
//         var self = this;
//         this.ws = new WebSocket('ws://' + window.location.host + '/ws');
//         this.ws.addEventListener('message', function (e) {
//             var msg = JSON.parse(e.data);
//             self.chatContent += '<div class="chip">'
//                 + '<img src="' + self.gravatarURL(msg.email) + '">' // Avatar
//                 + msg.username
//                 + '</div>'
//                 + emojione.toImage(msg.message) + '<br/>'; // Parse emojis

//             var element = document.getElementById('chat-messages');
//             element.scrollTop = element.scrollHeight; // Auto scroll to the bottom
//         });
//     },

//     methods: {
//         send: function () {
//             if (this.newMsg != '') {
//                 this.ws.send(
//                     JSON.stringify({
//                         email: this.email,
//                         username: this.username,
//                         message: $('<p>').html(this.newMsg).text() // Strip out html
//                     }
//                     ));
//                 this.newMsg = ''; // Reset newMsg
//             }
//         },

//         join: function () {
//             if (!this.email) {
//                 Materialize.toast('You must enter an email', 2000);
//                 return
//             }
//             if (!this.username) {
//                 Materialize.toast('You must choose a username', 2000);
//                 return
//             }
//             this.email = $('<p>').html(this.email).text();
//             this.username = $('<p>').html(this.username).text();
//             this.joined = true;
//         },

//         gravatarURL: function (email) {
//             return 'http://www.gravatar.com/avatar/' + CryptoJS.MD5(email);
//         }
//     }
// });

document.addEventListener("DOMContentLoaded", function () {

    var ws = new WebSocket('ws://' + window.location.host + '/ws');


    function TicTacToe(element) {
        var current = 0,
            players = ["x", "o"],
            player,
            field = document.createElement("table"),
            caption = document.createElement("caption"),
            labels = [
                ["oben links", "oben mittig", "oben rechts"],
                ["Mitte links", "Mitte mittig", "Mitte rechts"],
                ["unten links", "unten mittig", "unten rechts"]
            ],
            messages = {
                "o's-turn": "Spieler O ist am Zug.",
                "x's-turn": "Spieler X ist am Zug.",
                "o-wins": "Spieler O gewinnt.",
                "x-wins": "Spieler X gewinnt.",
                "you-win": "You Win",
                "you-lose": "You Lose",
                "draw": "Das Spiel endet unentschieden.",
                "instructionNoOppon": "Waiting for opponent!",
                "instructionOppon": "Opponent was found! You are player ",
                "select": "wählen",
                "new game?": "Neues Spiel?"
            },
            finished, games, b, c, i, r, tr;

        ws.addEventListener('message', function (e) {
            var msg = JSON.parse(e.data);
            var instrElem = document.getElementById("instruction");
            if (msg.type === "field") {
                var element = document.getElementById(msg.message);
                element.className = msg.player; // Klassennamen vergeben
                element.innerHTML = msg.player;
                if (msg.player === "x") {
                    caption.innerHTML = messages["o's-turn"];
                } else {
                    caption.innerHTML = messages["x's-turn"];
                }
            } else if (msg.type === "oponentConnected") {
                player = msg.message
                instrElem.innerHTML = messages["instructionOppon"] + player
            } else if (msg.type === "won") {
                if (player === msg.message) {
                    instrElem.innerHTML = messages["you-win"];
                } else {
                    instrElem.innerHTML = messages["you-lose"];
                }
            } else if (msg.type === "winRow") {
                finished = true;
                var winRow = JSON.parse(msg.message);
                var tds = field.getElementsByTagName("td");
                highlightCells([
                    tds[winRow.f1], tds[winRow.f2], tds[winRow.f3]
                ]);

                showNewGameButton();

            } else if (msg.type == "full") {
                finished = true;
                field.className = "game-over";
                instrElem.innerHTML = messages["draw"];
                caption.innerHTML = messages["draw"];

                showNewGameButton();

            } else if (msg.type == "newGameStarted") {
                instrElem.innerHTML = messages["instructionOppon"] + player
                startNewGame();
            }
        });

        function showNewGameButton() {
            var buttons = field.getElementsByTagName("button");
            while (buttons.length) {
                buttons[0].parentNode.removeChild(buttons[0]);
            }

            // new game?
            buttons = document.createElement("button");
            buttons.innerHTML = messages["new game?"];

            caption.appendChild(document.createTextNode(" "));
            caption.appendChild(buttons);

            buttons.addEventListener("click", function (event) {
                sendMessage("", "newGame");
            });
        }

        function startNewGame() {
            var buttons = field.getElementsByTagName("button");
            while (buttons.length) {
                buttons[0].parentNode.removeChild(buttons[0]);
            } 
            caption.innerHTML = messages["x's-turn"];
            var cells = field.getElementsByTagName("td"),
                button, cell;

            // reset game
            current = 0;
            finished = false;
            field.removeAttribute("class");

            for (r = 0; r < 3; r++) {
                for (c = 0; c < 3; c++) {
                    // reset cell
                    cell = cells[r * 3 + c];
                    cell.removeAttribute("class");
                    cell.innerHTML = "";

                    // re-insert button
                    button = document.createElement("button");
                    button.innerHTML = labels[r][c] + " " + messages["select"];

                    cell.appendChild(button);
                }
            }
        }

        function sendMessage(message, type) {
            ws.send(
                JSON.stringify({
                    message: message,
                    type: type
                }
                ));
            newMsg = ''; // Reset newMsg
        }

        function highlightCells(cells) {
            cells.forEach(function (node) {
                var el = document.createElement("strong");

                el.innerHTML = node.innerHTML;

                node.innerHTML = "";
                node.appendChild(el);
                node.classList.add("highlighted");
            });
        }

        // click / tap verarbeiten
        function sendRequest(event) {
            // Tabellenzelle bestimmen
            var td = event.target;

            // Button oder Zelle?
            while (td.tagName.toLowerCase() != "td"
                && td != field
            ) {
                td = td.parentNode;
            }

            // Zelle bei Bedarf markieren
            if (!finished
                && td.tagName.toLowerCase() == "td"
                && td.className.length < 1
            ) {

                sendMessage(td.getAttribute("id"), "field");

            }
        }

        // Spielanleitung ins Dokument einfügen
        b = document.createElement("p");
        b.setAttribute("id", "instruction");
        b.innerHTML = messages["instructionNoOppon"];
        element.appendChild(b);

        // Tabelle ins Dokument einfügen
        element.appendChild(field);

        // Tabelle aufbauen
        field.appendChild(caption); // Beschriftung
        field.appendChild(document.createElement("tbody"));

        // Hinweis einrichten
        caption.innerHTML = messages[
            players[current] + "'s-turn"
        ];

        tdId = 0
        for (r = 0; r < 3; r++) {
            // neue Tabellenzeile
            tr = document.createElement("tr");

            field.lastChild.appendChild(tr);

            for (c = 0; c < 3; c++) {
                // neue Tabellenzelle
                td = document.createElement("td");
                td.setAttribute("id", "field" + tdId);
                tr.appendChild(td);

                // Klickbutton
                b = document.createElement("button");
                b.innerHTML = labels[r][c] + " " + messages["select"];

                tr.lastChild.appendChild(b);
                tdId++;
            }
        }

        // Ereignis bei Tabelle überwachen
        field.addEventListener("click", sendRequest);
    }

    // finde alle Spiel-Platzhalter
    games = document.querySelectorAll(".tic-tac-toe");

    for (i = 0; i < games.length; i++) {
        TicTacToe(games[i]); // aktuelles Fundstück steht in games[i]
    }
});