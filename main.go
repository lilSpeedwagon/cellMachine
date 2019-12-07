package main

import (
	"cellMachine/pkg/gui"
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
	composer := make(chan utils.FieldComposer)

	core := gui.Uicore{CloseApp: closeApp, Composer: composer}
	go ui.Main(core.Init)

	Log.Println("Ready")

	<-closeApp
	Log.Println("Closing application...")
}
