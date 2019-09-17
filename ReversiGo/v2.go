package main

import (
	"math/rand"
	"strings"
)

type board struct {
	layout []string
	turn   string
	round  string
}

func (b board) String() string {
	ret := "\n"
	for x := 0; x < 8; x++ {
		ret += strings.Join(b.layout[(x*8):(x*8)+8], "") + "\n"
	}
	return ret
}

func (b *board) validMoves() []square {
	player := b.turn
	enemy := "2"
	if player == "2" {
		enemy = "1"
	}

	squares := []square{}
	if b.square(3, 3).stone == "0" {
		squares = append(squares, square{
			stone: player,
			x:     3,
			y:     3,
		})
	}
	if b.square(4, 3).stone == "0" {
		squares = append(squares, square{
			stone: player,
			x:     4,
			y:     3,
		})
	}
	if b.square(3, 4).stone == "0" {
		squares = append(squares, square{
			stone: player,
			x:     3,
			y:     4,
		})
	}
	if b.square(4, 4).stone == "0" {
		squares = append(squares, square{
			stone: player,
			x:     4,
			y:     4,
		})
	}
	if len(squares) > 0 {
		return squares
	}

	for x := 0; x < 8; x++ {
		for y := 0; y < 8; y++ {
			s := b.square(x, y)
			if s.stone != player {
				continue
			}
			//get adjcent moves
			adj := s.adjacent()
			for _, a := range adj {
				if b.square(a.dest()).stone != enemy {
					continue
				}
				for a.next() {
					piece := b.square(a.dest()).stone
					if piece == "0" {
						//this is a valid move
						destx, desty := a.dest()
						squares = append(squares, square{
							stone: player,
							x:     destx,
							y:     desty,
						})
						break
					}
					if piece == player {
						//this is not a valid move
						break
					}
				}
			}
		}
	}
	return squares
}

func (b *board) square(x, y int) square {
	return square{
		x:     x,
		y:     y,
		stone: b.layout[x*8+y],
	}
}

func newBoard(b board, s square) board {
	newB := make([]string, 64)
	copy(newB, b.layout)
	newB[s.x*8+s.y] = s.stone
	newTurn := "1"
	if b.turn == "1" {
		newTurn = "2"
	}
	return board{
		layout: newB,
		turn:   newTurn,
	}
}

type square struct {
	x     int
	y     int
	stone string
}

func validSquare(x, y int) bool {
	// return x*8+y >= 0 && x*8+y < 64 && x >= 0 && y >= 0 && x < 8 && y < 8
	return x >= 0 && y >= 0 && x < 8 && y < 8
}

func (s *square) adjacent() []move {
	moves := []move{}
	for x := -1; x <= 1; x++ {
		for y := -1; y <= 1; y++ {
			if x == 0 && y == 0 {
				continue
			}
			if validSquare(s.x+x, s.y+y) {
				moves = append(moves, move{
					stone: s.stone,
					orgX:  s.x,
					orgY:  s.y,
					steps: 1,
					dirX:  x,
					dirY:  y,
				})
			}
		}
	}
	return moves
}

type move struct {
	stone string
	orgX  int
	orgY  int
	steps int
	dirX  int
	dirY  int
}

func (m *move) next() bool {
	m.steps++
	x, y := m.dest()
	if !validSquare(x, y) {
		m.steps--
		return false
	}
	return true
}

func (m *move) square() square {
	x, y := m.dest()
	return square{
		x:     x,
		y:     y,
		stone: m.stone,
	}
}

func (m *move) dest() (int, int) {
	return m.orgX + (m.steps * m.dirX),
		m.orgY + (m.steps * m.dirY)
}

func findMove(b board) square {
	best := -1
	moves := b.validMoves()
	move := square{
		x:     rand.Intn(7),
		y:     rand.Intn(7),
		stone: b.turn,
	}
	if len(moves) > 0 {
		move = moves[0]
	}
	for _, m := range moves {
		if scoreMove(newBoard(b, m), 3) > best {
			move = m
		}
	}
	return move
}

func scoreMove(b board, depth int) int {
	if depth == 0 {
		// e := "1"
		// if b.turn == "1" {
		// 	e = "2"
		// }
		// score := 0
		// for x := 0; x < 8; x++ {
		// 	for y := 0; y < 8; y++ {
		// 		if b.square(x, y).stone == b.turn {
		// 			score++
		// 		}
		// 	}
		// }
		len := len(b.validMoves())
		// corner := 0
		// if b.square(0, 0).stone == b.turn {
		// 	corner += 10
		// }
		// if b.square(7, 7).stone == b.turn {
		// 	corner += 10
		// }
		// if b.square(0, 7).stone == b.turn {
		// 	corner += 10
		// }
		// if b.square(7, 0).stone == b.turn {
		// 	corner += 10
		// }
		// xs := 0
		// if b.square(1, 1).stone == b.turn {
		// 	xs += 8
		// }
		// if b.square(1, 6).stone == b.turn {
		// 	xs += 8
		// }
		// if b.square(6, 6).stone == b.turn {
		// 	xs += 8
		// }
		// if b.square(6, 1).stone == b.turn {
		// 	xs += 8
		// }
		// xsp := 0
		// if b.square(1, 1).stone == e {
		// 	xsp += 8
		// }
		// if b.square(1, 6).stone == e {
		// 	xsp += 8
		// }
		// if b.square(6, 6).stone == e {
		// 	xsp += 8
		// }
		// if b.square(6, 1).stone == e {
		// 	xsp += 8
		// }
		return len // + corner - xs + xsp
	}
	depth--

	var ret int
	ret = -1000000
	moves := b.validMoves()
	if len(moves) == 0 {
		score := scoreMove(b, 0)
		return score
	}
	for _, m := range moves {
		score := scoreMove(newBoard(b, m), depth)
		if score > ret {
			ret = score
		}
	}
	return ret
}
