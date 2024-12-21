package game

import (
	"errors"
	"fmt"
)

// Game represents the Battleship game session
type Game struct {
	PlayerBoards map[string]*Board
	Players      []string
	CurrentTurn  int
}

// NewGame initializes a new game session
func NewGame() *Game {
	return &Game{
		PlayerBoards: make(map[string]*Board),
	}
}

// AddPlayer initializes a board for a new player
func (g *Game) AddPlayer(playerID string) error {
	if _, exists := g.PlayerBoards[playerID]; exists {
		return errors.New("player already exists")
	}

	g.PlayerBoards[playerID] = NewBoard()
	return nil
}

func (g *Game) GetCurrentPlayer() string {
	if len(g.Players) == 0 {
		return ""
	}
	return g.Players[g.CurrentTurn%len(g.Players)]
}

// NextTurn advances the turn to the next player
func (g *Game) NextTurn() {
	g.CurrentTurn = (g.CurrentTurn + 1) % len(g.Players)
}

// GetBoard returns the player's board
func (g *Game) GetBoard(playerID string) (*Board, error) {
	board, exists := g.PlayerBoards[playerID]
	if !exists {
		return nil, errors.New("player does not exist")
	}
	return board, nil
}

// AttackPlayer handles an attack on a specific player's board
func (g *Game) AttackPlayer(attackerID, defenderID string, target Coordinate) (AttackResult, error) {
	currentPlayer := g.GetCurrentPlayer()
	if attackerID != currentPlayer {
		return Invalid, errors.New("it's not your turn")
	}

	defenderBoard, exists := g.PlayerBoards[defenderID]
	if !exists {
		return Invalid, errors.New("defender does not exist")
	}

	// Process the attack
	result, err := defenderBoard.Attack(target)
	if err != nil {
		return Invalid, err
	}

	// Check if the defender has lost
	if result == Sunk && len(defenderBoard.Ships) == 0 {
		return Sunk, fmt.Errorf("game over: %s has no ships left", defenderID)
	}

	g.NextTurn()

	return result, nil
}
