package sim

import (
	"cellMachine/pkg/utils"
	"log"
	"os"
	"time"
)

var (
	Log     *log.Logger
	Warning *log.Logger
	Error   *log.Logger

	turnDelay = time.Second / 2
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
	field        *CellField
	ComposerChan chan<- utils.FieldComposer
	turnTimer    time.Ticker
	turnCounter  int
}

func (simulator *Simulator) Init(w, h int, composerChan chan utils.FieldComposer) {
	initLog()
	Log.Println("Simulation init")

	simulator.ComposerChan = composerChan
	simulator.field = NewField(w, h)

	color := utils.Color{1.0, 1.0, 0, 0}
	for i := 23; i <= 26; i++ {
		for j := 23; j <= 26; j++ {
			simulator.field.Cells[i][j].Entity = NewEntity(color, 0.8)
		}
	}

	simulator.ComposerChan <- simulator.field.makeComposer()

	Log.Println("Ready.")
}

func count(w, h, x, y int, cells [][]Cell) int {
	count := 0
	for i := x - 1; i <= x+1; i++ {
		for j := y - 1; j <= y+1; j++ {
			posX := (i + w) % w
			posY := (j + h) % h
			if cells[posX][posY].Entity != nil {
				count++
			}
		}
	}
	return count
}

func (sim *Simulator) turn() {
	sim.turnCounter++
	Log.Println("Turn ", sim.turnCounter)
	w := sim.field.W
	h := sim.field.H
	newField := NewField(w, h)

	color := utils.Color{1.0, 1.0, 0, 0}
	for i := 0; i < w; i++ {
		for j := 0; j < h; j++ {
			count := count(w, h, i, j, sim.field.Cells)
			if count == 2 || count == 3 {
				newField.Cells[i][j].Entity = NewEntity(color, 0.8)
			}
		}
	}

	sim.field = newField
	composer := sim.field.makeComposer()
	sim.ComposerChan <- composer
}

func (sim *Simulator) Start() {
	Log.Println("Starting simulation...")
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
