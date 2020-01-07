package main

import (
	"cellMachine/pkg/gui"
	"cellMachine/pkg/sim"
	"cellMachine/pkg/utils"
	"fmt"
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

func showInfo() {
	fmt.Println("<< Cell Machine >>")
	fmt.Println("Egor Sorokin, 2019")
	fmt.Println()
}

func main() {
	showInfo()
	initLog()
	Log.Println("Application initialization...")

	closeApp := make(chan bool)
	composerChan := make(chan utils.FieldComposer)
	readyChan := make(chan utils.Ready)

	core := gui.Uicore{CloseApp: closeApp, ComposerChan: composerChan, ReadyChan: readyChan}
	go ui.Main(core.Init)

	w := 50
	h := 50

	var simulator sim.Simulator
	simulator.Init(w, h, composerChan)

	// waiting for UI initialisation
	<-readyChan
	go simulator.Start()

	<-closeApp
	simulator.Stop()
	close(closeApp)
	close(composerChan)
	Log.Println("Closing application...")
}
