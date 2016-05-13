package main 

import (
	"fmt"
	"math/rand"
	"time"
	"os"
	"sort"
)

// All types 

type Location struct {
	Description string
	Transitions []string
	Events []string
}

type Event struct {
	Type string
	Chance int
	Description string
	Health int
	Evt string
}

type Game struct {
	Welcome string
	Health int
	Current Location string
}

type Character struct {
	Name string
	Health int
	Evasion int
	Alive bool
	Speed int
	Weap int
	Npc bool
}

type Weapon struct {
	minAtt int
	maxAtt int
	Name string
}

type Players []Character

var locationMap = map[string]*Location{
	// various locations of the game
	"Bridge":      {"You are on the bridge of a spaceship sitting in the Captain's chair.", []string{"Ready Room", "Turbo Lift"}, []string{"alienAttack"}},
	"Ready Room":  {"The Captain's ready room.", []string{"Bridge"}, []string{}},
	"Turbo Lift":  {"A Turbo Lift that takes you anywhere in the ship.", []string{"Bridge", "Lounge", "Engineering"}, []string{"android"}},
	"Engineering": {"You are in engineering where you see the star drive", []string{"Turbo Lift"}, []string{"alienAttack"}},
	"Lounge":      {"You are in the lounge, you feel very relaxed", []string{"Turbo Lift"}, []string{"relaxing"}},
}

var evts = map[string]*Event{
	// Possible events based on location
	"alienAttack":     {Chance: 20, Description: "An alien beams in front of you and shoots you with a ray gun.", Health: -50, Evt: "doctorTreatment"},
	"doctorTreatment": {Chance: 10, Description: "The doctor rushes in and inject you with a health boost.", Health: +30, Evt: ""},
	"android":         {Chance: 50, Description: "Data is in the turbo lift and says hi to you", Health: 0, Evt: ""},
	"relaxing":        {Chance: 100, Description: "In the lounge you are so relaxed that your health improves.", Health: +10, Evt: ""},
}

var ennemies = map[int]*Character{
	1: {Name: "Klingon", Health: 50, Alive: true, Weap: 2},
	2: {Name: "Romulan", Health: 55, Alive: true, Weap: 3},
}

var Weaps = map[int]*Weapon{
	// 3 types of weapons
	1: {Name: "Phaser", minAtt: 5, maxAtt: 15},
	2: {Name: "Klingon Disruptor", minAtt: 1, maxAtt: 15},
	3: {Name: "Romulan Disruptor", minAtt: 3, maxAtt: 12},
}

func (e *Event) ProcessEvent() int {
	// event struct
	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)

	if e.Chance >= r1.Intn(100) {
		hp := e.Health
		if e.Type == "Combat" {
			fmt.Println("Combat Event")
		}

		fmt.Printf("\t%s\n", e.Description)

		if e.Evt !+ "" {
			hp = hp + evts[e.Evt].ProcessEvent()
		}
		return hp
	}

	return 0

}

func (p *Character) Equip(w int) {
	// Equip weapon selected
	p.Weap = w
}

func (p *Character) Attack() int {
	// UAttack function
	return Weaps[p.Weap].Fire()
}

func (w *Weapon) Fire() int {
	// corresponds to combat with weapons
	return w.minAtt + rand.Intn(w.maxAtt - w.minAtt)
}

// Players
func (slice Players) Len() int {
	return len(slice)
}

func (slice Players) Less(i, j int) bool {
	return slice[i].Speed > slice[j].Speed
}

func (slice Players) Swap(i, j int) {
	slice[i], slice[j] = slice[j], slice[i]
}

// Game logic
func (g *Game) Play() {
    CurrentLocation = locationMap["Bridge"]
	fmt.Println(g.Welcome)

    for {
		fmt.Println(CurrentLocation.Description)
		g.ProcessEvents(CurrentLocation.Events)

		if g.Health <= 0 {
			fmt.Println("You are dead, game over!!!")
			return
		}

		// Print Health
		fmt.Printf("Health: %d\n", g.Health)
		fmt.Println("You can go to these places:")

		for index, loc := range CurrentLocation.Transitions {
			fmt.Printf("\t%d - %s\n", index+1, loc)
		}

		i := 0

		// User inputted event
		for i < 1 || i > len(CurrentLocation.Transitions) {
			fmt.Printf("%s%d%s\n", "Where do you want to go (0 - to quit), [1...", len(CurrentLocation.Transitions), "]: ")
			fmt.Scan(&i)
		}
		newLoc := i - 1   

		// Go to user specified location                                         
		CurrentLocation = locationMap[CurrentLocation.Transitions[newLoc]]
	}
}

//Combat functions
func RunBattle(players Players) {
    sort.Sort(players)
    
	round := 1
	numAlive := players.Len()
	playerAction := 0
	for {
	    for x := 0; x < players.Len(); x++ {
    		players[x].Evasion = 0      // Reset evasion for all characters
    	}
		DisplayInfo("Combat round", round, "begins...")
        for x := 0; x < players.Len(); x++ {
            if players[x].Alive != true {
                continue
            }
            playerAction = 0
            if !players[x].Npc {
                DisplayInfo("DO you want to")
                DisplayInfo("\t1 - Run")
                DisplayInfo("\t2 - Evade")
                DisplayInfo("\t3 - Attack")
                GetUserInput(&playerAction)
            }
            if playerAction == 2 {
                players[x].Evasion = rand.Intn(15)
                DisplayInfo("Evasion set to:", players[x].Evasion)
            }
            tgt := selectTarget(players, x)
            if tgt != -1 {
                DisplayInfo("player: ", x, "target: ", tgt)
                attp1 := players[x].Attack()
                players[tgt].Health = players[tgt].Health - attp1
                if players[tgt].Health <= 0 {
                    players[tgt].Alive = false
                    numAlive--
                }
                DisplayInfo(players[x].Name+" attacks and does", attp1, "points of damage with his", Weaps[players[x].Weap].Name, "to the ennemy.")
            }
        }
		if endBattle(players) || playerAction == 1 {
			break
		} else {
			DisplayInfo(players)
			round++
		}
	}
}

func DisplayInfof(format string, args ...interface{}) {
	fmt.Fprintf(Out, format, args...)
}

func DisplayInfo(args ...interface{}) {
	fmt.Fprintln(Out, args...)
}

func GetUserInput(i *int) {
	fmt.Fscan(In, i)
}

func endBattle(players []Character) bool {
	count := make([]int, 2)
	count[0] = 0
	count[1] = 0
	for _, pla := range players {
		if pla.Alive {
			if pla.Npc == false {
				count[0]++
			} else {
				count[1]++
			}
		}
	}
	if count[0] == 0 || count[1] == 0 {
		return true
	} else {
		return false
	}
}

func selectTarget(players []Character, x int) int {
	y := x
	for {
		y = y + 1
		if y >= len(players) {
			y = 0
		}
		if (players[y].Npc != players[x].Npc) && players[y].Alive {
			return y
		}
		if y == x {
			return -1
		}
	}
	return -1
}

func (g *Game) ProcessEvents(events []string) {
	for _, evtName := range events {
		g.Health += evts[evtName].ProcessEvent()
	}
}

func main() {
	g := &Game{Health: 100, Welcome: "Welcome to the Starship Enterprise\n\n", shield: "Minor Shield", weapon: "Minor Raygun"}
	g.Play()
}