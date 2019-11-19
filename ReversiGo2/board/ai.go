package board

import (
	"math/rand"
	"sort"
)

//ValueBoard calculates a boards value using minimax and the Board.value() function
func ValueBoard(b *Board, player, depth int, currentBest float64, stop chan bool, value chan float64) (float64, bool) {
	select {
	case <-stop:
		if depth == 0 {
			score := b.Value(player)
			value <- score
			return score, true
		}
		return 0, false
	default:
		//Keep moving down the tree, we havent been stoped yet
		var target float64
		var prune bool
		moves := b.ValidMoves(player)
		if len(moves) == 0 {
			score := b.Value(player)
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

func ValueBoardDepth(b *Board, player, depth, maxDepth int, currentBest float64) float64 {
	if depth == maxDepth {
		return b.Value(player)
	}

	moves := b.ValidMoves(player)
	if len(moves) == 0 {
		return b.Value(player)
	}

	var target float64
	var prune bool
	nextCurrentBest := currentBest
	prunedMoves := pruneMoves(b, moves, depth, player, 3)

	for i, move := range prunedMoves {
		nb := move.board
		score := ValueBoardDepth(nb, player, depth+1, maxDepth, nextCurrentBest)
		target = miniMax(depth, score, target, i)
		nextCurrentBest, prune = alphBetaPrune(depth, currentBest, nextCurrentBest, score)
		if prune {
			break
		}
	}

	return target
}

func StocasticBestMove(b *Board, player int, moves []Move) Move {
	value := 0.0
	values := []float64{}
	total := 0.0
	move := moves[0]

	for _, m := range moves {
		x, y := m.XY()
		if x == 0 && y == 0 {
			return m
		}
		if x == 0 && y == 7 {
			return m
		}
		if x == 7 && y == 7 {
			return m
		}
		if x == 7 && y == 0 {
			return m
		}
		nb := b.Move(m)
		v := nb.Value(player)
		values = append(values, v)
		total += v
		if v > value {
			value = v
			move = m
		}
	}

	sel := rand.Intn(int(total))
	if float32(sel) < 0.99*float32(total) {
		return move
	}

	runningT := 0.0
	for i, v := range values {
		runningT += v
		if float64(sel) <= runningT {
			return moves[i]
		}
	}

	return move
}

func miniMax(depth int, score, target float64, i int) float64 {
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
	score float64
}

type byScore []valueBoard

func (a byScore) Len() int           { return len(a) }
func (a byScore) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a byScore) Less(i, j int) bool { return a[i].score < a[j].score }

func pruneMoves(b *Board, moves []Move, depth, player, max int) []valueBoard {
	ret := []valueBoard{}
	for _, move := range moves {
		nb := b.Move(move)
		score := nb.Value(player)
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

func alphBetaPrune(depth int, currentBest, nextCurrentBest, score float64) (float64, bool) {
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
	aSquares   int
}
