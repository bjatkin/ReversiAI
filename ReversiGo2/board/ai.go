package board

var PositionsCalculated int

//ValueBoard calculates a boards value using minimax and the Board.value() function
func ValueBoard(b *Board, player, depth, currentBest int, stop chan bool, value chan int) int {
	select {
	case _, open := <-stop:
		if !open {
			return 0
		}
		//we've reached max depth, calculate up from here
		PositionsCalculated++
		score := b.value(player)
		if depth == 0 {
			value <- score
		}
		return score
	default:
		//Keep moving down the tree, we havent been stoped yet

		var target int
		moves := b.ValidMoves(player)
		if depth%3 == 0 {
			scores, moves := prune(b, moves, player)
		}

		newCurrentBest := currentBest
		for i, move := range moves {
			nb := b.Move(move)
			score := ValueBoard(&nb, player, depth+1, newCurrentBest, stop, value)
			target = miniMax(depth, score, target, i)
		}

		if depth == 0 {
			value <- target
		}

		if len(moves) == 0 {
			PositionsCalculated++
			return b.value(player)
		}
		return target
	}
}

func miniMax(depth, score, target, i int) int {
	if i == 0 {
		return score
	}

	if depth%2 == 0 && score < target { //Min
		return score
	}
	if score > target { //Max
		return score
	}
	return target
}

func prune(b *Board, moves []Move, player int) ([]int, []Move) {
	scores := []int{}
	for _, move := range moves {
		nb := b.Move(move)
		scores = append(scores, nb.value(player))
	}
	return scores, moves
}

type boardStats struct {
	win        bool
	loss       bool
	stoneCount int
	frontier   int
	sweet16    int
	nEdge      int
	sEdge      int
	eEdge      int
	wEdge      int
	corner     int
	mobility   int
	xSquares   int
	cSquares   int
	bSquares   int
}

func TestValue(b Board, player int) int {
	return b.value(player)
}

func (b Board) value(player int) int {
	playerStats := boardStats{}
	turn := 0

	for x := 0; x < 8; x++ {
		for y := 0; y < 8; y++ {
			squ := b.square(x, y)
			if squ.stone == 0 {
				turn++
			}
			if squ.stone != player {
				continue
			}
			playerStats.stoneCount++
			if squ.is(xSquare) {
				playerStats.xSquares++
			}
			if squ.is(sweet16) {
				playerStats.sweet16++
			}
			if squ.is(corner) {
				playerStats.corner++
			}
		}
	}

	if turn < 5 {
		return playerStats.stoneCount
	}
	return -playerStats.stoneCount + playerStats.sweet16 - 5*(playerStats.xSquares-playerStats.corner) + playerStats.corner
}

var corner = []square{
	square{x: 0, y: 0},
	square{x: 7, y: 0},
	square{x: 7, y: 7},
	square{x: 0, y: 7},
}
var sweet16 = []square{
	square{x: 2, y: 2},
	square{x: 2, y: 3},
	square{x: 2, y: 4},
	square{x: 2, y: 5},

	square{x: 3, y: 2},
	square{x: 3, y: 3},
	square{x: 3, y: 4},
	square{x: 3, y: 5},

	square{x: 4, y: 2},
	square{x: 4, y: 3},
	square{x: 4, y: 4},
	square{x: 4, y: 5},

	square{x: 5, y: 2},
	square{x: 5, y: 3},
	square{x: 5, y: 4},
	square{x: 5, y: 5},
}

var xSquare = []square{
	square{x: 1, y: 1},
	square{x: 6, y: 1},
	square{x: 6, y: 6},
	square{x: 1, y: 6},
}

func (s square) is(squares []square) bool {
	for _, squ := range squares {
		if s.x == squ.x && s.y == squ.y {
			return true
		}
	}
	return false
}
