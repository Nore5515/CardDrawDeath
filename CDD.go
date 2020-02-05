package main

import (
	"fmt"
	"strconv"
	"os"
	"math/rand"
	"time"
)

type Card struct {
	action string
	count int
}

func stoI (s string) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
    }
    return i
}

func seeHand () string {
	return ""
}

//creates the DEFAULT deck
func genDeck() []Card{
	//var deck [30]Card
	//updated to use slices
	deck := make([]Card,30, 100)

	for i := 0; i < 10; i++{deck[i] = Card{"DRAW", 1}}
	for i := 10; i < 15; i++{deck[i] = Card{"MILL", 3}}
	for i := 15; i < 20; i++{deck[i] = Card{"SKIP", 2}}
	for i := 20; i < 23; i++{deck[i] = Card{"DRAW", 5}}
	for i := 23; i < 26; i++{deck[i] = Card{"SKIP", 5}}
	for i := 26; i < 28; i++{deck[i] = Card{"MILL", 6}}
	for i := 28; i < 29; i++{deck[i] = Card{"DRAW", 10}}
	for i := 29; i < 30; i++{deck[i] = Card{"SKIP", 8}}

	return deck
}

// Pass it however many cards you want to draw
// Then the deck you want to draw it from
// It returns a slice of Cards, randomly selected from the deck
func drawCards(x int, deck []Card) []Card{
	chosenCards := make([]Card, x)
	rand.Seed(time.Now().UTC().UnixNano())
	z := -1

	for i := 0; i < x; i++{
			z = rand.Intn(len(deck))
			chosenCards[i] = deck[z]
			deck[z] = Card{"", -1}
	}

	return chosenCards
}

func remove (deck []Card, index int) []Card{
	if (len(deck) > 0){
		deck[index] = deck[len(deck)-1]
		deck[len(deck)-1] = Card{"", -1}
		deck = deck[:len(deck)-1]
		return deck
	} else{
		fmt.Println("Phbbt")
		fmt.Println("0 element array in remove!")
		var x []Card
		return x
	}
}

func main() {
	rand.Seed(time.Now().UTC().UnixNano())

	// Game Deck of 30 cards
	//var deck [30]Card
	//deck[29] = Card{"DRAW", 4}
	deck := genDeck()

	//playing := true

	//practice while loop
	/*
		fmt.Println("Hi!");
		playing = false
	}*/

	c := Card{"MILL", 4}
	fmt.Println(c)
	fmt.Println(deck)
	deck = remove (deck[:], 0)
	fmt.Println(deck)

}



//fmt.Println("How many sides do you want your dice to have?")
//var sidePrompt string
//fmt.Scanln(&sidePrompt)
