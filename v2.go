package main

type board struct {
	layout []string
	turn   string
}

func (b *board) validMoves() []move {
	player := b.turn

	moves := []move{}
	if b.square(3, 3).stone == "0" {
		moves = append(moves, move{
			stone: player,
			orgX:  3,
			orgY:  3,
		})
	}
	if b.square(4, 3).stone == "0" {
		moves = append(moves, move{
			stone: player,
			orgX:  4,
			orgY:  3,
		})

	}
	if b.square(3, 4).stone == "0" {
		moves = append(moves, move{
			stone: player,
			orgX:  3,
			orgY:  4,
		})
	}
	if b.square(4, 4).stone == "0" {
		moves = append(moves, move{
			stone: player,
			orgX:  4,
			orgY:  4,
		})
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
				for a.next() {
					piece := b.square(a.dest()).stone
					if piece == "0" {
						//this is a valid move
						moves = append(moves, a)
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

func newBoard(b board, s square) board {
	b.layout[s.x*8+s.y] = s.stone
	newTurn := "1"
	if b.turn == "1" {
		newTurn = "2"
	}
	return board{
		layout: b.layout,
		turn:   newTurn,
	}
}

type square struct {
	x     int
	y     int
	stone string
}

func validSquare(x, y int) bool {
	return x*8+y >= 0 && x*8+y < 64
}

func (s *square) adjacent() []move {
	moves := []move{}
	for x := -1; x <= 1; x++ {
		for y := -1; y <= 1; y++ {
			if x == 0 && y == 0 {
				continue
			}
			if validSquare(x, y) {
				moves = append(moves, move{
					stone: s.stone,
					orgX:  s.x,
					orgY:  s.y,
					steps: 0,
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

func findMove(b board) move {
	best := -1
	moves := b.validMoves()
	move := moves[0]
	for _, m := range moves {
		if scoreMove(newBoard(b, m.square())) > best {
			move = m
		}
	}
	return move
}

func scoreMove(b board) int {
	return len(b.validMoves())
}
