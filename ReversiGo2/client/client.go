package client

import (
	rb "Projects/School/ReversiBot/ReversiGo2/board"
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"
)

//Client is a client that connects to a reversi server
type Client struct {
	socket     net.Conn
	LastSent   int64
	LastUpdate int64
	Ack        bool
}

//Message is an incomming message from the reversi server
type Message struct {
	Board    [64]int
	Turn     int
	Round    int
	GameOver bool
	p1Time   float64
	p2Time   float64
}

//GetConnection returns a client connected to a reversi server
func GetConnection(dest string, player int) *Client {
	connection, err := net.Dial("tcp", fmt.Sprintf("%s:%d", dest, 3333+player))
	if err != nil {
		fmt.Printf("There was an error connecting to the server: %s\n", err.Error())
		return &Client{}
	}
	return &Client{socket: connection}
}

//Receive sends incomming messages through the injected channel
func (client *Client) Receive(message chan Message) {
	for {
		incomming := make([]byte, 4096)
		length, err := client.socket.Read(incomming)
		if err != nil {
			client.socket.Close()
			break
		}

		//We got some incomming data
		if length > 0 {
			fmt.Println("reciving message")
			data := strings.Split(string(incomming), "\n")
			if !client.Ack {
				data := strings.Split(data[0], " ")
				pNum, err := strconv.ParseInt(data[0], 10, 32)
				if err != nil {
					fmt.Printf("Invalid player number recieved: '%s'\n", data[0])
					client.socket.Close()
					break
				}
				gameTime, err := strconv.ParseFloat(data[1], 32)
				if err != nil {
					fmt.Printf("Invalid game time recieved: '%s'\n", data[1])
					client.socket.Close()
					break
				}
				client.Ack = true
				fmt.Printf("My Player: %d\nGame Time: %f\n", pNum, gameTime)
				continue
			}

			if data[0] == "-999" {
				message <- Message{
					GameOver: true,
				}
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

			board := [64]int{}
			for i, s := range data[4:68] {
				d, err := strconv.ParseInt(s, 10, 8)
				if err != nil {
					fmt.Printf("Game state not updated, invalid board recieved: '%+v', Error: %s\n", data, err.Error())
				}
				board[i] = int(d)
			}

			message <- Message{
				Board:  board,
				Turn:   int(turn),
				Round:  int(round),
				p1Time: p1t,
				p2Time: p2t,
			}

			client.LastUpdate = time.Now().UnixNano()
			fmt.Printf("turn: %d\nround: %d\nBoard: %s", int(turn), int(round), rb.Board(board))
		}
	}
}

//SendMove sends a move to the Reversi server
func (client *Client) SendMove(x, y int) {
	move := fmt.Sprintf("%d\n%d", x, y)
	client.socket.Write([]byte(move + "\n")) //Send the move
}
