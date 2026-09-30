package game

import (
	"context"
	"fmt"
)

type nothing struct{}

type Game struct {
	currentFragment string
	wordList        []string
	validWords      map[string]nothing
	validFragments  map[string]nothing
}

type Move struct {
	Letter    rune
	Challenge bool
	Call      bool
}

func NewGame(wordList []string) *Game {
	return &Game{
		wordList: wordList,
	}
}

func (g *Game) Initialize(minLength int) {

	// Build the list of valid "loosing" words.
	// Words from the wordList that are as long or longer than the min length.
	g.validWords = map[string]nothing{}
	for _, word := range g.wordList {
		if len(word) >= minLength {
			g.validWords[word] = nothing{}
		}
	}

	// Build the valid fragment list
	// This is the list of all strings that are a prefix of a word in the validWord list
	// but are not themselves in the list.
	g.validFragments = map[string]nothing{}
	for word := range g.validWords {
		for i := 1; i < len(word); i++ {
			fragment := word[:i]
			_, exists := g.validWords[fragment]
			if !exists {
				g.validFragments[fragment] = nothing{}
			}

		}
	}

	// an empty string is a valid fragment
	g.validFragments[""] = nothing{}

}

func (g *Game) Play(c rune) {
	g.currentFragment = g.currentFragment + string(c)
}

func (g *Game) Evaluate(ctx context.Context) error {
	_, validWord := g.validWords[g.currentFragment]
	_, validFragment := g.validFragments[g.currentFragment]

	// This is a valid word. You lose. Not a strong position
	if validWord {
		return fmt.Errorf("not strong because it's a valid word: fragment=%s", g.currentFragment)
	}

	// This is not a valid fragment. The other play will challenge you. Not a strong position
	if !validFragment {
		return fmt.Errorf("not strong because this is not a valid fragment: fragment=%s", g.currentFragment)
	}

	suggestedMove := g.SuggestMove(ctx)
	// If the other player gives up (challenges), then we will win. We know this is valid fragment
	if !suggestedMove.Challenge {
		return fmt.Errorf("the other play will keep playing")
	}
	return nil

}

// Simulate just returneds a copy of the current game
// Note: The underlying word lists are not copied.
func (g *Game) Simulate() *Game {
	return &Game{
		currentFragment: g.currentFragment,
		wordList:        g.wordList,
		validWords:      g.validWords,
		validFragments:  g.validFragments,
	}

}

func (g *Game) SuggestMove(ctx context.Context) Move {

	if _, validWord := g.validWords[g.currentFragment]; validWord {
		return Move{Call: true}
	}

	for i := range 26 {
		letter := rune(97 + 25 - i)

		select {
		case <-ctx.Done():
			//fmt.Println("Oh no... ran out of time while thinking...")
			return Move{Letter: letter, Challenge: false}
		default:
		}

		simulation := g.Simulate()
		simulation.Play(letter)
		err := simulation.Evaluate(ctx)
		isStrongMove := err == nil
		if isStrongMove {
			return Move{
				Letter:    letter,
				Challenge: false,
			}
		}
	}

	// Can't win just give-up / challenge
	return Move{
		Challenge: true,
	}

}

func (g *Game) Challenge() (string, error) {
	if _, validWord := g.validWords[g.currentFragment]; validWord {
		return g.currentFragment, nil
	}

	for i := range 26 {
		letter := rune(i + 97)
		fragment := g.currentFragment + string(letter)
		if _, validWord := g.validWords[fragment]; validWord {
			return fragment, nil
		}
		if _, validFrfragment := g.validFragments[g.currentFragment]; !validFrfragment {
			continue
		}

		simulation := g.Simulate()
		simulation.Play(letter)
		word, err := simulation.Challenge()
		if err == nil {
			return word, nil
		}
	}

	return "", fmt.Errorf("Can't thing of a word!!")
}

func (g *Game) GetCurrentFragment() string {
	return g.currentFragment
}
