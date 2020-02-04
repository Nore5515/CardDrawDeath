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

func main() {
	rand.Seed(time.Now().UTC().UnixNano())

	// Game Deck of 30 cards
	var deck [30]Card
	deck[29] = Card{"DRAW", 4}

	playing := true

	for playing{
		fmt.Println("Hi!");
		playing = false
	}

	c := Card{"MILL", 4}
	fmt.Println(c)
	fmt.Println(deck)

}



//fmt.Println("How many sides do you want your dice to have?")
//var sidePrompt string
//fmt.Scanln(&sidePrompt)
