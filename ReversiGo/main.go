package main

import (
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
)

type Client struct {
	socket   net.Conn
	data     chan []byte
	lastSent int64
	wait     int64
}

type Message struct {
	Board    [64]int
	Turn     int
	Round    int
	GameOver bool
	p1Time   float64
	p2Time   float64
}

type Game struct {
	PlayerNum int
	EnemyNum  int
	Time      float64
	Ack       bool
	GameOver  bool
	Board     Board
}

func (client *Client) receive(myGame *Game) {
	for {
		message := make([]byte, 4096)
		length, err := client.socket.Read(message)
		if err != nil {
			client.socket.Close()
			break
		}

		//We got some incomming data
		if length > 0 {
			fmt.Println("reciving message")
			data := strings.Split(string(message), "\n")
			if !myGame.Ack {
				data := strings.Split(data[0], " ")
				pNum, err := strconv.ParseInt(data[0], 10, 32)
				myGame.PlayerNum = int(pNum)
				if err != nil {
					fmt.Printf("Invalid player number recieved: '%s'\n", data[0])
					client.socket.Close()
					break
				}
				myGame.Time, err = strconv.ParseFloat(data[1], 32)
				if err != nil {
					fmt.Printf("Invalid game time recieved: '%s'\n", data[1])
					client.socket.Close()
					break
				}
				myGame.Ack = true
				continue
			}

			if data[0] == "-999" {
				myGame.GameOver = true
				return //The game is over, stop reading stuff
			}
			// update the game state
			turn, err := strconv.ParseInt(data[0], 10, 32)
			if err != nil {
				fmt.Printf("Game state not updated, invalid turn recieved '%s', Error: %s\n", data[0], err.Error())
				continue
			}
			round, err := strconv.ParseInt(data[1], 10, 32)
			if err != nil {
				fmt.Printf("Game state not updated, invalid round recieved '%s', Error: %s\n", data[1], err.Error())
				continue
			}
			p1t, err := strconv.ParseFloat(data[2], 32)
			if err != nil {
				fmt.Printf("Game state not updated, invalid p1Time recieved: '%s', Error: %s\n", data[2], err.Error())
				continue
			}
			p2t, err := strconv.ParseFloat(data[3], 32)
			if err != nil {
				fmt.Printf("Game state not updated, invalid p2Time recieved: '%s', Error: %s\n", data[3], err.Error())
				continue
			}
			//convert the board
			board := [64]int{}
			// row := 0
			for i, s := range data[4:68] {
				d, err := strconv.ParseInt(s, 10, 8)
				if err != nil {
					fmt.Printf("Game state not updated, invalid board recieved: '%+v', Error: %s\n", data, err.Error())
				}
				board[i] = int(d)
				// fmt.Printf("Index: %d, %d, %d\n", ((8-(i%8))+row*8)-1, i%8, row)
				// board[63-(((8-(i%8))+row*8)-1)] = int(d)
				// if (i+1)%8 == 0 {
				// 	row++
				// }
			}

			message := Message{
				Board:  board,
				Turn:   int(turn),
				Round:  int(round),
				p1Time: p1t,
				p2Time: p2t,
			}

			myGame.Board = Board{
				layout: message.Board,
				turn:   message.Turn,
				round:  message.Round,
			}
			fmt.Printf("turn: %d\nround: %d\nboard: %s", message.Turn, message.Round, myGame.Board)
		}
	}
}

func (client *Client) SendMove(row, col, round int) {
	now := time.Now().UnixNano()
	if client.lastSent == 0 {
		client.lastSent = now - client.wait - 1
	}
	if now-client.lastSent < client.wait {
		return
	}
	client.lastSent = now
	move := fmt.Sprintf("%d\n%d", row, col)
	client.socket.Write([]byte(move)) //Send the move
	client.socket.Write([]byte("\n")) //Finish the message
	fmt.Printf("\n-------------------------------\nPlayed Move: (%d, %d) ROUND: %d\n", row, col, round)
}

func (client *Client) RequestUpdate() {
	now := time.Now().UnixNano()
	if client.lastSent == 0 {
		client.lastSent = now - client.wait - 1
	}
	client.lastSent = now
	client.socket.Write([]byte("-1\n-1\n")) //Request a new simple update
	fmt.Println("Requesting update from the server")
}

func main() {
	// b := Board{
	// 	layout: [64]int{
	// 		0, 0, 0, 0, 0, 0, 0, 0,
	// 		0, 0, 0, 0, 0, 0, 0, 0,
	// 		0, 0, 0, 0, 0, 0, 0, 0,
	// 		0, 0, 0, 0, 0, 0, 0, 0,
	// 		0, 0, 0, 0, 0, 0, 0, 0,
	// 		0, 0, 0, 0, 0, 0, 0, 0,
	// 		0, 0, 0, 0, 0, 0, 0, 0,
	// 		0, 0, 0, 0, 0, 0, 0, 0,
	// 	},
	// 	turn:  1,
	// 	round: 1,
	// }
	// findMove(&b, 2)
	// return
	if len(os.Args) < 3 {
		fmt.Printf("Please specify both the address of the server and the player number\n")
		return
	}

	fmt.Printf("connecting to the Reversi server @ %s...\n", os.Args[1])
	port, err := strconv.Atoi(os.Args[2])
	if err != nil {
		fmt.Printf("Player number must be 1 or 2: %s", err.Error())
	}
	connection, err := net.Dial("tcp", fmt.Sprintf("%s:%d", os.Args[1], 3333+port))
	if err != nil {
		fmt.Printf("There was an error connecting to the server: %s\n", err.Error())
		return
	}
	sec := 1000000000.0 //one second in nano seconds
	client := &Client{socket: connection, wait: int64(0.01 * sec)}
	myGame := Game{}
	go client.receive(&myGame)
	timeout := int64(5 * sec)
	currentRound := 1
	for {
		if time.Now().UnixNano()-client.lastSent > timeout && client.lastSent != 0 {
			client.RequestUpdate() //need to jog the server
			currentRound++
		}
		if myGame.GameOver {
			fmt.Println("Game Over")
			return
		}

		//Right now this is where the logic to send a turn lives, this is porbably no where it should live
		if myGame.Board.turn == myGame.PlayerNum && myGame.Ack && !myGame.GameOver && myGame.Board.round == currentRound {
			fmt.Printf("Search for a round %d move\n", currentRound)
			nb := NewBoard(&myGame.Board) //do this to prevent weird errors with the board updating
			nb.turn = myGame.PlayerNum
			move := findMove(&nb, 7)
			client.SendMove(move.squares[0].x, move.squares[0].y, currentRound)
			currentRound += 2
		}
	}
}
