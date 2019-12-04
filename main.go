package main

import (
	"cellMachine/pkg/gui"
	"log"
	"os"
)

var (
	Log *log.Logger
	Warning *log.Logger
	Error *log.Logger
)

func initLog()	{
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
	core := gui.Uicore{}
	core.Init()
	core.ShowWindow()
	Log.Println("Ready")
}