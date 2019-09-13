package main

import (
	"fmt"
	"math/rand"
	"net"
	"os"
	"strconv"
	"strings"
)

type Client struct {
	socket net.Conn
	data   chan []byte
}

type Message struct {
	Board    []string
	Turn     int
	GameOver bool
	Round    int
	p1Time   float64
	p2Time   float64
}

type Game struct {
	PlayerNum int
	Time      float64
	Ack       bool
	State     Message
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
			data := strings.Split(string(message), "\n")
			if !myGame.Ack {
				myGame.Ack = true
				data := strings.Split(data[0], " ")
				PlayerNum, err := strconv.Atoi(data[0])
				if err != nil {
					fmt.Printf("Invalid player number: '%s'\n", data[1])
					client.socket.Close()
					break
				}
				myGame.PlayerNum = PlayerNum
				myGame.Time, err = strconv.ParseFloat(data[1], 32)
				if err != nil {
					fmt.Printf("Invalid game time recieved: '%s'\n", data[1])
					client.socket.Close()
					break
				}
				continue
			}

			// update the game state
			turn, err := strconv.Atoi(data[0])
			if turn == -999 {
				myGame.State.GameOver = true
				return //The game is over, stop reading stuff
			}
			if err != nil {
				fmt.Printf("Game state not updated, invalid turn recieved: '%s', Error: %s\n", data[0], err.Error())
				continue
			}
			round, err := strconv.Atoi(data[1])
			if err != nil {
				fmt.Printf("Game state not updated, invalid round recieved: '%s', Error: %s\n", data[1], err.Error())
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

			myGame.State = Message{
				Board:  data[4:68],
				Turn:   turn,
				Round:  round,
				p1Time: p1t,
				p2Time: p2t,
			}
		}
	}
}

func (client *Client) SendMove(row, col int) {
	move := fmt.Sprintf("%d\n%d", row, col)
	client.socket.Write([]byte(move)) //Send the move
	client.socket.Write([]byte("\n")) //Finish the message
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

	client := &Client{socket: connection}
	myGame := Game{}
	go client.receive(&myGame)
	for {
		if myGame.State.GameOver {
			fmt.Println("Game Over")
			return
		}

		//Right now this is where the logic to send a turn lives, this is porbably no where it should live
		if myGame.State.Turn == myGame.PlayerNum {
			fmt.Printf("It's my turn!... turn: %d, round: %d\n", myGame.State.Turn, myGame.State.Round)
			r, c := rand.Intn(8), rand.Intn(8)
			client.SendMove(r, c)
			fmt.Printf("I played (%d, %d)\n", r, c)
		}
	}
}
