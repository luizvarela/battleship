package game

import "errors"

const BoardSize = 10

type AttackResult string

const (
	Hit     AttackResult = "HIT"
	Miss    AttackResult = "MISS"
	Sunk    AttackResult = "SUNK"
	Invalid AttackResult = "INVALID"
)

type Orientation string

const (
	Horizontal Orientation = "H"
	Vertical   Orientation = "V"
)

// Ship represents a ship with its size and placement status
type Ship struct {
	Name        string
	Length      int
	Coordinates []Coordinate
}

// Coordinate represents a position on the board
type Coordinate struct {
	X int
	Y int
}

// Board represents a player's game board
type Board struct {
	Grid      [BoardSize][BoardSize]string
	Ships     []Ship
	ShipCount int
}

// NewBoard initializes a new empty game board
func NewBoard() *Board {
	return &Board{
		Grid:      [BoardSize][BoardSize]string{},
		Ships:     make([]Ship, 0),
		ShipCount: 0,
	}
}

// PlaceShip attempts to place a ship on the board
func (b *Board) PlaceShip(ship Ship, start Coordinate, orientation Orientation) error {
	if ship.Length <= 0 {
		return errors.New("invalid ship length")
	}

	// Check if the ship can be placed
	coordinates := make([]Coordinate, ship.Length)
	for i := 0; i < ship.Length; i++ {
		var x, y int
		if orientation == Horizontal {
			x, y = start.X, start.Y+i
		} else if orientation == Vertical {
			x, y = start.X+i, start.Y
		}

		if x < 0 || y < 0 || x >= BoardSize || y >= BoardSize {
			return errors.New("ship placement out of bounds")
		}

		if b.Grid[x][y] != "" {
			return errors.New("ship placement overlaps with another ship")
		}

		coordinates[i] = Coordinate{X: x, Y: y}
	}

	// Place the ship
	for _, coord := range coordinates {
		b.Grid[coord.X][coord.Y] = "S" // "S" marks a ship part
	}

	ship.Coordinates = coordinates
	b.Ships = append(b.Ships, ship)
	b.ShipCount++

	return nil
}

func (b *Board) Attack(target Coordinate) (AttackResult, error) {
	if target.X < 0 || target.Y < 0 || target.X >= BoardSize || target.Y >= BoardSize {
		return Invalid, errors.New("attack out of bounds")
	}

	cell := b.Grid[target.X][target.Y]
	switch cell {
	case "H", "M": // Already attacked
		return Invalid, errors.New("cell already attacked")
	case "S": // Hit a ship
		b.Grid[target.X][target.Y] = "H" // Mark as hit
		// Check if the ship is sunk
		for i, ship := range b.Ships {
			for _, coord := range ship.Coordinates {
				if coord == target {
					// Remove the hit coordinate from the ship
					b.Ships[i].Coordinates = removeCoordinate(ship.Coordinates, coord)
					if len(b.Ships[i].Coordinates) == 0 {
						return Sunk, nil
					}
					return Hit, nil
				}
			}
		}
	case "": // Miss
		b.Grid[target.X][target.Y] = "M" // Mark as miss
		return Miss, nil
	}

	return Invalid, nil
}

// AllShipsSunk returns true if all ships on the board have been sunk
func (b *Board) AllShipsSunk() bool {
	// Check each ship
	for _, ship := range b.Ships {
		// If any ship has remaining coordinates (not fully hit), return false
		if len(ship.Coordinates) > 0 {
			return false
		}
	}
	// All ships have been sunk
	return true
}

// Helper function to remove a coordinate from a slice
func removeCoordinate(coords []Coordinate, target Coordinate) []Coordinate {
	for i, coord := range coords {
		if coord == target {
			return append(coords[:i], coords[i+1:]...)
		}
	}
	return coords
}
