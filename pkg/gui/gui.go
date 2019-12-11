package gui

import (
	"cellMachine/pkg/utils"
	"github.com/andlabs/ui"
	_ "github.com/andlabs/ui/winmanifest"
	"log"
	"math"
	"os"
	"time"
)

var (
	Log         *log.Logger
	Warning     *log.Logger
	Error       *log.Logger
	redrawDelay = 100 * time.Millisecond
)

func initUILog() {
	Log = log.New(os.Stdout,
		"UI LOG: ",
		log.Ldate|log.Ltime|log.Lshortfile)

	Warning = log.New(os.Stdout,
		"UI WARNING: ",
		log.Ldate|log.Ltime|log.Lshortfile)

	Error = log.New(os.Stdout,
		"UI ERROR: ",
		log.Ldate|log.Ltime|log.Lshortfile)
}

func NewBrush(color utils.Color) ui.DrawBrush {
	var brush ui.DrawBrush
	brush.A = color.A
	brush.R = color.R
	brush.G = color.G
	brush.B = color.B
	return brush
}

type Point struct {
	x, y float64
}

var (
	strokeBrush = ui.DrawBrush{
		R: 0.0,
		G: 0.0,
		B: 0.0,
		A: 0.6,
	}
	strokeParams = ui.DrawStrokeParams{
		Thickness: 1.0,
	}
)

type Uicore struct {
	CloseApp     chan<- bool
	ComposerChan <-chan utils.FieldComposer
	mainwin      *ui.Window
	area         *ui.Area
	redrawTimer  time.Ticker
}

func (core *Uicore) Init() {
	initUILog()
	Log.Println("UI initialization...")

	core.mainwin = ui.NewWindow("cell machine", 800, 800, true)
	core.mainwin.SetMargined(true)
	core.mainwin.OnClosing(func(*ui.Window) bool {
		Log.Println("Closing...")
		core.mainwin.Destroy()
		ui.Quit()
		return false
	})
	vbox := ui.NewVerticalBox()
	core.mainwin.SetChild(vbox)
	core.mainwin.OnClosing(core.OnCloseWindow)
	core.mainwin.Show()

	ui.OnShouldQuit(func() bool {
		core.mainwin.Destroy()
		return true
	})

	areaHandler := AreaHandler{ComposerChannel: core.ComposerChan}
	core.area = ui.NewArea(&areaHandler)
	vbox.Append(core.area, true)

	core.redrawTimer = *time.NewTicker(redrawDelay)
	go func() {
		for range core.redrawTimer.C {
			core.area.QueueRedrawAll()
		}
	}()

	Log.Println("UI is ready.")
}

func (core *Uicore) ShowWindow() {
	core.mainwin.Show()
}

func (core *Uicore) OnCloseWindow(window *ui.Window) bool {
	Log.Println("Closing window...")
	core.CloseApp <- true
	return true
}

func drawLine(from, to Point) *ui.DrawPath {
	path := ui.DrawNewPath(ui.DrawFillModeWinding)
	path.NewFigure(from.x, from.y)
	path.LineTo(to.x, to.y)
	path.End()
	return path
}

func drawRect(from Point, w, h float64) *ui.DrawPath {
	path := ui.DrawNewPath(ui.DrawFillModeWinding)
	path.AddRectangle(from.x, from.y, w, h)
	path.End()
	return path
}

func drawCircle(center Point, radius float64) *ui.DrawPath {
	path := ui.DrawNewPath(ui.DrawFillModeWinding)
	path.NewFigureWithArc(center.x, center.y, radius, 0, 2*math.Pi, false)
	path.End()
	return path
}

func handleComposer(composer utils.FieldComposer, params *ui.AreaDrawParams) {
	cellWidth := params.AreaWidth / float64(composer.W)
	cellHeight := params.AreaHeight / float64(composer.H)

	for i := range composer.Cells {
		for j := range composer.Cells[i] {
			cellComposer := &composer.Cells[i][j]
			brush := NewBrush(cellComposer.BackColor)
			cellPath := drawRect(Point{cellWidth * float64(i), cellHeight * float64(j)}, cellWidth, cellHeight)
			params.Context.Fill(cellPath, &brush)
			if cellComposer.Composer.Size > 0 {
				center := Point{cellWidth * (float64(i) + 0.5), cellHeight * (float64(j) + 0.5)}
				radius := float64(cellComposer.Composer.Size) * cellWidth * 0.5
				if cellWidth > cellHeight {
					radius = float64(cellComposer.Composer.Size) * cellHeight * 0.5
				}
				entityPath := drawCircle(center, radius)
				brush = NewBrush(cellComposer.Composer.Color)
				params.Context.Fill(entityPath, &brush)
			}
			//params.Context.Stroke(path, &strokeBrush, &strokeParams)
		}
	}

	// lines
	for i := 1; i < composer.W; i++ {
		x := cellWidth * float64(i)
		path := drawLine(Point{x, 0}, Point{x, params.AreaHeight})
		params.Context.Stroke(path, &strokeBrush, &strokeParams)
	}
	for i := 1; i < composer.H; i++ {
		y := cellHeight * float64(i)
		path := drawLine(Point{0, y}, Point{params.AreaWidth, y})
		params.Context.Stroke(path, &strokeBrush, &strokeParams)
	}
}

type AreaHandler struct {
	ComposerChannel <-chan utils.FieldComposer
	Composer        utils.FieldComposer
}

func (handler *AreaHandler) Draw(a *ui.Area, p *ui.AreaDrawParams) {
	Log.Println("Draw call.")

	// non blocking composer receiving
	select {
	case handler.Composer = <-handler.ComposerChannel:
		Log.Println("New field composer received.")
	default:
	}

	if handler.Composer.Cells != nil {
		handleComposer(handler.Composer, p)
	}
}

func (AreaHandler) MouseEvent(a *ui.Area, me *ui.AreaMouseEvent) {
	// do nothing
}

func (AreaHandler) MouseCrossed(a *ui.Area, left bool) {
	// do nothing
}

func (AreaHandler) DragBroken(a *ui.Area) {
	// do nothing
}

func (AreaHandler) KeyEvent(a *ui.Area, ke *ui.AreaKeyEvent) (handled bool) {
	// reject all keys
	return false
}
