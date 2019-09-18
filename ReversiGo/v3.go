package main

import (
	"fmt"
	"math/rand"
	"strconv"
)

type Ray struct {
	x, y   int
	dx, dy int
	s      int
}

func (r *Ray) Step() bool {
	r.s++
	nx := r.x + (r.s * r.dx)
	ny := r.y + (r.s * r.dy)
	if nx < 0 || nx > 7 || ny < 0 || ny > 7 {
		r.s--
		return false
	}
	return true
}

func (r *Ray) Dest() (int, int) {
	return r.x + (r.s * r.dx), r.y + (r.s * r.dy)
}

func (r *Ray) Squares(stone int) []Square {
	ret := []Square{}
	for d := 1; d < r.s; d++ {
		ret = append(ret, Square{
			x:     r.x + d*r.dx,
			y:     r.y + d*r.dy,
			stone: stone,
		})
	}
	return ret
}

type Square struct {
	x, y, stone int
}

func (s *Square) Adj() []Ray {
	rays := []Ray{}
	for x := -1; x <= 1; x++ {
		for y := -1; y <= 1; y++ {
			if x == 0 && y == 0 {
				continue
			}
			if s.x+x < 0 || s.x+x > 7 || s.y+y < 0 || s.y+y > 8 {
				continue
			}

			rays = append(rays, Ray{
				x:  s.x,
				y:  s.y,
				s:  1,
				dx: x,
				dy: y,
			})

		}
	}
	return rays
}

type Move struct {
	squares [64]Square
	index   int
}

func (m Move) String() string {
	return fmt.Sprintf("%v", m.squares[:m.index])
}

func (m *Move) add(squares ...Square) {
	for _, s := range squares {
		m.squares[m.index] = s
		m.index++
	}
}

type Board struct {
	layout [64]int
	turn   int
	round  int
}

func (b *Board) Move(m Move) {
	for _, m := range m.squares {
		b.layout[m.x*8+m.y] = m.stone
	}
}

func (b Board) String() string {
	ret := "\n"
	for x := 0; x < 8; x++ {
		for y := 0; y < 8; y++ {
			ret += strconv.Itoa(b.layout[x*8+y]) + " "
		}
		ret = "\n"
	}
	return ret
}

func (b *Board) Square(x, y int) Square {
	if x < 0 || x > 7 || y < 0 || y > 7 {
		return Square{}
	}

	return Square{
		x:     x,
		y:     y,
		stone: b.layout[x*8+y],
	}
}

func (b *Board) ValidMoves() []Move {
	player := b.turn
	enemy := 2
	if player == 2 {
		enemy = 1
	}

	moves := []Move{}
	if b.Square(3, 3).stone == 0 {
		m := Move{}
		m.add(Square{x: 3, y: 3, stone: player})
		moves = append(moves, m)
	}
	if b.Square(4, 3).stone == 0 {
		m := Move{}
		m.add(Square{x: 4, y: 3, stone: player})
		moves = append(moves, m)
	}
	if b.Square(3, 4).stone == 0 {
		m := Move{}
		m.add(Square{x: 3, y: 4, stone: player})
		moves = append(moves, m)
	}
	if b.Square(4, 4).stone == 0 {
		m := Move{}
		m.add(Square{x: 4, y: 4, stone: player})
		moves = append(moves, m)
	}
	if len(moves) > 0 {
		return moves
	}

	for x := 0; x < 8; x++ {
		for y := 0; y < 8; y++ {
			s := b.Square(x, y)
			if s.stone != 0 {
				continue
			}
			m := Move{}
			m.add(Square{x: x, y: y, stone: player})
			valid := false
			//get adjcent rays
			rays := s.Adj()
			for _, r := range rays {
				if b.Square(r.Dest()).stone != enemy {
					continue
				}
				for r.Step() {
					piece := b.Square(r.Dest()).stone
					if piece == player {
						//this is a valid move
						valid = true
						for _, s := range r.Squares(player) {
							m.add(s)
						}
						break
					}
					if piece == 0 {
						//this is not a valid move
						break
					}
				}
			}

			if valid {
				moves = append(moves, m)
			}
		}
	}

	return moves
}

