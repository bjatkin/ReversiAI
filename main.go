package main

import (
	"fmt"
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
	EnemyNum  int
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
		if myGame.State.Turn == myGame.PlayerNum && len(myGame.State.Board) > 0 {
			fmt.Printf("It's my turn!... turn: %d, round: %d\nboard: %s\n", myGame.State.Turn, myGame.State.Round, myGame.State.Board)
			move := findMove(&myGame, myGame.State.Board)
			r, c := move.x, move.y
			client.SendMove(r, c)
			fmt.Printf("I played (%d, %d)\n", r, c)
		}
	}
}

type MovePos struct {
	x      int
	y      int
	player int
}

func findMove(Game *Game, board []string) MovePos {
	m := ValidMoves(Game.PlayerNum, board)
	vs := []int{}
	moves := []MovePos{}
	for _, i := range m {
		vs, moves = append(vs, valueBoard(Game, updateBoard(i, board), 5)), append(moves, i)
	}

	score := Max(vs)
	for i, j := range moves {
		if vs[i] == score {
			return j
		}
	}

	return MovePos{}
}

func valueBoard(Game *Game, board []string, depth int) int {
	moves := ValidMoves(Game.PlayerNum, board)
	if depth <= 0 || len(moves) == 0 {
		//calculate the leaf values
		return len(moves)
	}
	values := []int{}
	fmt.Printf("I'm gonna check some moves here %v\n", moves)
	for _, m := range moves {
		fmt.Printf("I'm here...\n")
		values = append(values, valueBoard(Game, updateBoard(m, board), depth-1))
	}
	fmt.Printf("I'm done checking moves now\n")

	fmt.Printf("\nValues: %v\n", values)
	if depth%2 == 0 {
		return Min(values)
	}
	return Max(values)
}

func Max(i []int) int {
	if len(i) == 0 {
		return 0
	}
	max := i[0]
	for _, v := range i {
		if v > max {
			max = v
		}
	}
	return max
}

func Min(i []int) int {
	if len(i) == 0 {
		return 0
	}
	min := i[0]
	for _, v := range i {
		if v < min {
			min = v
		}
	}
	return min
}

func updateBoard(m MovePos, board []string) []string {
	board[m.x*8+m.y] = string(m.player)
	return board
}

func ValidMoves(PlayerNum int, board []string) []MovePos {
	//Need to take into accoun the firest four rounds
	EnemyNum := 1
	if PlayerNum == 1 {
		EnemyNum = 2
	}
	moves := []MovePos{}

	// check for center moves
	if Pos(board, 3, 3) == "0" {
		moves = append(moves, MovePos{x: 3, y: 3, player: PlayerNum})
	}
	if Pos(board, 4, 3) == "0" {
		moves = append(moves, MovePos{x: 4, y: 3, player: PlayerNum})
	}
	if Pos(board, 3, 4) == "0" {
		moves = append(moves, MovePos{x: 3, y: 4, player: PlayerNum})
	}
	if Pos(board, 4, 4) == "0" {
		moves = append(moves, MovePos{x: 4, y: 4, player: PlayerNum})
	}
	for x := 0; x < 8; x++ {
		for y := 0; y < 8; y++ {
			if Pos(board, x, y) == strconv.Itoa(PlayerNum) {
				for _, m := range Ajacent(board, x, y) {
					val, _ := strconv.ParseInt(m.Val, 10, 32)
					if int(val) == EnemyNum {
						for m.Next(board) {
							if m.Val == "0" {
								//This is a valid move
								moves = append(moves, MovePos{
									x:      m.x,
									y:      m.y,
									player: PlayerNum,
								})
								break
							}

							if m.Val == strconv.Itoa(PlayerNum) {
								break
							}
						}
					}
				}
			}
		}
	}

	return moves
}

func Pos(board []string, x, y int) string {
	return board[x*8+y]
}

type Move struct {
	Val  string
	x    int
	y    int
	xDir int
	yDir int
}

func (m *Move) Next(board []string) bool {
	x := m.x + m.xDir
	y := m.y + m.yDir
	pos := x*8 + y
	if pos < 0 || pos > 63 || x < 0 || y < 0 {
		return false
	}
	m.Val = Pos(board, x, y)
	m.x = x
	m.y = y
	m.xDir = m.xDir
	m.yDir = m.yDir
	return true
}

func Ajacent(board []string, px, py int) []Move {
	ret := []Move{}
	for x := -1; x <= 1; x++ {
		for y := -1; y <= 1; y++ {
			if x == 0 && y == 0 {
				continue
			}
			x1, y1 := px+x, py+y
			if x1*8+y1 < 0 || x1*8+y1 > 63 {
				continue
			}

			ret = append(ret, Move{
				Val:  Pos(board, px+x, py+y),
				x:    px + x,
				y:    py + y,
				xDir: x,
				yDir: y,
			})
		}
	}

	return ret
}
