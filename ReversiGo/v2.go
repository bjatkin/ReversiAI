package main

/*
import (
	"fmt"
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

func (b *board) validMoves() []move {
	player := b.turn
	enemy := "2"
	if player == "2" {
		enemy = "1"
	}

	moves := []move{}
	if b.square(3, 3).stone == "0" {
		moves = append(moves, move{
			orgX:  3,
			orgY:  3,
			stone: player,
		})
	}
	if b.square(4, 3).stone == "0" {
		moves = append(moves, move{
			orgX:  4,
			orgY:  3,
			stone: player,
		})
	}
	if b.square(3, 4).stone == "0" {
		moves = append(moves, move{
			orgX:  3,
			orgY:  4,
			stone: player,
		})
	}
	if b.square(4, 4).stone == "0" {
		moves = append(moves, move{
			orgX:  4,
			orgY:  4,
			stone: player,
		})
	}
	if len(moves) > 0 {
		return moves
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
					destx, desty := a.dest()
					piece := b.square(a.dest()).stone
					if piece == "0" {
						//this is a valid move
						moves = append(moves, a)
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
	return moves
}

func (b *board) square(x, y int) square {
	return square{
		x:     x,
		y:     y,
		stone: b.layout[x*8+y],
	}
}

func newBoard(b board, m move) board {
	newB := make([]string, 64)
	copy(newB, b.layout)
	//Update the board
	for s := 0; s <= m.steps; s++ {
		x := m.orgX + (s * m.dirX)
		y := m.orgY + (s * m.dirY)
		fmt.Printf("update %d, %d, %d\n", x, y, s)
		newB[x*8+y] = m.stone
	}

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
	stone    string
	orgX     int
	orgY     int
	steps    int
	dirX     int
	dirY     int
	children []move
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

func findMove(b board) move {
	best := -1
	moves := b.validMoves()
	move := move{ //pick a random move if we can't find a valid one
		orgX:  rand.Intn(7),
		orgY:  rand.Intn(7),
		stone: b.turn,
	}
	if len(moves) > 0 {
		move = moves[0]
	}
	for _, m := range moves {
		fmt.Printf("Start with board: %s\n", newBoard(b, m))
		if scoreMove(newBoard(b, m), b.turn, 1) > best {
			move = m
		}
	}
	return move
}

func scoreMove(b board, player string, depth int) int {
	if depth == 0 {
		enemy := "1"
		if player == "1" {
			enemy = "2"
		}
		// score := 0
		// for x := 0; x < 8; x++ {
		// 	for y := 0; y < 8; y++ {
		// 		if b.square(x, y).stone == b.turn {
		// 			score++
		// 		}
		// 	}
		// }
		b.turn = player
		mlen := len(b.validMoves())
		b.turn = enemy
		nlen := len(b.validMoves())
		corner := 0
		if b.square(0, 0).stone == player {
			corner += 10
		}
		if b.square(7, 7).stone == player {
			corner += 10
		}
		if b.square(0, 7).stone == player {
			corner += 10
		}
		if b.square(7, 0).stone == player {
			corner += 10
		}
		bcorner := 0
		if b.square(0, 0).stone == enemy {
			bcorner += 8
		}
		if b.square(7, 7).stone == enemy {
			bcorner += 8
		}
		if b.square(0, 7).stone == enemy {
			bcorner += 8
		}
		if b.square(7, 0).stone == enemy {
			bcorner += 8
		}
		xs := 0
		if b.square(1, 1).stone == player {
			xs += 10
		}
		if b.square(1, 6).stone == player {
			xs += 10
		}
		if b.square(6, 6).stone == player {
			xs += 10
		}
		if b.square(6, 1).stone == player {
			xs += 10
		}
		xsp := 0
		if b.square(1, 1).stone == enemy {
			xsp += 8
		}
		if b.square(1, 6).stone == enemy {
			xsp += 8
		}
		if b.square(6, 6).stone == enemy {
			xsp += 8
		}
		if b.square(6, 1).stone == enemy {
			xsp += 8
		}
		fmt.Printf("mlen: %d, nlen:%d, corner:%d, bcorner:%d, xs: %d, xsp: %d\n", mlen, nlen, corner, bcorner, xs, xsp)
		return mlen - nlen + corner - bcorner + xsp - xs
	}
	depth--

	var ret int
	ret = 1000000
	max := b.turn == player
	if max {
		ret = -1000000
	}
	moves := b.validMoves()
	if len(moves) == 0 {
		score := scoreMove(b, player, 0)
		return score
	}
	for _, m := range moves {
		fmt.Printf("Looking at board: %s\n", newBoard(b, m))
		score := scoreMove(newBoard(b, m), player, depth)
		//Max
		if max {
			if score > ret {
				ret = score
			}
		} else {
			if score < ret {
				ret = score
			}
		}
	}
	return ret
}
*/
