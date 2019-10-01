package board

import (
	"sort"
)

//ValueBoard calculates a boards value using minimax and the Board.value() function
func ValueBoard(b *Board, player, depth, currentBest int, stop chan bool, value chan int) (int, bool) {
	select {
	case <-stop:
		if depth == 0 {
			score := b.value(player)
			value <- score
			return score, true
		}
		return 0, false
	default:
		//Keep moving down the tree, we havent been stoped yet
		var target int
		var prune bool
		moves := b.ValidMoves(player)
		if len(moves) == 0 {
			score := b.value(player)
			if depth == 0 {
				value <- score
			}
			return score, true
		}

		nextCurrentBest := currentBest
		prunedMoves := pruneMoves(b, moves, depth, player, 4)

		for i, move := range prunedMoves {
			nb := move.board
			score, calculated := ValueBoard(nb, player, depth+1, nextCurrentBest, stop, value)
			if !calculated {
				score = move.score
			}
			target = miniMax(depth, score, target, i)
			nextCurrentBest, prune = alphBetaPrune(depth, currentBest, nextCurrentBest, score)
			if prune { //We can stop searching due to alpha beta pruning
				break
			}
		}

		if depth == 0 {
			value <- target
		}

		return target, true
	}
}

func miniMax(depth, score, target, i int) int {
	if i == 0 {
		return score
	}

	if depth%2 == 0 && score < target { //Min
		return score
	}
	if depth%2 != 0 && score > target { //Max
		return score
	}
	return target
}

type valueBoard struct {
	board *Board
	move  *Move
	score int
}

type byScore []valueBoard

func (a byScore) Len() int           { return len(a) }
func (a byScore) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a byScore) Less(i, j int) bool { return a[i].score < a[j].score }

func pruneMoves(b *Board, moves []Move, depth, player, max int) []valueBoard {
	ret := []valueBoard{}
	for _, move := range moves {
		nb := b.Move(move)
		score := nb.value(player)
		ret = append(ret, valueBoard{&nb, &move, score})
	}
	if max == 0 || len(ret) <= max {
		return ret
	}

	//Sort the valueBoards
	sort.Sort(byScore(ret))

	//Min layer so return the min
	if depth%2 == 0 {
		return ret[:max]
	}

	//Max layer so return the max
	return ret[len(ret)-max:]
}

func alphBetaPrune(depth, currentBest, nextCurrentBest, score int) (int, bool) {
	//adjust the currentBestValue
	if depth%2 == 0 && currentBest == MaxScore+1 {
		currentBest = MinScore - 1
	}

	//I'm min and above/ bellow is a max
	if depth%2 == 0 && score < currentBest {
		return 0, true //Prune all remaining nodes
	}
	if depth%2 == 0 {
		return score, false
	}

	//I'm a max and above/ bellow is a min
	if depth%2 != 0 && score > currentBest {
		return 0, true
	}
	return score, false
}

type boardStats struct {
	win        bool
	loss       bool
	stoneCount int
	frontier   int
	sweet16    int
	top4       int
	edges      int
	corner     int
	mobility   int
	xSquares   int
	cSquares   int
	bSquares   int
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
			if squ.is(top4) {
				playerStats.top4++
			}
			if squ.is(edge) {
				playerStats.edges++
			}
			if squ.isFrontier(&b) {
				playerStats.frontier++
			}
		}
	}
	playerStats.mobility = len(b.ValidMoves(player))

	if turn < 5 { //End game
		return playerStats.stoneCount
	}

	if turn > 40 { //Early Game
		return -playerStats.stoneCount +
			playerStats.mobility +
			5*playerStats.top4 +
			playerStats.sweet16 +
			-playerStats.edges +
			-playerStats.frontier +
			-500*playerStats.xSquares
	}

	return -playerStats.stoneCount +
		playerStats.sweet16 +
		playerStats.top4 +
		-500*(playerStats.xSquares-playerStats.corner) +
		3*playerStats.corner +
		-3*playerStats.frontier
}

var edge = []square{
	square{x: 0, y: 2},
	square{x: 0, y: 3},
	square{x: 0, y: 4},
	square{x: 0, y: 5},

	square{x: 7, y: 2},
	square{x: 7, y: 3},
	square{x: 7, y: 4},
	square{x: 7, y: 5},

	square{x: 2, y: 0},
	square{x: 3, y: 0},
	square{x: 4, y: 0},
	square{x: 5, y: 0},

	square{x: 2, y: 7},
	square{x: 3, y: 7},
	square{x: 4, y: 7},
	square{x: 5, y: 7},
}

var corner = []square{
	square{x: 0, y: 0},
	square{x: 7, y: 0},
	square{x: 7, y: 7},
	square{x: 0, y: 7},
}

var top4 = []square{
	square{x: 4, y: 4},
	square{x: 3, y: 3},
	square{x: 4, y: 3},
	square{x: 3, y: 4},
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

func (s square) isFrontier(b *Board) bool {
	adj := s.adj()
	for _, ray := range adj {
		if b.square(ray.destx, ray.desty).stone == 0 {
			return true
		}
	}
	return false
}
