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
	myGame, err := game.LoadAndInitializeGame(os.Args[1])
	if err != nil {
		return err
	}
	fmt.Println("init complete")

	ctx := context.Background()

	scanner := bufio.NewScanner(os.Stdin)

	for {
		if myGame.GetCurrentFragment() != "" {
			fmt.Println("Current word:", myGame.GetCurrentFragment())
		}

		// Ask the user for their next letter
		fmt.Println("Entter the letter you want to play (or press enter to challenge):")
		scanner.Scan()
		line := scanner.Text()
		if len(line) == 0 {
			word, err := myGame.Challenge()
			if err != nil {
				fmt.Println("You win!!")
				break
			}
			fmt.Println("My word:", word)
			break
		}

		// Play the letter they gave
		myGame.Play(rune(line[0]))

		// Ask the computer for their next move
		ctx, cancel := context.WithTimeout(ctx, time.Duration(maxMs)*time.Millisecond)
		move := myGame.SuggestMove(ctx)
		cancel()
		if move.Challenge {
			fmt.Println("Challenge!")
			break
		}
		if move.Call {
			fmt.Println("You lose! That's a word.")
			break
		}

		myGame.Play(move.Letter)
	}

	return nil
}
