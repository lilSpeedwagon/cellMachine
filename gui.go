package main

import (
	"github.com/andlabs/ui"
	_ "github.com/andlabs/ui/winmanifest"
)

type uicore struct {
	mainwin *ui.Window
}

func (core *uicore) init()	{
	core.mainwin = ui.NewWindow("cell machine", 800, 800, true)
	core.mainwin.SetMargined(true)

	core.mainwin.OnClosing(func(*ui.Window) bool {
		core.mainwin.Destroy()
		ui.Quit()
		return false
	})

	ui.OnShouldQuit(func() bool {
		core.mainwin.Destroy()
		return true
	})
}

func setupUI() {
	mainwin := ui.NewWindow("cell machine", 800, 800, true)
	mainwin.SetMargined(true)
	mainwin.OnClosing(func(*ui.Window) bool {
		mainwin.Destroy()
		ui.Quit()
		return false
	})
	ui.OnShouldQuit(func() bool {
		mainwin.Destroy()
		return true
	})

	grid := ui.NewGrid()
	hbox := ui.NewHorizontalBox()
	hbox.SetPadded(true)
	mainwin.SetChild(hbox)

	vbox := ui.NewVerticalBox()
	vbox.SetPadded(true)
	hbox.Append(vbox, false)

	area := ui.NewArea(areaHandler{})

	fontButton = ui.NewFontButton()
	fontButton.OnChanged(func(*ui.FontButton) {
		area.QueueRedrawAll()
	})
	vbox.Append(fontButton, false)

	form := ui.NewForm()
	form.SetPadded(true)
	// TODO on OS X if this is set to 1 then the window can't resize; does the form not have the concept of stretchy trailing space?
	vbox.Append(form, false)

	alignment = ui.NewCombobox()
	// note that the items match with the values of the uiDrawTextAlign values
	alignment.Append("Left")
	alignment.Append("Center")
	alignment.Append("Right")
	alignment.SetSelected(0)		// start with left alignment
	alignment.OnSelected(func(*ui.Combobox) {
		area.QueueRedrawAll()
	})
	form.Append("Alignment", alignment, false)

	hbox.Append(area, true)

	mainwin.Show()
}

