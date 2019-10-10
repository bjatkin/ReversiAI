package board

import (
	"sort"
)

//ValueBoard calculates a boards value using minimax and the Board.value() function
func ValueBoard(b *Board, player, depth, currentBest int, stop chan bool, value chan int) (int, bool) {
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
		var target int
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
