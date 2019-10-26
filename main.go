package main

import (
	"fmt"
	"strings"

	"github.com/mrpineapples/deck"
)

// Hand represents a players hand.
type Hand []deck.Card

func (h Hand) String() string {
	strs := make([]string, len(h))
	for i := range h {
		strs[i] = h[i].String()
	}
	return strings.Join(strs, ", ")
}

// DealerString returns the dealers hand, only showing the first card.
func (h Hand) DealerString() string {
	return h[0].String() + ", **HIDDEN**"
}

// Score returns the current players score.
// It handles the values of Ace's if necessary.
func (h Hand) Score() int {
	minScore := h.MinScore()
	if minScore > 11 {
		return minScore
	}

	for _, c := range h {
		if c.Rank == deck.Ace {
			// At this point Ace == 1 we add 10 to make it equal 11
			return minScore + 10
		}
	}

	return minScore
}

// MinScore returns the minimum possible score (i.e. Ace is always 1).
func (h Hand) MinScore() int {
	score := 0
	for _, c := range h {
		score += min(int(c.Rank), 10)
	}
	return score
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// Shuffle combines three decks and randomizes the order.
func Shuffle(gs GameState) GameState {
	ret := clone(gs)
	ret.Deck = deck.New(deck.Deck(3), deck.Shuffle)
	return ret
}

// Deal is called at the start of a game and gives each player two cards.
func Deal(gs GameState) GameState {
	ret := clone(gs)
	ret.Player = make(Hand, 0, 5)
	ret.Dealer = make(Hand, 0, 5)
	var card deck.Card
	for i := 0; i < 2; i++ {
		card, ret.Deck = draw(ret.Deck)
		ret.Player = append(ret.Player, card)
		card, ret.Deck = draw(ret.Deck)
		ret.Dealer = append(ret.Dealer, card)
	}
	ret.State = StatePlayerTurn
	return ret
}

// Hit takes a card from the deck and adds it to the current hand.
func Hit(gs GameState) GameState {
	ret := clone(gs)
	hand := ret.CurrentPlayer()
	var card deck.Card
	card, ret.Deck = draw(ret.Deck)
	*hand = append(*hand, card)
	if hand.Score() > 21 {
		Stand(ret)
	}
	return ret
}

// Stand ends the current players turn.
func Stand(gs GameState) GameState {
	ret := clone(gs)
	ret.State++
	return ret
}

// EndGame displays the score and winner/loser of the current game.
func EndGame(gs GameState) GameState {
	ret := clone(gs)
	pScore, dScore := ret.Player.Score(), ret.Dealer.Score()
	fmt.Println("==FINAL HANDS==")
	fmt.Println("Player:", ret.Player, "\nScore:", pScore)
	fmt.Println("Dealer:", ret.Dealer, "\nScore:", dScore)
	switch {
	case pScore > 21:
		fmt.Println("You busted! You lose 😢")
	case dScore > 21:
		fmt.Println("Dealer busted! You win 🎉")
	case pScore > dScore:
		fmt.Println("You win 🎉")
	case dScore > pScore:
		fmt.Println("You lose 😢")
	case dScore == pScore:
		fmt.Println("Draw!")
	}
	fmt.Println()
	ret.Player = nil
	ret.Dealer = nil
	return ret
}

func main() {
	var gs GameState
	gs = Shuffle(gs)

	// Play 5 games of blackjack
	for i := 0; i < 5; i++ {
		gs = Deal(gs)

		var input string
		for gs.State == StatePlayerTurn {
			fmt.Println("Player:", gs.Player)
			fmt.Println("Dealer:", gs.Dealer.DealerString())
			fmt.Println("What will you do? (h)it, (s)tand")
			fmt.Scanf("%s\n", &input)
			switch input {
			case "h":
				gs = Hit(gs)
			case "s":
				gs = Stand(gs)
			default:
				if input != "s" {
					fmt.Printf("\"%s\" is not a valid option 🤕. Try again.\n\n", input)
				}
			}
		}

		for gs.State == StateDealerTurn {
			if gs.Dealer.Score() <= 16 || (gs.Dealer.Score() == 17 && gs.Dealer.MinScore() != 17) {
				gs = Hit(gs)
			} else {
				gs = Stand(gs)
			}
		}

		gs = EndGame(gs)
	}
}

func draw(cards []deck.Card) (deck.Card, []deck.Card) {
	return cards[0], cards[1:]
}

// State represents the current phase of the game.
type State int8

// The three game states.
const (
	StatePlayerTurn = iota
	StateDealerTurn
	StateHandOver
)

// GameState is an object which can represent the game's state at any given momemnt.
type GameState struct {
	Deck   []deck.Card
	State  State
	Player Hand
	Dealer Hand
}

// CurrentPlayer a pointer to the hand of the current player.
func (gs *GameState) CurrentPlayer() *Hand {
	switch gs.State {
	case StatePlayerTurn:
		return &gs.Player
	case StateDealerTurn:
		return &gs.Dealer
	default:
		panic("It's currently not any player's turn 🤨")
	}
}

func clone(gs GameState) GameState {
	ret := GameState{
		Deck:   make([]deck.Card, len(gs.Deck)),
		State:  gs.State,
		Player: make(Hand, len(gs.Player)),
		Dealer: make(Hand, len(gs.Dealer)),
	}
	copy(ret.Deck, gs.Deck)
	copy(ret.Player, gs.Player)
	copy(ret.Dealer, gs.Dealer)
	return ret
}
