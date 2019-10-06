package generate

import (
	rb "Projects/School/ReversiBot/ReversiGo2/board"
	"bufio"
	"fmt"
	"io"
	"math/rand"
	"os"
	"strconv"
	"strings"
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

//Add this to make random things more random
var rando int64

//Game takes a game position and plays it randomly to the
//	end returning all positions and the winner
func Game(game *rb.Board, player int) PlayedGame {
	p1 := 1
	p2 := 2
	if player == 2 {
		p1 = 2
		p2 = 1
	}

	p1Moves := []rb.Move{}
	p2Moves := []rb.Move{}
	ret := [64]*rb.Board{}
	//get rand set up
	rand.Seed(time.Now().Unix() + rando)
	rando++

	var i int
	for run := true; run; run = (len(p1Moves)+len(p2Moves) > 0) {
		p1Moves = game.ValidMoves(p1)
		p1MCount := len(p1Moves)
		//Play a random move
		if p1MCount > 0 {
			move := p1Moves[rand.Intn(p1MCount)]
			newBoard := game.Move(move)
			ret[i] = &newBoard
			i++
			game = &newBoard
		}

		p2Moves = game.ValidMoves(p2)
		p2MCount := len(p2Moves)
		//Play a random move
		if p2MCount > 0 {
			move := p2Moves[rand.Intn(p2MCount)]
			newBoard := game.Move(move)
			ret[i] = &newBoard
			i++
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
		return PlayedGame{ret, 0}
	}

	if p1Score > p2Score {
		return PlayedGame{ret, p1}
	}

	return PlayedGame{ret, p2}
}

//WinLossTie is the number of black wins, white wins and ties
type WinLossTie struct {
	black, white, tie int
}

//SlowScoreBoard takes a board position and play it randomly to the end
//	itter times and returns the black wins, white wins and ties
func SlowScoreBoard(b *rb.Board, player, itter int) WinLossTie {
	ret := WinLossTie{}
	for i := 0; i < itter; i++ {
		game := Game(b, player)
		if game.winner == 0 {
			ret.tie++
		}
		if game.winner == 1 {
			ret.black++
		}
		if game.winner == 2 {
			ret.white++
		}
	}
	return ret
}

//PlayedGame is a set of the 64 boards that make up a game
type PlayedGame struct {
	positions [64]*rb.Board
	winner    int
}

//SaveGames saves a set of games into a file
func SaveGames(games []PlayedGame, dest io.Writer) {
	gameCount := len(games)
	for line := 0; line < gameCount; line++ {
		header := fmt.Sprintf("Game(%d) Winner(%d)\n", line, games[line].winner)
		var game string
		for _, pos := range games[line].positions {
			if pos != nil {
				game += fmt.Sprintf("%v\n", [64]int(*pos))
			}
		}
		dest.Write([]byte(header + game))
	}
}

//CreateSavedGames generates count games and saves them to the file specified
func CreateSavedGames(count int, fileName string) {
	var allGames []PlayedGame
	for x := 0; x < count; x++ {
		allGames = append(allGames, Game(&StartBoard, 1))
		if x%1000 == 0 {
			fmt.Printf("generating game %d\n", x)
		}
	}

	//Save all the games I just played
	file, err := os.Create(fileName)
	if err != nil {
		fmt.Printf("Error Creating Save File%s\n", err.Error())
		return
	}
	defer file.Close()

	fmt.Printf("Saving Games...\n")
	SaveGames(allGames, file)
}

//SaveScoredPositions get's positions from a file and scores them saving them to a separate file
func SaveScoredPositions(dataFile, saveFile string, skip, maxRoutines int, append bool, startAt string) error {
	rand.Seed(time.Now().Unix() + rando)
	rando++
	data, err := os.Open(dataFile)
	if err != nil {
		fmt.Printf("Error Opening data file %s\n", err.Error())
		return err
	}
	defer data.Close()

	var dest *os.File
	if !append {
		dest, err = os.Create(saveFile)
	} else {
		dest, err = os.OpenFile(saveFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	}
	if err != nil {
		fmt.Printf("Error Creating Save File %s\n", err.Error())
		return err
	}
	defer dest.Close()
	if !append {
		dest.WriteString("Game_Position, 1_Black_Wins, 1_White_Wins, 1_Tie, 2_Black_Wins, 2_White_Wins, 2_Tie\n")
	}

	scanner := bufio.NewScanner(data)
	count := 0
	routineCount := 0
	type posScore struct {
		wlt [2]WinLossTie
		pos rb.Board
	}
	scores := make(chan posScore, routineCount*2)
	start := false
	for scanner.Scan() {
		line := scanner.Text()

		//should we start yet?
		if !start {
			if line == startAt {
				start = true
			}
			continue
		}
		//Game line contains no position
		if line[0] == 'G' {
			continue
		}

		//First 4 lines are always the same and can be easily manually added
		if line == "[0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 1 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0]" ||
			line == "[0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 1 0 0 0 0 0 0 0 2 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0]" ||
			line == "[0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 1 0 0 0 0 0 0 0 2 1 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0]" ||
			line == "[0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 1 2 0 0 0 0 0 0 2 1 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0]" {
			continue
		}

		if rand.Intn(skip) != 1 {
			continue
		}

		board, err := boardString(line)
		if err != nil {
			fmt.Printf("Error creating board %s\n", err.Error())
			continue
		}

		if routineCount > maxRoutines {
			select {
			case pScore := <-scores:
				routineCount--
				dest.WriteString(fmt.Sprintf("%v,%d,%d,%d,%d,%d,%d\n",
					[64]int(pScore.pos),
					pScore.wlt[0].black,
					pScore.wlt[0].white,
					pScore.wlt[0].tie,
					pScore.wlt[1].black,
					pScore.wlt[1].white,
					pScore.wlt[1].tie,
				))
			}
		}

		fmt.Printf("Scoring Position #%d...\n", count)
		go func(s chan posScore) {
			data := posScore{}
			data.wlt[0] = SlowScoreBoard(&board, 1, 1000)
			data.wlt[1] = SlowScoreBoard(&board, 2, 1000)
			data.pos = board
			s <- data
		}(scores)
		routineCount++

		count++
	}

	for routineCount > 0 {
		select {
		case pScore := <-scores:
			routineCount--
			dest.Write([]byte(fmt.Sprintf("%v,%d,%d,%d,%d,%d,%d\n",
				[64]int(pScore.pos),
				pScore.wlt[0].black,
				pScore.wlt[0].white,
				pScore.wlt[0].tie,
				pScore.wlt[1].black,
				pScore.wlt[1].white,
				pScore.wlt[1].tie,
			)))
		}
	}

	return nil
}

func boardString(pos string) (rb.Board, error) {
	ret := rb.Board{}
	var err error

	stones := strings.Split(pos[1:128], " ")
	for i, stone := range stones {
		ret[i], err = strconv.Atoi(stone)
		if err != nil {
			return ret, err
		}
	}
	return ret, nil
}
