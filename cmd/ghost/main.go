package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"time"

	"github.com/ryanjyoder/ghost/game"
)

func main() {
	err := run()
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
		return
	}
}
func run() error {
	if len(os.Args) < 3 {
		return fmt.Errorf("first argument should be the wordlist filename")
	}

	var maxMs int64
	fmt.Sscanf(os.Args[2], "%d", &maxMs)

	fmt.Println("Reading Word list and initializing the game")
	newGame, err := game.LoadAndInitializeGame(os.Args[1])
	if err != nil {
		return err
	}
	fmt.Println("init complete")

	ctx := context.Background()

	scanner := bufio.NewScanner(os.Stdin)

	for {
		if newGame.GetCurrentFragment() != "" {
			fmt.Println("Current word:", newGame.GetCurrentFragment())
		}
		fmt.Println("Entter the letter you want to play (or press enter to challenge):")
		scanner.Scan()
		line := scanner.Text()
		if len(line) == 0 {
			word, err := newGame.Challenge()
			if err != nil {
				fmt.Println("You win!!")
				break
			}
			fmt.Println("My word:", word)
			break
		}
		newGame.Play(rune(line[0]))

		ctx, cancel := context.WithTimeout(ctx, time.Duration(maxMs)*time.Millisecond)
		move := newGame.SuggestMove(ctx)
		cancel()
		if move.Challenge {
			fmt.Println("Challenge!")
			break
		}
		if move.Call {
			fmt.Println("You lose! That's a word.")
			break
		}

		newGame.Play(move.Letter)
	}

	return nil
}
