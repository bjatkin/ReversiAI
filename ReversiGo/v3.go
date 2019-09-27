package main

import (
	"fmt"
	"sort"
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
	return fmt.Sprintf("(x: %d, y: %d)", m.squares[0].x, m.squares[0].y)
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

func (b *Board) Move(m Move) *Board {
	if b.turn == 1 {
		b.turn = 2
	} else {
		b.turn = 1
	}

	for i := 0; i < m.index; i++ {
		m := m.squares[i]
		b.layout[m.x*8+m.y] = m.stone
	}
	return b
}

func (b Board) String() string {
	ret := "\n"
	for x := 7; x >= 0; x-- {
		for y := 0; y < 8; y++ {
			ret += strconv.Itoa(b.layout[x*8+y]) + " "
		}
		ret += "\n"
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
	pscore := 0
	escore := 0
	for x := 0; x < 8; x++ {
		for y := 0; y < 8; y++ {
			s := b.Square(x, y).stone
			if s == player {
				pscore++
			}
			if s == enemy {
				escore++
			}
		}
	}
	if pscore+escore == 64 {
		if pscore > 32 {
			return 1000000 //You Win!
		}
		return -1000000 //You Loose :(
	}
	if pscore == 0 {
		return -1000000 //You Loose :(
	}
	if escore == 0 {
		return 1000000 //You Win!
	}
	b.turn = player
	plen := len(b.ValidMoves())
	b.turn = enemy
	elen := len(b.ValidMoves())
	pcorner := 0
	if b.Square(0, 0).stone == player {
		pcorner++
	}
	if b.Square(7, 7).stone == player {
		pcorner++
	}
	if b.Square(0, 7).stone == player {
		pcorner++
	}
	if b.Square(7, 0).stone == player {
		pcorner++
	}
	ecorner := 0
	if b.Square(0, 0).stone == enemy {
		ecorner++
	}
	if b.Square(7, 7).stone == enemy {
		ecorner++
	}
	if b.Square(0, 7).stone == enemy {
		ecorner++
	}
	if b.Square(7, 0).stone == enemy {
		ecorner++
	}
	pxs := 0
	if b.Square(1, 1).stone == player {
		pxs++
	}
	if b.Square(1, 6).stone == player {
		pxs++
	}
	if b.Square(6, 6).stone == player {
		pxs++
	}
	if b.Square(6, 1).stone == player {
		pxs++
	}
	exs := 0
	if b.Square(1, 1).stone == enemy {
		exs++
	}
	if b.Square(1, 6).stone == enemy {
		exs++
	}
	if b.Square(6, 6).stone == enemy {
		exs++
	}
	if b.Square(6, 1).stone == enemy {
		exs++
	}
	if pscore+escore < 16 {
		return plen - elen - 100*(pxs-pcorner) + 3*exs + 30*(pcorner-ecorner) - 2*pscore
	}
	if pscore+escore < 30 {
		return plen - elen - 100*(pxs-pcorner) + 6*exs + 50*(pcorner-ecorner) - pscore
	}
	if pscore+escore < 48 {
		return 50*(pcorner-ecorner) - 100*(pxs-pcorner) + pscore
	}

	return pscore - escore
}

func NewBoard(b *Board) Board {
	newB := b.layout
	return Board{
		layout: newB,
		turn:   b.turn,
	}
}

type ByScore []Board

func (b ByScore) Len() int           { return len(b) }
func (b ByScore) Swap(i, j int)      { b[i], b[j] = b[j], b[i] }
func (b ByScore) Less(i, j int) bool { return b[i].Value() < b[j].Value() }

func ScoreMove(b *Board, player int, currentBest int, depth int) int {
	if depth <= 0 {
		b.turn = player
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
		b.turn = player
		score := b.Value()
		return score
	}

	// forward-ish pruning for deeper searches
	var boards []Board
	if len(moves) > 4 {
		sortBoards := []Board{}
		for _, m := range moves {
			nb := NewBoard(b)
			nb.turn = player
			sortBoards = append(sortBoards, *nb.Move(m))
		}
		sort.Sort(ByScore(sortBoards))

		boards = sortBoards[:4] //Take the top 4 boards
		if !max {
			boards = sortBoards[len(sortBoards)-4:] //Take the bottom 4 boards
		}
	} else {
		for _, m := range moves {
			nb := NewBoard(b)
			nb.Move(m)
			boards = append(boards, nb)
		}
	}

	for _, nb := range boards {
		score := ScoreMove(&nb, player, ret, depth)
		if max {
			if score > ret {
				ret = score
			}
			//Alpha Beta pruning
			if score > currentBest {
				break
			}
		} else {
			if score < ret {
				ret = score
			}
			//Alpha Beta pruning
			if score < currentBest {
				break
			}
		}
	}
	return ret
}

func findMove(b *Board, depth int) Move {
	best := -100000000000
	moves := b.ValidMoves()
	move := Move{}

	if len(moves) == 0 {
		fmt.Printf("I forfit my turn\n")
		return move
	}
	if len(moves) == 1 {
		return moves[0]
	}
	for _, m := range moves {
		nb := NewBoard(b)
		nb.Move(m)
		s := ScoreMove(&nb, b.turn, -100000000000, depth)
		if s > best {
			best = s
			move = m
		}
	}
	return move
}
