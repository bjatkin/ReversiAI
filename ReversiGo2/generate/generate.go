package generate

import (
	rb "Projects/School/ReversiBot/ReversiGo2/board"
	"math/rand"
	"time"
)

//StartBoard is a blank starting board
var StartBoard = rb.Board{
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
}

//Add this so we make sure every run of Game is really really random
var rando int64 = 0

//Game takes a game position and plays it randomly to the
//	end returning all positions and the winner
func Game(game *rb.Board, player int) ([]*rb.Board, int) {
	p1 := 1
	p2 := 2
	if player == 2 {
		p1 = 2
		p2 = 1
	}

	p1Moves := []rb.Move{}
	p2Moves := []rb.Move{}
	ret := []*rb.Board{}
	//get rand set up
	rand.Seed(time.Now().Unix() + rando)
	rando++
	for run := true; run; run = (len(p1Moves)+len(p2Moves) > 0) {
		p1Moves = game.ValidMoves(p1)
		p1MCount := len(p1Moves)
		//Play a random move
		if p1MCount > 0 {
			move := p1Moves[rand.Intn(p1MCount)]
			newBoard := game.Move(move)
			ret = append(ret, &newBoard)
			game = &newBoard
		}

		p2Moves = game.ValidMoves(p2)
		p2MCount := len(p2Moves)
		//Play a random move
		if p2MCount > 0 {
			move := p2Moves[rand.Intn(p2MCount)]
			newBoard := game.Move(move)
			ret = append(ret, &newBoard)
			game = &newBoard
		}
	}
	p1Score := 0
	p2Score := 0
	for x := 0; x < 8; x++ {
		for y := 0; y < 8; y++ {
			stone := game.Square(x, y).Stone
			if stone == p1 {
				p1Score++
			}
			if stone == p2 {
				p2Score++
			}
		}
	}
	if p1Score == p2Score {
		return ret, 0
	}

	if p1Score > p2Score {
		return ret, p1
	}

	return ret, p2
}

//SlowScoreBoard takes a board position and play it randomly to the end
//	itter times and returns the black wins, white wins and ties
func SlowScoreBoard(b *rb.Board, player, itter int) (int, int, int) {
	black := 0
	white := 0
	tie := 0
	for i := 0; i < itter; i++ {
		_, score := Game(b, player)
		if score == 0 {
			tie++
		}
		if score == 1 {
			black++
		}
		if score == 2 {
			white++
		}
	}

	return black, white, tie
}
