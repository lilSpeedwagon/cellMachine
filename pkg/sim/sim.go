package sim

import (
	"cellMachine/pkg/Cell"
	"cellMachine/pkg/utils"
	"log"
	"math/rand"
	"os"
	"time"
)

var (
	Log     *log.Logger
	Warning *log.Logger
	Error   *log.Logger

	turnDelay = time.Second / 10
)

func initLog() {
	Log = log.New(os.Stdout,
		"SIMLOG: ",
		log.Ldate|log.Ltime|log.Lshortfile)

	Warning = log.New(os.Stdout,
		"SIMWARNING: ",
		log.Ldate|log.Ltime|log.Lshortfile)

	Error = log.New(os.Stdout,
		"SIMERROR: ",
		log.Ldate|log.Ltime|log.Lshortfile)
}

type Simulator struct {
	field        *Cell.CellField
	ComposerChan chan<- utils.FieldComposer
	turnTimer    time.Ticker
	turnCounter  int
}

func (simulator *Simulator) Init(w, h int, composerChan chan utils.FieldComposer) {
	initLog()
	Log.Println("Simulation init")

	simulator.ComposerChan = composerChan
	simulator.field = Cell.NewField(w, h)

	e := *Cell.NewEntity()

	for i := 23; i <= 26; i++ {
		for j := 23; j <= 26; j++ {
			simulator.field.PutEntity(e, i, j)
		}
	}

	simulator.ComposerChan <- simulator.field.MakeComposer()

	Log.Println("Ready.")
}

func (sim *Simulator) turn() {
	sim.turnCounter++
	Log.Println("Turn ", sim.turnCounter)
	sim.field.Update()
	composer := sim.field.MakeComposer()
	sim.ComposerChan <- composer
}

func (sim *Simulator) Start() {
	Log.Println("Starting simulation...")
	rand.Seed(time.Now().UnixNano())
	sim.turnCounter = 0
	sim.turnTimer = *time.NewTicker(turnDelay)
	go func() {
		for range sim.turnTimer.C {
			sim.turn()
		}
	}()
}

func (sim *Simulator) Stop() {
	Log.Println("Stopping simulation...")
}
