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

type Player struct {
	name string 	//your name!
	hand []Card 	//your hand!
	drawing int  	//cards to draw at the end of your turn
	milling int 	//cards to mill before drawing on your turn
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

// Pass it however many cards you want to draw, and your hand
// Then the deck you want to draw it from
// It returns your updated hand with the randomly selected cards from the deck, then the deck without those cards in it
func drawCards(x int, hand []Card, deck []Card) ([]Card, []Card){
	rand.Seed(time.Now().UTC().UnixNano())
	z := -1

	for i := 0; i < x; i++{
		z = rand.Intn(len(deck))	//generates a random index of the deck
		hand = append(hand, deck[z])	//adds it to the slice
		deck = remove (deck[:], z)		//remove that card from the deck
	}

	return hand, deck
}

// remove X amount of cards from the deck, and put into the discard pile
// ...but right now, just remove, as we don't have a discard pile.
// return the deck.
func millCards(x int, deck []Card) ([]Card) {
	rand.Seed(time.Now().UTC().UnixNano())
	z := -1

	for i := 0; i < x; i++{
		z = rand.Intn(len(deck))	//generates a random index of the deck
		deck = remove (deck[:], z)		//remove that card from the deck
	}

	return deck
}

//works! returns the modified deck
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

//takes in a player and deck, and takes the player's turn. then returns the modified player and deck.
func takeTurn(p Player, d []Card) (Player, []Card) {
	// the player mills all cards he's supposed to first
	fmt.Println(p.name, "is milling", p.milling, "cards")	
	d = millCards(p.milling, d)
	// the player drawing all of his cards.
	fmt.Println(p.name, "is drawing", p.drawing, "cards")
	p.hand, d = drawCards(p.drawing, p.hand, d)
	// then resets his draw to 1, and mill to 0
	p.drawing = 1
	p.milling = 0
	return p, d
}

//take in a player and the card they want to play, then return the player
func playCard(p Player, c Card) (Player){
	if (c.action == "DRAW"){
		fmt.Println(p.name, "played a DRAW", c.count, "card!")
		p.drawing = p.drawing + c.count
	} else if (c.action == "SKIP"){
		fmt.Println(p.name, "played a SKIP", c.count, "card!")
		p.drawing = p.drawing - c.count
		if (p.drawing < 0){p.drawing = 0}
	} else if (c.action == "MILL"){
		fmt.Println(p.name, "played a MILL", c.count, "card!")	
		p.milling = p.milling + c.count
	}
	for i := 0; i < len(p.hand); i++{
		if (p.hand[i] == c){
			p.hand = remove(p.hand, i)
			i = len(p.hand)+1
		}
	}
	return p
}

func main() {
	rand.Seed(time.Now().UTC().UnixNano())

	// Game Deck of 30 cards
	//var deck [30]Card
	//deck[29] = Card{"DRAW", 4}
	deck := genDeck()
	p1 := Player{"Noah", []Card{},1, 0}


	//hand := []Card{}

	//playing := true

	//practice while loop
	/*
		fmt.Println("Hi!");
		playing = false
	}*/

	fmt.Println("DECK:", len(deck))
	p1, deck = takeTurn(p1, deck)
	fmt.Println(p1.hand)
	fmt.Println("DECK:", len(deck))
	p1 = playCard (p1, p1.hand[0])
	p1, deck = takeTurn(p1, deck)
	fmt.Println(p1.hand)
	fmt.Println("DECK:", len(deck))
}



//fmt.Println("How many sides do you want your dice to have?")
//var sidePrompt string
//fmt.Scanln(&sidePrompt)
