package main

import (
	"fmt"
	"strconv"
	"os"
	"math/rand"
	"time"
	"github.com/gorilla/mux"
	"net/http"
	"log"
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

type Data struct {
	p1 Player
	p2 Player
	deck []Card
}

var gameData Data

// WEB STUFF ALERT WOOP WOOP
func get(w http.ResponseWriter, r *http.Request) {
		bs := []byte(strconv.Itoa(len(gameData.p1.hand)))
		w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    w.Write([]byte(bs))
}

func put(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusAccepted)
    w.Write([]byte(`{"message": "put called"}`))
}
// WEB STUFF IS GONE

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
// Also! It returns true lastly if it worked, and false if you drew a non-existant card
func drawCards(x int, hand []Card, deck []Card) ([]Card, []Card, bool){
	rand.Seed(time.Now().UTC().UnixNano())
	z := -1

	for i := 0; i < x; i++{
		if (len(deck) == 0){
			fmt.Println("^^^LOSER^^^!")
			return hand, deck, false
		}
		z = rand.Intn(len(deck))	//generates a random index of the deck
		hand = append(hand, deck[z])	//adds it to the slice
		deck = remove (deck[:], z)		//remove that card from the deck
	}

	return hand, deck, true
}

// remove X amount of cards from the deck, and put into the discard pile
// ...but right now, just remove, as we don't have a discard pile.
// return the deck, and true if it works (false if it didn't)
func millCards(x int, deck []Card) ([]Card, bool) {
	rand.Seed(time.Now().UTC().UnixNano())
	z := -1

	for i := 0; i < x; i++{
		if (len(deck) == 0){
			fmt.Println("^^^LOSER^^^!")
			return deck, false
		}
		z = rand.Intn(len(deck))	//generates a random index of the deck
		deck = remove (deck[:], z)		//remove that card from the deck
	}

	return deck, true
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
//ALSO return true if it worked, and false if it didn't
func takeTurn(p Player, d []Card) (Player, []Card, bool) {
	notDead := true
	// the player mills all cards he's supposed to first
	fmt.Println(p.name, "is milling", p.milling, "cards")
	d, notDead = millCards(p.milling, d)
	// the player drawing all of his cards.
	if (notDead){
		fmt.Println(p.name, "is drawing", p.drawing, "cards")
		p.hand, d, notDead = drawCards(p.drawing, p.hand, d)
	}
	// then resets his draw to 1, and mill to 0

	p.drawing = 1
	p.milling = 0
	return p, d, notDead
}

//take in a player and the card they want to play, then return the player
func playCard(p Player, victim Player, c Card) (Player, Player){
	if (c.action == "DRAW"){			//FORCES OPPONENT TO DRAW THAT MANY
		fmt.Println(p.name, "played a DRAW", c.count, "card!")
		victim.drawing = victim.drawing + c.count
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
	return p, victim
}

//ANN VARIABLE POTENTIAL
// When to Mill (aiMillSensitivity, 0 means mill whenever possible, higher value means only mill when safer)
// When to Skip (aiSkipSensitivity, 0 means they skip only if they'll draw the last card in the deck, higher value means they skip to a higher base deck value(???))
// When to Attack (aiDrawSensitvity, 0 means they will attack only enough for a lethal hit, higher value means more overkill)
// The AI should attack when they have enough cards to kill, mill when they don't, and skip if they're drawing more than they can handle
func aiPlay(ai Player, victim Player, deck []Card) (Player, Player){
	maxDraw, maxMill, maxSkip := 0, 0, 0		//the max amount that the AI can draw, mill and skip on its turn
	aiDrawSensitvity, aiMillSensitivity, aiSkipSensitivity := 0, 0, 0
	for i := 0; i < len(ai.hand); i++{
		if (ai.hand[i].action == "DRAW"){maxDraw += ai.hand[i].count}
		if (ai.hand[i].action == "MILL"){maxMill += ai.hand[i].count}
		if (ai.hand[i].action == "SKIP"){maxSkip += ai.hand[i].count}
	}

	//if it is going to die, skip until it can live! or at least try to.
	for (ai.drawing > len(deck) + aiSkipSensitivity && maxSkip > 0){
		maxSkip -= nextCardWith(ai, "SKIP").count
		ai, victim = playCard(ai, victim, nextCardWith(ai, "SKIP"))
	}

	//if it's in no danger of dying, mill the deck down to get it closer to 0.
	//so, if it's milling less cards than are remaining in the deck after it draws, it'll play some mill cards
	for (ai.milling < len(deck)-aiMillSensitivity-ai.drawing && maxMill > 0){
		//SO LONG AS PLAYING THE CARD WOULD NOT PUSH AI.MILLING OVER LEN(deck)-ai.drawing, IT WILL PLAY
		//otherwise it'd just commit suicide
		if (len(deck)-ai.drawing-ai.milling-nextCardWith(ai, "MILL").count >= 0){
			maxMill -= nextCardWith(ai, "MILL").count
			ai, victim = playCard(ai, victim, nextCardWith(ai, "MILL"))
		} else{
			maxMill = 0
		}
	}

	//finally, if they feel safe in an attack, they'll go for it!
	if ((len(deck) + aiDrawSensitvity) - victim.drawing - victim.milling - maxDraw  < 0){
		//fmt.Println("(", len(deck), "+", aiDrawSensitvity, ") -", victim.drawing, "-", victim.milling, "-", maxDraw, " = ", (len(deck) + aiDrawSensitvity) - victim.drawing - victim.milling - maxDraw)
		for (maxDraw > 0){
			maxDraw -= nextCardWith(ai, "DRAW").count
			ai, victim = playCard(ai, victim, nextCardWith(ai, "DRAW"))
		}
	}

	return ai, victim
}

func nextCardWith(p Player, search string) Card{
	for i := 0; i < len(p.hand); i++{
		if (p.hand[i].action == search){
			return p.hand[i]
		}
	}
	return Card{"",0}
}

func main() {
	rand.Seed(time.Now().UTC().UnixNano())

	//deck := genDeck()
	//gameData.p1 := Player{"Noah", []Card{},1, 0}

	// WEB ALERT WOOP WOOP
	// routers are...uh...important for web stuff?
	// then take in the GET function (seen above) and "handle" it...
	// then start the server!
	r := mux.NewRouter()


	// Main game loop
	webInitialized := false
	playing := true
	inGame := false
	for playing {
		fmt.Print(":")
		var response string
		response = "exit" 	//default response
		fmt.Scanln(&response)	//get player's response
		// If you want to quit...
		if (response == "exit" || response == "quit"){playing = false}
		// If you want to start a new game...
		if (response == "start" || response == "new" || response == "begin"){
			fmt.Println("What is your name?")
			fmt.Scanln(&response)	//get player's response
			gameData.p1 = Player{response, []Card{}, 3, 0}	//note! player starts with draw 3, for a starting hand size.
			fmt.Println("Let's begin!\n====================")
			// Game actually begins here!
			fmt.Println("Shuffling deck...")
			deck := genDeck()									//create deck
			fmt.Println("Drawing hand...")
			gameData.p1, deck, _ = takeTurn(gameData.p1, deck)			//draw hand
			fmt.Println("Generating opponent...")
			p2 := Player{"roBob", []Card{}, 3, 0}		//creates a second player, AI opponent NOTE: current AI just plays first card in it's hand
			p2, deck, _ = takeTurn(p2, deck)
			fmt.Println("Beginning game loop...\n\n\n\n")
			inGame = true
			for inGame {
				fmt.Println("===============\nDECK:",len(deck))
				fmt.Println("Your Hand:", gameData.p1.hand)
				// TESTING
				r.HandleFunc("/", get).Methods(http.MethodGet)
				//go log.Fatal(http.ListenAndServe(":8080", r))
				// blocking funvytion must be in coroutine
				// alo breaks on second loop
				if (!webInitialized){
						webInitialized = true
						go func() { log.Fatal(http.ListenAndServe(":8080", r)) }()
				}
				// TESTING
				fmt.Println("Your Opponent's Hand:", len(p2.hand))
				fmt.Println(gameData.p1.name, "will Draw", gameData.p1.drawing, "and Mill", gameData.p1.milling)
				fmt.Println(p2.name, "will Draw", p2.drawing, "and Mill", p2.milling)
				fmt.Println("---------------")
				fmt.Println("Play a card?")
				response = "no"		//keeps the default to no, so if u just hit enter it ends ur turn
				fmt.Scanln(&response)	//get player's response
				if (response == "no" || response == "next" || response == "done" || response == "skip" || response == "next"){	//if you wanna end your turn
					gameData.p1, deck, inGame = takeTurn(gameData.p1, deck)
					if (!inGame){break;}	//if you lose, break the game loop
					fmt.Println("DECK:",len(deck))
					fmt.Println("---------------\nOPPONENTS TURN")
					p2, gameData.p1 = aiPlay(p2, gameData.p1, deck)
					//p2, gameData.p1 = playCard(p2, gameData.p1, p2.hand[0])
					p2, deck, inGame = takeTurn(p2, deck)
					if (!inGame){break;}	//if he loses, break the game loop
				} else if (response == "exit" || response == "quit"){	//if you wanna quit
					inGame = false
				} else if (response == "play" || response == "yes"){	//if you wanna play a card
					fmt.Println("What is the index of the card you want to play?")
					fmt.Scanln(&response)	//get player's response
					gameData.p1, p2 = playCard(gameData.p1, p2, gameData.p1.hand[stoI(response)])
				} else if _, err := strconv.Atoi(response); err == nil {
					gameData.p1, p2 = playCard(gameData.p1, p2, gameData.p1.hand[stoI(response)])
				}
			}
			// End of game loop.

		}
	}
}



//fmt.Println("How many sides do you want your dice to have?")
//var sidePrompt string
//fmt.Scanln(&sidePrompt)
