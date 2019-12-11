package main

import (
	"cellMachine/pkg/gui"
	"cellMachine/pkg/sim"
	"cellMachine/pkg/utils"
	"github.com/andlabs/ui"
	"log"
	"os"
)

var (
	Log     *log.Logger
	Warning *log.Logger
	Error   *log.Logger
)

func initLog() {
	Log = log.New(os.Stdout,
		"LOG: ",
		log.Ldate|log.Ltime|log.Lshortfile)

	Warning = log.New(os.Stdout,
		"WARNING: ",
		log.Ldate|log.Ltime|log.Lshortfile)

	Error = log.New(os.Stdout,
		"ERROR: ",
		log.Ldate|log.Ltime|log.Lshortfile)
}

func main() {
	initLog()
	Log.Println("Application initialization...")

	closeApp := make(chan bool)
	composerChan := make(chan utils.FieldComposer)

	core := gui.Uicore{CloseApp: closeApp, ComposerChan: composerChan}
	go ui.Main(core.Init)

	w := 50
	h := 50

	var simulator sim.Simulator
	simulator.Init(w, h, composerChan)

	Log.Println("Ready")

	go simulator.Start()

	<-closeApp
	simulator.Stop()
	close(closeApp)
	close(composerChan)
	Log.Println("Closing application...")
}
