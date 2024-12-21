package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/luizvarela/battleship/game"
)

// Player represents a connected player
type Player struct {
	ID   string
	Conn *websocket.Conn
}

// GameServer manages WebSocket connections and game state
type GameServer struct {
	Players map[string]*Player
	Game    *game.Game
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all origins for simplicity
	},
}

// NewGameServer creates a new game server
func NewGameServer() *GameServer {
	return &GameServer{
		Players: make(map[string]*Player),
		Game:    game.NewGame(),
	}
}

type PlaceShipMessage struct {
	Type        string `json:"type"`
	Ship        string `json:"ship"`
	Orientation string `json:"orientation"`
	Row         int    `json:"row"`
	Col         int    `json:"col"`
}

type BoardUpdate struct {
	Type  string     `json:"type"`
	Board [][]string `json:"board"`
}

type AttackMessage struct {
	Type string `json:"type"`
	Row  int    `json:"row"`
	Col  int    `json:"col"`
}

// HandleConnection handles incoming WebSocket connections
func (s *GameServer) HandleConnection(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Failed to upgrade connection: %v\n", err)
		return
	}

	defer conn.Close()

	// Assign a player ID (e.g., based on a query parameter or unique ID)
	playerID := fmt.Sprintf("Player%d", len(s.Players)+1)
	player := &Player{ID: playerID, Conn: conn}
	s.Players[playerID] = player
	s.Game.AddPlayer(playerID)

	log.Printf("Player %s connected\n", playerID)

	board, _ := s.Game.GetBoard(playerID)
	if board != nil {
		// Convert board to client format
		clientBoard := make([][]string, len(board.Grid))
		for i := range board.Grid {
			clientBoard[i] = make([]string, len(board.Grid[i]))
			for j := range board.Grid[i] {
				if board.Grid[i][j] != "" {
					clientBoard[i][j] = "ship"
				} else {
					clientBoard[i][j] = "white"
				}
			}
		}

		// Send initial board state
		update := BoardUpdate{
			Type:  "update_board",
			Board: clientBoard,
		}

		if err := conn.WriteJSON(update); err != nil {
			log.Printf("Error sending initial board state: %v\n", err)
		}
	}

	// Listen for messages from the player
	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Printf("Error reading message: %v\n", err)
			delete(s.Players, playerID)
			break
		}

		var placeShipMessage PlaceShipMessage
		err = json.Unmarshal(message, &placeShipMessage)
		if err != nil {
			log.Printf("Error unmarshalling message: %v\n", err)
			continue
		}

		if placeShipMessage.Type == "place_ship" {
			orientation := game.Horizontal
			if placeShipMessage.Orientation == "vertical" {
				orientation = game.Vertical
			}

			// Create ship based on name
			var shipLength int
			switch placeShipMessage.Ship {
			case "carrier":
				shipLength = 5
			case "battleship":
				shipLength = 4
			case "cruiser":
				shipLength = 3
			case "submarine":
				shipLength = 3
			case "destroyer":
				shipLength = 2
			default:
				log.Printf("Invalid ship type: %s\n", placeShipMessage.Ship)
				continue
			}

			ship := game.Ship{Name: placeShipMessage.Ship, Length: shipLength}
			coord := game.Coordinate{X: placeShipMessage.Row, Y: placeShipMessage.Col}

			// Place ship
			board, _ := s.Game.GetBoard(playerID)
			err := board.PlaceShip(ship, coord, orientation)
			if err != nil {
				log.Printf("Error placing ship: %v\n", err)
				conn.WriteJSON(map[string]string{
					"type":    "error",
					"message": err.Error(),
				})
				continue
			}

			// Convert board to client format
			clientBoard := make([][]string, len(board.Grid))
			for i := range board.Grid {
				clientBoard[i] = make([]string, len(board.Grid[i]))
				for j := range board.Grid[i] {
					if board.Grid[i][j] != "" {
						clientBoard[i][j] = "ship"
					} else {
						clientBoard[i][j] = "white"
					}
				}
			}

			// Send board update to client
			update := BoardUpdate{
				Type:  "update_board",
				Board: clientBoard,
			}

			if err := conn.WriteJSON(update); err != nil {
				log.Printf("Error sending board update: %v\n", err)
				continue
			}

			// Send success message
			conn.WriteJSON(map[string]string{
				"type": "log",
				"text": fmt.Sprintf("Successfully placed %s at (%d, %d)", placeShipMessage.Ship, coord.X, coord.Y),
			})
		} else if placeShipMessage.Type == "attack" {
			var attackMsg AttackMessage
			err = json.Unmarshal(message, &attackMsg)
			if err != nil {
				log.Printf("Error unmarshalling attack message: %v\n", err)
				continue
			}

			// Get the opponent's ID (assuming 2 players)
			var opponentID string
			for id := range s.Players {
				if id != playerID {
					opponentID = id
					break
				}
			}

			if opponentID == "" {
				conn.WriteJSON(map[string]string{
					"type":    "error",
					"message": "No opponent found",
				})
				continue
			}

			// Get opponent's board and perform attack
			opponentBoard, _ := s.Game.GetBoard(opponentID)
			coord := game.Coordinate{X: attackMsg.Row, Y: attackMsg.Col}
			result, err := opponentBoard.Attack(coord)

			if err != nil {
				conn.WriteJSON(map[string]string{
					"type":    "error",
					"message": err.Error(),
				})
				continue
			}

			// Send attack result to attacker (showing hit/miss on enemy board)
			attackerResult := map[string]interface{}{
				"type":   "attack_result",
				"row":    attackMsg.Row,
				"col":    attackMsg.Col,
				"result": string(result),
			}
			conn.WriteJSON(attackerResult)

			// Send board update to opponent (showing hit/miss on their own board)
			if opponent, exists := s.Players[opponentID]; exists {
				opponentUpdate := BoardUpdate{
					Type:  "update_board",
					Board: convertBoardToClientFormat(opponentBoard),
				}
				opponent.Conn.WriteJSON(opponentUpdate)
			}

			// Check if game is over
			if opponentBoard.AllShipsSunk() {
				gameOver := map[string]string{
					"type":   "game_over",
					"winner": playerID,
				}
				// Broadcast game over to both players
				conn.WriteJSON(gameOver)
				if opponent, exists := s.Players[opponentID]; exists {
					opponent.Conn.WriteJSON(gameOver)
				}
			}

			// Log the attack
			log.Printf("Player %s attacked (%d, %d): %s\n", playerID, attackMsg.Row, attackMsg.Col, result)
		}

		log.Printf("Received from %s: %s\n", playerID, string(message))
	}
}

func main() {
	server := NewGameServer()
	http.HandleFunc("/ws", server.HandleConnection)

	log.Println("Starting server on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func convertBoardToClientFormat(board *game.Board) [][]string {
	clientBoard := make([][]string, game.BoardSize)
	for i := range clientBoard {
		clientBoard[i] = make([]string, game.BoardSize)
		for j := range clientBoard[i] {
			switch board.Grid[i][j] {
			case "S":
				clientBoard[i][j] = "ship"
			case "H":
				clientBoard[i][j] = "HIT"
			case "M":
				clientBoard[i][j] = "MISS"
			default:
				clientBoard[i][j] = "white"
			}
		}
	}
	return clientBoard
}
