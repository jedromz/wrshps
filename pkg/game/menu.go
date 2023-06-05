package game

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"runtime"
)

func clear() {
	var c *exec.Cmd
	if runtime.GOOS == "windows" {
		c = exec.Command("cmd", "/c", "cls")
	} else {
		c = exec.Command("clear")
	}
	c.Stdout = os.Stdout
	c.Run()
}

func (a *App) Menu(ctx context.Context) {
	for {
		clear()

		fmt.Println("Welcome to Warships!")
		fmt.Println("1. Start bot game")
		fmt.Println("2. Start PvP game")
		fmt.Println("3. Enter player info")
		fmt.Println("4. Show player ranking")
		fmt.Println("5. Show player stats")
		fmt.Println("6. Show Player Lobby")
		fmt.Println("7. Exit")
		fmt.Println("0. Return to menu")

		var choice int
		_, err := fmt.Scanln(&choice)
		if err != nil {
			fmt.Println("An error occurred:", err)
			continue
		}

		switch choice {
		case 1:
			a.StartBotGame(ctx)
		case 2:
			a.StartPlayerGame(ctx)
		case 3:
			a.EnterPlayerInfo(ctx)
		case 4:
			stats, err := a.game.GetTopPlayerStats()
			if err != nil {
				fmt.Println("An error occurred:", err)
				continue
			}
			fmt.Println(stats)
		case 5:
			a.GetPlayerStats(ctx)
		case 7:
			fmt.Println("Bye!")
			os.Exit(0)
		case 6:
			a.PrintLobby()
		default:
			fmt.Println("Invalid option. Please enter a number between 0 and 6.")
		}
		fmt.Println("Press ane key to continue...")
		fmt.Scanln()

	}
}
