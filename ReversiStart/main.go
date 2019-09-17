package main

import (
	"bufio"
	"fmt"
	"os/exec"
	"strings"
)

func main() {
	startServer := exec.Command("sh", "./ServerStart.sh")
	serverOut, err := startServer.StdoutPipe()
	if err != nil {
		fmt.Printf("Error stdOut: %s\n", err.Error())
		return
	}

	startJBot := exec.Command("sh", "./JavaStart.sh")
	startGBot := exec.Command("sh", "./GoStart.sh")
	err = startServer.Start()
	if err != nil {
		fmt.Printf("Error: %s\n", err.Error())
		return
	}

	scanner := bufio.NewScanner(serverOut)
	for scanner.Scan() {
		out := scanner.Text()
		fmt.Println(out)
		switch out {
		case "Set up the connections:3334":
			fmt.Printf("start java\n")
			err = startJBot.Start()
			if err != nil {
				fmt.Printf("Error JBot: %s\n", err.Error())
				return
			}
		case "Set up the connections:3335":
			fmt.Printf("start golang\n")
			err = startGBot.Start()
			if err != nil {
				fmt.Printf("Error GBot: %s\n", err.Error())
				return
			}
		case "Game Over!":
			fmt.Println("exit the server")
			procs, _ := exec.Command("ps").Output()
			proc := ""
			for _, p := range strings.Split(string(procs), "\n") { //find the pid to kill the server process
				arr := strings.Split(p, " ")
				if len(arr) < 8 {
					continue
				}
				pid := arr[0]
				process1, process2 := arr[6], arr[7]
				if process1 == "/usr/bin/java" && process2 == "Reversi" {
					proc = pid
				}
			}

			if proc == "" {
				fmt.Println("The proccess could not be found!")
				return
			}
			exec.Command("kill", proc).Run()
			return
		}
	}
}
