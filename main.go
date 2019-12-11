package main

import (
	"cellMachine/pkg/gui"
	"cellMachine/pkg/utils"
	"github.com/andlabs/ui"
	"log"
	"os"
	"time"
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

	Log.Println("Ready")

	composer := utils.DefaultFieldComposer(30, 30)
	composerChan <- composer

	f := true
	i := 0
	for f {
		cell := utils.DefaultCellComposer()
		cell.BackColor = utils.Color{
			A: 1.0,
			R: 1.0,
			G: 0,
			B: 0,
		}
		composer.Cells[1][i] = cell
		composerChan <- composer
		i++
		time.Sleep(time.Second / 2)
	}

	<-closeApp
	close(closeApp)
	close(composerChan)
	Log.Println("Closing application...")
}
