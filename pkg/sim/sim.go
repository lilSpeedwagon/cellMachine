package sim

import (
	"cellMachine/pkg/Cell"
	"cellMachine/pkg/utils"
	"log"
	"math/rand"
	"os"
	"time"
)

const (
	turnDelay  = time.Second / 20
	baseWidth  = 40
	baseHeight = 40
)

var (
	Log     *log.Logger
	Warning *log.Logger
	Error   *log.Logger
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

type SimulationInfo struct {
	turnCounter     uint64
	mutationCounter uint64
	entityCounter   uint64
}

func (info *SimulationInfo) Turns() uint64 {
	return info.turnCounter
}
func (info *SimulationInfo) Mutations() uint64 {
	return info.mutationCounter
}

func (info *SimulationInfo) Reset() {
	info.turnCounter = 0
	info.mutationCounter = 0
}

type Simulator struct {
	field     *Cell.CellField
	turnTimer time.Ticker
	ready     bool
	info      SimulationInfo

	composerChan chan<- utils.FieldComposer
}

func (sim *Simulator) Init(configPath string, composerChan chan utils.FieldComposer) {
	initLog()
	Log.Println("Simulation init")

	sim.composerChan = composerChan

	var err error
	sim.field, err = initFieldByJSON(configPath)
	if err != nil {
		Error.Printf(err.Error())
		panic(err.Error())
	}

	sim.sendAsync()

	sim.ready = true

	Log.Println("Ready.")
}

func (sim *Simulator) turn() {
	if sim.ready == false {
		return
	}
	sim.ready = false

	sim.info.turnCounter++

	sim.field.Update()
	sim.info.mutationCounter = Cell.MutationCounter
	sim.info.entityCounter = sim.field.EntityCount()

	sim.sendAsync()
	sim.ready = true
}

func (sim *Simulator) sendAsync() {
	composer := sim.field.MakeComposer()
	composer.Turns = sim.info.turnCounter
	composer.Mutations = sim.info.mutationCounter
	composer.Entities = sim.info.entityCounter
	select {
	case sim.composerChan <- composer:
	default:
	}
}

func (sim *Simulator) Start() {
	Log.Println("Starting simulation...")
	rand.Seed(time.Now().UnixNano())
	sim.info.Reset()
	sim.turnTimer = *time.NewTicker(turnDelay)
	go func() {
		for range sim.turnTimer.C {
			sim.turn()
		}
	}()
}

func (sim *Simulator) Stop() {
	Log.Println("Stopping simulation...")
	sim.turnTimer.Stop()
}
