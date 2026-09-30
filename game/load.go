package game

import (
	"bufio"
	"fmt"
	"os"
)

const (
	DefaultMinLength = 4
)

func LoadAndInitializeGame(filename string) (*Game, error) {
	words, err := loadWordList(filename)
	if err != nil {
		return nil, err
	}

	newGame := NewGame(words)
	newGame.Initialize(DefaultMinLength)

	return newGame, nil

}

func loadWordList(filename string) ([]string, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("could not open file: filename:%s %w", filename, err)
	}

	result := []string{}

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		result = append(result, line)
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading file:%w", err)
	}

	return result, nil
}
