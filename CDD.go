package main

import (
	"fmt"
	"strconv"
	"os"
	"math/rand"
	"time"
	"github.com/gorilla/mux"
	"github.com/gorilla/handlers"
	"net/http"
	"log"
	"encoding/json"
)

type Card struct {
	Action string
	Count int
}

type Player struct {
	Name string 	`json:"name"`//your Name!
	Hand []Card 	//your Hand!
	Drawing int  	//cards to draw at the end of your turn
	Milling int 	//cards to mill before Drawing on your turn
}

type Data struct {
	P1 Player		`json:"play1"`
	P2 Player		`json:"play2"`
	Deck []Card	`json:"deck"`
}

var gameData Data

// WEB STUFF ALERT WOOP WOOP
func get(w http.ResponseWriter, r *http.Request) {
		// TODO: Find a way to convert Data into a json stringy thing
		//bs := []byte(strconv.Itoa(len(gameData.P1.Hand)))
		bs, _ := json.Marshal(gameData)
		w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    w.Write([]byte(bs))
		fmt.Println("GETTING")
}

// enables...cors?
func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
}

func put(w http.ResponseWriter, r *http.Request) {
		fmt.Println("PUTTING")
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

//creates the DEFAULT Deck
func genDeck() []Card{
	//var Deck [30]Card
	//updated to use slices
	Deck := make([]Card,30, 100)

	for i := 0; i < 10; i++{Deck[i] = Card{"DRAW", 1}}
	for i := 10; i < 15; i++{Deck[i] = Card{"MILL", 3}}
	for i := 15; i < 20; i++{Deck[i] = Card{"SKIP", 2}}
	for i := 20; i < 23; i++{Deck[i] = Card{"DRAW", 5}}
	for i := 23; i < 26; i++{Deck[i] = Card{"SKIP", 5}}
	for i := 26; i < 28; i++{Deck[i] = Card{"MILL", 6}}
	for i := 28; i < 29; i++{Deck[i] = Card{"DRAW", 10}}
	for i := 29; i < 30; i++{Deck[i] = Card{"SKIP", 8}}

	return Deck
}

// Pass it however many cards you want to draw, and your Hand
// Then the Deck you want to draw it from
// It returns your updated Hand with the randomly selected cards from the Deck, then the Deck without those cards in it
// Also! It returns true lastly if it worked, and false if you drew a non-existant card
func drawCards(x int, Hand []Card, Deck []Card) ([]Card, []Card, bool){
	rand.Seed(time.Now().UTC().UnixNano())
	z := -1

	for i := 0; i < x; i++{
		if (len(Deck) == 0){
			fmt.Println("^^^LOSER^^^!")
			return Hand, Deck, false
		}
		z = rand.Intn(len(Deck))	//generates a random index of the Deck
		Hand = append(Hand, Deck[z])	//adds it to the slice
		Deck = remove (Deck[:], z)		//remove that card from the Deck
	}

	return Hand, Deck, true
}

// remove X amount of cards from the Deck, and put into the discard pile
// ...but right now, just remove, as we don't have a discard pile.
// return the Deck, and true if it works (false if it didn't)
func millCards(x int, Deck []Card) ([]Card, bool) {
	rand.Seed(time.Now().UTC().UnixNano())
	z := -1

	for i := 0; i < x; i++{
		if (len(Deck) == 0){
			fmt.Println("^^^LOSER^^^!")
			return Deck, false
		}
		z = rand.Intn(len(Deck))	//generates a random index of the Deck
		Deck = remove (Deck[:], z)		//remove that card from the Deck
	}

	return Deck, true
}

//works! returns the modified Deck
func remove (Deck []Card, index int) []Card{
	if (len(Deck) > 0){
		Deck[index] = Deck[len(Deck)-1]
		Deck[len(Deck)-1] = Card{"", -1}
		Deck = Deck[:len(Deck)-1]
		return Deck
	} else{
		fmt.Println("Phbbt")
		fmt.Println("0 element array in remove!")
		var x []Card
		return x
	}
}

//takes in a player and Deck, and takes the player's turn. then returns the modified player and Deck.
//ALSO return true if it worked, and false if it didn't
func takeTurn(p Player, d []Card) (Player, []Card, bool) {
	notDead := true
	// the player mills all cards he's supposed to first
	fmt.Println(p.Name, "is Milling", p.Milling, "cards")
	d, notDead = millCards(p.Milling, d)
	// the player Drawing all of his cards.
	if (notDead){
		fmt.Println(p.Name, "is Drawing", p.Drawing, "cards")
		p.Hand, d, notDead = drawCards(p.Drawing, p.Hand, d)
	}
	// then resets his draw to 1, and mill to 0

	p.Drawing = 1
	p.Milling = 0
	return p, d, notDead
}

//take in a player and the card they want to play, then return the player
func playCard(p Player, victim Player, c Card) (Player, Player){
	if (c.Action == "DRAW"){			//FORCES OPPONENT TO DRAW THAT MANY
		fmt.Println(p.Name, "played a DRAW", c.Count, "card!")
		victim.Drawing = victim.Drawing + c.Count
	} else if (c.Action == "SKIP"){
		fmt.Println(p.Name, "played a SKIP", c.Count, "card!")
		p.Drawing = p.Drawing - c.Count
		if (p.Drawing < 0){p.Drawing = 0}
	} else if (c.Action == "MILL"){
		fmt.Println(p.Name, "played a MILL", c.Count, "card!")
		p.Milling = p.Milling + c.Count
	}
	for i := 0; i < len(p.Hand); i++{
		if (p.Hand[i] == c){
			p.Hand = remove(p.Hand, i)
			i = len(p.Hand)+1
		}
	}
	return p, victim
}

//ANN VARIABLE POTENTIAL
// When to Mill (aiMillSensitivity, 0 means mill whenever possible, higher value means only mill when safer)
// When to Skip (aiSkipSensitivity, 0 means they skip only if they'll draw the last card in the Deck, higher value means they skip to a higher base Deck value(???))
// When to Attack (aiDrawSensitvity, 0 means they will attack only enough for a lethal hit, higher value means more overkill)
// The AI should attack when they have enough cards to kill, mill when they don't, and skip if they're Drawing more than they can Handle
func aiPlay(ai Player, victim Player, Deck []Card) (Player, Player){
	maxDraw, maxMill, maxSkip := 0, 0, 0		//the max amount that the AI can draw, mill and skip on its turn
	aiDrawSensitvity, aiMillSensitivity, aiSkipSensitivity := 0, 0, 0
	for i := 0; i < len(ai.Hand); i++{
		if (ai.Hand[i].Action == "DRAW"){maxDraw += ai.Hand[i].Count}
		if (ai.Hand[i].Action == "MILL"){maxMill += ai.Hand[i].Count}
		if (ai.Hand[i].Action == "SKIP"){maxSkip += ai.Hand[i].Count}
	}

	//if it is going to die, skip until it can live! or at least try to.
	for (ai.Drawing > len(Deck) + aiSkipSensitivity && maxSkip > 0){
		maxSkip -= nextCardWith(ai, "SKIP").Count
		ai, victim = playCard(ai, victim, nextCardWith(ai, "SKIP"))
	}

	//if it's in no danger of dying, mill the Deck down to get it closer to 0.
	//so, if it's Milling less cards than are remaining in the Deck after it draws, it'll play some mill cards
	for (ai.Milling < len(Deck)-aiMillSensitivity-ai.Drawing && maxMill > 0){
		//SO LONG AS PLAYING THE CARD WOULD NOT PUSH AI.Milling OVER LEN(Deck)-ai.Drawing, IT WILL PLAY
		//otherwise it'd just commit suicide
		if (len(Deck)-ai.Drawing-ai.Milling-nextCardWith(ai, "MILL").Count >= 0){
			maxMill -= nextCardWith(ai, "MILL").Count
			ai, victim = playCard(ai, victim, nextCardWith(ai, "MILL"))
		} else{
			maxMill = 0
		}
	}

	//finally, if they feel safe in an attack, they'll go for it!
	if ((len(Deck) + aiDrawSensitvity) - victim.Drawing - victim.Milling - maxDraw  < 0){
		//fmt.Println("(", len(Deck), "+", aiDrawSensitvity, ") -", victim.Drawing, "-", victim.Milling, "-", maxDraw, " = ", (len(Deck) + aiDrawSensitvity) - victim.Drawing - victim.Milling - maxDraw)
		for (maxDraw > 0){
			maxDraw -= nextCardWith(ai, "DRAW").Count
			ai, victim = playCard(ai, victim, nextCardWith(ai, "DRAW"))
		}
	}

	return ai, victim
}

func nextCardWith(p Player, search string) Card{
	for i := 0; i < len(p.Hand); i++{
		if (p.Hand[i].Action == search){
			return p.Hand[i]
		}
	}
	return Card{"",0}
}

func main() {
	rand.Seed(time.Now().UTC().UnixNano())

	//Deck := genDeck()
	//gameData.P1 := Player{"Noah", []Card{},1, 0}

	// WEB ALERT WOOP WOOP
	// routers are...uh...important for web stuff?
	// then take in the GET function (seen above) and "Handle" it...
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
			fmt.Println("What is your Name?")
			fmt.Scanln(&response)	//get player's response
			gameData.P1 = Player{response, []Card{}, 3, 0}	//note! player starts with draw 3, for a starting Hand size.
			fmt.Println("Let's begin!\n====================")
			// Game actually begins here!
			fmt.Println("Shuffling Deck...")
			gameData.Deck = genDeck()									//create Deck
			fmt.Println("Drawing Hand...")
			gameData.P1, gameData.Deck, _ = takeTurn(gameData.P1, gameData.Deck)			//draw Hand
			fmt.Println("Generating opponent...")
			gameData.P2 = Player{"roBob", []Card{}, 3, 0}		//creates a second player, AI opponent NOTE: current AI just plays first card in it's Hand
			gameData.P2, gameData.Deck, _ = takeTurn(gameData.P2, gameData.Deck)
			fmt.Println("Beginning game loop...\n\n\n\n")
			inGame = true
			for inGame {
				fmt.Println("===============\nDeck:",len(gameData.Deck))
				fmt.Println("Your Hand:", gameData.P1.Hand)
				// TESTING
				r.HandleFunc("/", get).Methods(http.MethodGet)
				//r.HandleFunc("/player", getAllPlayers).Methods(http.MethodGet)
				//r.HandleFunc("/player/:playerName", getPlayerByName(name)).Methods(http.MethodGet)
				r.HandleFunc("/", put).Methods(http.MethodPut)
				//go log.Fatal(http.ListenAndServe(":8080", r))
				// blocking funvytion must be in coroutine
				// alo breaks on second loop
				if (!webInitialized){
						webInitialized = true
						headersOk := handlers.AllowedHeaders([]string{"X-Requested-With"})
						originsOk := handlers.AllowedOrigins([]string{"*"})
						methodsOk := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "OPTIONS"})
						cors := handlers.CORS(originsOk, headersOk, methodsOk)
						go func() { log.Fatal(http.ListenAndServe(":8080", cors(r))) }()
				}
				// TESTING
				fmt.Println("Your Opponent's Hand:", len(gameData.P2.Hand))
				fmt.Println(gameData.P1.Name, "will Draw", gameData.P1.Drawing, "and Mill", gameData.P1.Milling)
				fmt.Println(gameData.P2.Name, "will Draw", gameData.P2.Drawing, "and Mill", gameData.P2.Milling)
				fmt.Println("---------------")
				fmt.Println("Play a card?")
				response = "no"		//keeps the default to no, so if u just hit enter it ends ur turn
				fmt.Scanln(&response)	//get player's response
				if (response == "no" || response == "next" || response == "done" || response == "skip" || response == "next"){	//if you wanna end your turn
					gameData.P1, gameData.Deck, inGame = takeTurn(gameData.P1, gameData.Deck)
					if (!inGame){break;}	//if you lose, break the game loop
					fmt.Println("Deck:",len(gameData.Deck))
					fmt.Println("---------------\nOPPONENTS TURN")
					gameData.P2, gameData.P1 = aiPlay(gameData.P2, gameData.P1, gameData.Deck)
					//P2, gameData.P1 = playCard(P2, gameData.P1, P2.Hand[0])
					gameData.P2, gameData.Deck, inGame = takeTurn(gameData.P2, gameData.Deck)
					if (!inGame){break;}	//if he loses, break the game loop
				} else if (response == "exit" || response == "quit"){	//if you wanna quit
					inGame = false
				} else if (response == "play" || response == "yes"){	//if you wanna play a card
					fmt.Println("What is the index of the card you want to play?")
					fmt.Scanln(&response)	//get player's response
					gameData.P1, gameData.P2 = playCard(gameData.P1, gameData.P2, gameData.P1.Hand[stoI(response)])
				} else if _, err := strconv.Atoi(response); err == nil {
					gameData.P1, gameData.P2 = playCard(gameData.P1, gameData.P2, gameData.P1.Hand[stoI(response)])
				}
			}
			// End of game loop.

		}
	}
}



//fmt.Println("How many sides do you want your dice to have?")
//var sidePrompt string
//fmt.Scanln(&sidePrompt)
