package main

import (
	"bufio"
	"fmt"
	"os/exec"
	"strconv"
	"strings"
)

type History struct {
	history []string
	size    int
}

func (h *History) add(s string) {
	if len(h.history) < h.size {
		h.history = append(h.history, s)
	}
	h.history = h.history[1:]
	h.history = append(h.history, s)
}

func (h *History) get(i int) string {
	index := h.size - i
	if index > len(h.history) || index < 0 {
		return ""
	}

	return h.history[index]
}

func main() {
	win := 0
	loss := 0
	tie := 0
	for i := 0; i < 25; i++ {
		score := runServer()
		if score == 1 {
			win++
		}
		if score == -1 {
			loss++
		}
		if score == 0 {
			tie++
		}
		fmt.Printf("\n\nRound: %d, Win: %d, Loss: %d, Tie: %d\n\n", i+1, win, loss, tie)
	}
}

func runServer() int {
	startServer := exec.Command("sh", "./ServerStart.sh")
	serverOut, err := startServer.StdoutPipe()
	if err != nil {
		fmt.Printf("Error stdOut: %s\n", err.Error())
		return 0
	}

	startJBot := exec.Command("sh", "./JavaStart.sh")
	startGBot := exec.Command("sh", "./GoStart.sh")
	err = startServer.Start()
	if err != nil {
		fmt.Printf("Error: %s\n", err.Error())
		return 0
	}

	scanner := bufio.NewScanner(serverOut)
	history := History{size: 10}
	for scanner.Scan() {
		out := scanner.Text()
		fmt.Print(".")
		history.add(out)
		switch out {
		case "Set up the connections:3334":
			fmt.Printf("start java\n")
			err = startJBot.Start()
			if err != nil {
				fmt.Printf("Error JBot: %s\n", err.Error())
				return 0
			}
		case "Set up the connections:3335":
			fmt.Printf("start golang\n")
			err = startGBot.Start()
			if err != nil {
				fmt.Printf("Error GBot: %s\n", err.Error())
				return 0
			}
		case "Game Over!":
			fmt.Println("exit the server")
			bScore, _ := strconv.ParseInt(history.get(7)[7:], 10, 32)
			wScore, _ := strconv.ParseInt(history.get(6)[7:], 10, 32)
			winMes := "White wins"
			score := 1
			if bScore > wScore {
				winMes = "Balck wins"
				score = -1
			}
			if bScore == wScore {
				winMes = "It's a Tie"
				score = 0
			}
			fmt.Printf("Game Results: B: %d, W: %d, %s\n", bScore, wScore, winMes)
			procs, _ := exec.Command("ps").Output()
			proc := ""
			for _, p := range strings.Split(string(procs), "\n") { //find the pid to kill the server process
				arr := strings.Split(p, " ")
				if len(arr) < 9 {
					continue
				}
				process1, process2, process3 := arr[7], arr[8], arr[6]
				if process1 == "/usr/bin/java" && process2 == "Reversi" {
					proc = arr[1]
					break
				}

				if process3 == "/usr/bin/java" && process1 == "Reversi" {
					proc = arr[0]
					break
				}
			}

			if proc == "" {
				fmt.Println("The proccess could not be found!")
				return 0
			} else {
				fmt.Printf("Killing proc %s\n", proc)
			}

			exec.Command("kill", proc).Run()
			return score
		}
	}
	return 0
}