func (b *Board) Value() int {
	enemy := 1
	player := b.turn
	if player == 1 {
		enemy = 2
	}
	score := 0
	piceCount := 0
	for x := 0; x < 8; x++ {
		for y := 0; y < 8; y++ {
			s := b.Square(x, y).stone
			if s == b.turn {
				score++
			}
			if s != 0 {
				piceCount++
			}
		}
	}
	if piceCount == 64 {
		if score > 32 {
			return 1000000
		}
		return -1000000
	}
	b.turn = player
	plen := len(b.ValidMoves())
	b.turn = enemy
	elen := len(b.ValidMoves())
	pcorner := 0
	if b.Square(0, 0).stone == player {
		pcorner += 10
	}
	if b.Square(7, 7).stone == player {
		pcorner += 10
	}
	if b.Square(0, 7).stone == player {
		pcorner += 10
	}
	if b.Square(7, 0).stone == player {
		pcorner += 10
	}
	ecorner := 0
	if b.Square(0, 0).stone == enemy {
		ecorner += 8
	}
	if b.Square(7, 7).stone == enemy {
		ecorner += 8
	}
	if b.Square(0, 7).stone == enemy {
		ecorner += 8
	}
	if b.Square(7, 0).stone == enemy {
		ecorner += 8
	}
	pxs := 0
	if b.Square(1, 1).stone == player {
		pxs += 2
	}
	if b.Square(1, 6).stone == player {
		pxs += 2
	}
	if b.Square(6, 6).stone == player {
		pxs += 2
	}
	if b.Square(6, 1).stone == player {
		pxs += 2
	}
	exs := 0
	if b.Square(1, 1).stone == enemy {
		exs += 8
	}
	if b.Square(1, 6).stone == enemy {
		exs += 8
	}
	if b.Square(6, 6).stone == enemy {
		exs += 8
	}
	if b.Square(6, 1).stone == enemy {
		exs += 8
	}
	// fmt.Printf("mlen: %d, nlen:%d, corner:%d, bcorner:%d, xs: %d, xsp: %d\n", mlen, nlen, corner, bcorner, xs, xsp)
	if piceCount > 10 {
		return plen - elen - pxs + exs + pcorner - ecorner
	}
	if piceCount > 32 {
		return 2*pcorner + 2*ecorner + score
	}

	return plen - elen - score

}

func NewBoard(b *Board) Board {
	newB := b.layout
	newTurn := 1
	if b.turn == 1 {
		newTurn = 2
	}
	return Board{
		layout: newB,
		turn:   newTurn,
	}
}

func ScoreMove(b *Board, player int, currentBest int, depth int) int {
	if depth == 0 {
		return b.Value()
	}
	depth--

	var ret int
	ret = 1000000
	max := b.turn == player
	if max {
		ret = -1000000
	}
	moves := b.ValidMoves()
	if len(moves) == 0 {
		// fmt.Printf(" - No boards to look at for %s\n", b)
		score := b.Value()
		return score
	}
	for _, m := range moves {
		nb := NewBoard(b)
		nb.Move(m)
		score := ScoreMove(&nb, player, ret, depth)
		// fmt.Printf(" - Looking at board: %s\n - score: %d\n", nb, score)
		if max {
			if score > ret {
				ret = score
			}
		} else {
			//Alpha Beta pruning
			if score < currentBest {
				break
			}
			if score < ret {
				ret = score
			}
		}
	}
	return ret
}

func findMove(b *Board, depth int) Move {
	best := -1000000
	moves := b.ValidMoves()
	move := Move{} //pick a random move if we can't find a valid one
	move.add(Square{
		x:     rand.Intn(7),
		y:     rand.Intn(7),
		stone: b.turn,
	})

	if len(moves) > 0 {
		move = moves[0]
	}
	for _, m := range moves {
		nb := NewBoard(b)
		nb.Move(m)
		// fmt.Printf("Start with board: %s\n", nb)
		s := ScoreMove(&nb, b.turn, -1000000, depth)
		// fmt.Printf("\n\nSCORE: %d\n\n", s)
		if s > best {
			// fmt.Printf("\n\nCHOOSE: %d\n\n", s)
			best = s
			move = m
		}
	}
	return move
}
