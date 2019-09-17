package main

import (
	"fmt"
	"math/rand"
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
	Board    []string
	Turn     string
	Round    string
	GameOver bool
	p1Time   float64
	p2Time   float64
}

type Game struct {
	PlayerNum string
	EnemyNum  string
	Time      float64
	Ack       bool
	GameOver  bool
	Board     board
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
				myGame.PlayerNum = data[0]
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

			message := Message{
				Turn:   data[0],
				Round:  data[1],
				p1Time: p1t,
				p2Time: p2t,
			}

			myGame.Board = board{
				layout: data[4:68],
				turn:   message.Turn,
				round:  message.Round,
			}
			fmt.Printf("turn: %s\nround: %s\nboard: %s", message.Turn, message.Round, myGame.Board)
		}
	}
}

func (client *Client) SendMove(row, col int) {
	now := time.Now().UnixNano()
	if client.lastSent == 0 {
		client.lastSent = now
	}
	if now-client.lastSent < client.wait {
		return
	}
	client.lastSent = now
	move := fmt.Sprintf("%d\n%d", row, col)
	client.socket.Write([]byte(move)) //Send the move
	client.socket.Write([]byte("\n")) //Finish the message
	fmt.Printf("\nI played this move (%d, %d)\n", row, col)
}

func (client *Client) RequestUpdate() {
	client.socket.Write([]byte("\n")) //Request a new simple update
	fmt.Println("Requesting update from the server")
}

func main() {
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

	client := &Client{socket: connection, wait: 100000000}
	myGame := Game{}
	go client.receive(&myGame)
	timeout := int64(10000000000)
	for {
		if time.Now().UnixNano()-client.lastSent > timeout && client.lastSent != 0 {
			client.RequestUpdate() //need to jog the server
		}
		if myGame.GameOver {
			fmt.Println("Game Over")
			return
		}

		//Right now this is where the logic to send a turn lives, this is porbably no where it should live
		if myGame.Board.turn == myGame.PlayerNum && myGame.Ack && !myGame.GameOver {
			moves := myGame.Board.validMoves() //findMove(myGame.Board)
			move := square{
				x:     rand.Intn(7),
				y:     rand.Intn(7),
				stone: myGame.PlayerNum,
			}
			if len(moves) > 0 {
				move = moves[0]
			}
			client.SendMove(move.x, move.y)
		}
	}
}
