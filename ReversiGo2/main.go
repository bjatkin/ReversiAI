package main

import (
	rb "Projects/School/ReversiBot/ReversiGo2/board"
	rc "Projects/School/ReversiBot/ReversiGo2/client"
	rg "Projects/School/ReversiBot/ReversiGo2/generate"
	"fmt"
	"os"
	"strconv"
	"time"
)

func main() {
	b := rb.Board{
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 1, 0, 0, 1, 0, 0,
		2, 1, 1, 1, 1, 1, 1, 0,
		2, 1, 1, 1, 1, 1, 1, 0,
		2, 1, 2, 1, 2, 1, 0, 0,
		2, 2, 1, 1, 1, 1, 1, 1,
		2, 1, 2, 1, 1, 1, 1, 0,
		2, 2, 2, 2, 2, 2, 2, 2,
	}
	start := time.Now().Unix()
	count := 10000
	_, win, _ := rg.SlowScoreBoard(&b, 2, count)
	end := time.Now().Unix()
	fmt.Printf("success rate: %f\nwins: %d\nTime(sec): %d\n", float64(win)/float64(count), win, end-start)

	return

	if len(os.Args) < 3 {
		fmt.Printf("Please specify both the address of the server and the player number\n")
		return
	}

	fmt.Printf("connecting to the Reversi server @ %s...\n", os.Args[1])
	player, err := strconv.Atoi(os.Args[2])
	if err != nil || player < 0 || player > 2 {
		fmt.Printf("Player number must be 1 or 2: %s", err.Error())
		return
	}

	messages := make(chan rc.Message)
	client := rc.GetConnection(os.Args[1], player)
	go client.Receive(messages)

	for {
		select {
		case message := <-messages:
			//We got a message from the server
			if message.GameOver {
				fmt.Println("Game Over!")
				return
			}
			if message.Turn == player {
				//Get all the valid moves in the current board state
				board := rb.Board(message.Board)
				move, pass := findMove(&board, player, 5000*time.Millisecond)
				if pass {
					break //We have no valid moves
				}

				client.SendMove(move.XY())
			}
		}
	}
}

func findMove(b *rb.Board, player int, searchTime time.Duration) (rb.Move, bool) {
	moves := b.ValidMoves(player)
	if len(moves) == 0 {
		//No legal moves so we pass
		return rb.Move{}, true
	}
	if len(moves) == 1 {
		//No need to search this move, just play it
		return moves[0], false
	}

	stops := []chan bool{}
	scores := []chan int{}
	for i, m := range moves {
		stops = append(stops, make(chan bool, 1))
		scores = append(scores, make(chan int, 1))
		nb := b.Move(m)
		go rb.ValueBoard(&nb, player, 0, rb.MaxScore+1, stops[i], scores[i])
	}

	time.Sleep(searchTime)

	//Stop all the go routines and have them head back up
	for i := 0; i < len(moves); i++ {
		close(stops[i])
	}

	//Find the best scoring move
	max := rb.MinScore - 1
	move := moves[0]
	for i, m := range moves {
		score := <-scores[i]
		if score > max {
			max = score
			move = m
		}
		close(scores[i])
	}

	return move, false
}
