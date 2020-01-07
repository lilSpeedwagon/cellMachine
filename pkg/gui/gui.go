package gui

import (
	"cellMachine/pkg/utils"
	"github.com/andlabs/ui"
	_ "github.com/andlabs/ui/winmanifest"
	"log"
	"math"
	"os"
	"strconv"
	"time"
)

const (
	strMutations = "Mutations: "
	strEntities  = "Entities: "
	strTurns     = "Turns: "
	fieldW       = 800
	fieldH       = 800
	infoH        = 100
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
	ReadyChan    chan<- utils.Ready
	redrawTimer  time.Ticker

	mainwin       *ui.Window
	area          *ui.Area
	turnLabel     *ui.Label
	entityLabel   *ui.Label
	mutationLabel *ui.Label
}

func (core *Uicore) Init() {
	initUILog()
	Log.Println("UI initialization...")

	core.mainwin = ui.NewWindow("cell machine", fieldW, fieldH+infoH, true)
	core.mainwin.SetMargined(true)
	core.mainwin.OnClosing(core.OnCloseWindow)

	infoBox := ui.NewHorizontalBox()
	core.turnLabel = ui.NewLabel(strTurns)
	infoBox.Append(core.turnLabel, true)
	core.mutationLabel = ui.NewLabel(strMutations)
	infoBox.Append(core.mutationLabel, true)
	core.entityLabel = ui.NewLabel(strEntities)
	infoBox.Append(core.entityLabel, true)

	areaHandler := areaHandler{composerChannel: core.ComposerChan, core: core}
	core.area = ui.NewArea(&areaHandler)

	gameBox := ui.NewVerticalBox()
	gameBox.Append(core.area, true)
	gameBox.Append(infoBox, false)
	core.mainwin.SetChild(gameBox)

	ui.OnShouldQuit(func() bool {
		core.mainwin.Destroy()
		return true
	})

	core.ReadyChan <- utils.Ready{}
	Log.Println("UI is ready.")

	core.redrawTimer = *time.NewTicker(redrawDelay)
	go func() {
		for range core.redrawTimer.C {
			core.area.QueueRedrawAll()
		}
	}()
	core.mainwin.Show()
}

func (core *Uicore) OnCloseWindow(window *ui.Window) bool {
	Log.Println("Closing window...")
	core.CloseApp <- true
	ui.Quit()
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
			cellPath.Free()

			if cellComposer.Composer.Size > 0 {
				center := Point{cellWidth * (float64(i) + 0.5), cellHeight * (float64(j) + 0.5)}
				radius := float64(cellComposer.Composer.Size) * cellWidth * 0.5
				if cellWidth > cellHeight {
					radius = float64(cellComposer.Composer.Size) * cellHeight * 0.5
				}

				entityPath := drawCircle(center, radius)
				brush = NewBrush(cellComposer.Composer.Color)
				params.Context.Fill(entityPath, &brush)
				entityPath.Free()
			}
		}
	}

	// lines
	for i := 1; i < composer.W; i++ {
		x := cellWidth * float64(i)
		path := drawLine(Point{x, 0}, Point{x, params.AreaHeight})
		params.Context.Stroke(path, &strokeBrush, &strokeParams)
		path.Free()
	}
	for i := 1; i < composer.H; i++ {
		y := cellHeight * float64(i)
		path := drawLine(Point{0, y}, Point{params.AreaWidth, y})
		params.Context.Stroke(path, &strokeBrush, &strokeParams)
		path.Free()
	}
}

type areaHandler struct {
	composerChannel <-chan utils.FieldComposer
	core            *Uicore
	contextStored   bool
}

func (handler *areaHandler) Draw(a *ui.Area, p *ui.AreaDrawParams) {
	composer := <-handler.composerChannel
	handler.core.turnLabel.SetText(strTurns + strconv.FormatUint(composer.Turns, 10))
	handler.core.mutationLabel.SetText(strMutations + strconv.FormatUint(composer.Mutations, 10))
	handler.core.entityLabel.SetText(strEntities + strconv.FormatUint(composer.Entities, 10))
	if composer.Cells != nil {
		handleComposer(composer, p)
	}
}

func (areaHandler) MouseEvent(a *ui.Area, me *ui.AreaMouseEvent) {
	// do nothing
}

func (areaHandler) MouseCrossed(a *ui.Area, left bool) {
	// do nothing
}

func (areaHandler) DragBroken(a *ui.Area) {
	// do nothing
}

func (areaHandler) KeyEvent(a *ui.Area, ke *ui.AreaKeyEvent) (handled bool) {
	// reject all keys
	return false
}
