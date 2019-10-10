package board

import (
	"fmt"
	"strconv"
)

//MaxScore is the max possible score for a board
const MaxScore = 99999999

//MinScore is the min possible score for a board
const MinScore = -99999999

type ray struct {
	x, y         int
	dx, dy       int
	destx, desty int
	steps        int
}

func (r *ray) step() bool {
	r.steps++
	nx := r.x + (r.steps * r.dx)
	ny := r.y + (r.steps * r.dy)
	if nx < 0 || nx > 7 || ny < 0 || ny > 7 {
		r.steps--
		return false
	}
	r.destx = nx
	r.desty = ny
	return true
}

func (r *ray) squares(stone int) []Square {
	ret := []Square{}
	for d := 1; d < r.steps; d++ {
		ret = append(ret, Square{
			x:     r.x + d*r.dx,
			y:     r.y + d*r.dy,
			Stone: stone,
		})
	}
	return ret
}

//Square is a square on the reversi board
type Square struct {
	x, y, Stone int
}

func (s *Square) adj() []ray {
	rays := []ray{}
	for x := -1; x <= 1; x++ {
		for y := -1; y <= 1; y++ {
			if x == 0 && y == 0 {
				continue
			}
			if s.x+x < 0 || s.x+x > 7 || s.y+y < 0 || s.y+y > 7 {
				continue
			}

			rays = append(rays, ray{
				x:     s.x,
				y:     s.y,
				steps: 1,
				dx:    x,
				dy:    y,
				destx: x + s.x,
				desty: y + s.y,
			})

		}
	}
	return rays
}

//Move reperesents a move on the reversi board
type Move struct {
	squares [64]Square
	index   int
}

func newMove(x, y, player int) Move {
	ret := Move{
		squares: [64]Square{
			Square{x: x, y: y, Stone: player},
		},
		index: 1,
	}
	return ret
}

func (m Move) String() string {
	return fmt.Sprintf("(x: %d, y: %d)", m.squares[0].x, m.squares[0].y)
}

//XY gives the x, y coodinates to play for the move
func (m Move) XY() (int, int) {
	return m.squares[0].x, m.squares[0].y
}

func (m *Move) add(r ray, player int) {
	for d := 1; d < r.steps; d++ {
		s := Square{
			x:     r.x + d*r.dx,
			y:     r.y + d*r.dy,
			Stone: player,
		}
		m.squares[m.index] = s
		m.index++
	}
}

//Board is a representation of a reversi board
type Board [64]int

//Move creates a new board from the given board and applies the supplied move
func (b Board) Move(m Move) Board {
	nb := b

	for i := 0; i < m.index; i++ {
		m := m.squares[i]
		nb[m.x*8+m.y] = m.Stone
	}

	return nb
}

func (b Board) String() string {
	ret := "\n"
	for x := 7; x >= 0; x-- {
		for y := 0; y < 8; y++ {
			ret += strconv.Itoa(b[x*8+y]) + " "
		}
		ret += "\n"
	}

	return ret
}

//Square returns the value at the x, y square
func (b Board) Square(x, y int) Square {
	return Square{
		x:     x,
		y:     y,
		Stone: b[x*8+y],
	}
}

//ValidMoves returns all the valid moves for the given player on the board
func (b Board) ValidMoves(player int) []Move {
	enemy := 2
	if player == 2 {
		enemy = 1
	}

	moves := []Move{}
	if b.Square(3, 3).Stone == 0 {
		return []Move{newMove(3, 3, player)}
	}
	if b.Square(4, 3).Stone == 0 {
		return []Move{newMove(4, 3, player)}
	}
	if b.Square(4, 4).Stone == 0 {
		return []Move{newMove(4, 4, player)}
	}
	if b.Square(3, 4).Stone == 0 {
		return []Move{newMove(3, 4, player)}
	}

	for x := 0; x < 8; x++ {
		for y := 0; y < 8; y++ {
			squ := b.Square(x, y)
			if squ.Stone != 0 {
				//Move from empty squares only
				continue
			}
			m := newMove(x, y, player)
			valid := false

			//get adjcent rays
			rays := squ.adj()
			for _, r := range rays {
				if b.Square(r.destx, r.desty).Stone != enemy {
					continue
				}
				for r.step() {
					piece := b.Square(r.destx, r.desty).Stone
					if piece == player {
						//this is a valid move
						valid = true
						m.add(r, player)
						break
					}
					if piece == 0 {
						//this is not a valid move
						break
					}
					//If it's an enemy piece keep looking
				}
			}

			if valid {
				moves = append(moves, m)
			}
		}
	}

	return moves
}
