package gui

import (
	"cellMachine/pkg/utils"
	"github.com/andlabs/ui"
	_ "github.com/andlabs/ui/winmanifest"
	"log"
	"math"
	"os"
)

var (
	Log     *log.Logger
	Warning *log.Logger
	Error   *log.Logger
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
	mainwin *ui.Window
	area    *ui.Area
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

	ui.OnShouldQuit(func() bool {
		core.mainwin.Destroy()
		return true
	})

	vbox := ui.NewVerticalBox()
	core.mainwin.SetChild(vbox)

	core.area = ui.NewArea(AreaHandler{})
	vbox.Append(core.area, true)
	Log.Println("UI is ready.")

	core.mainwin.Show()
}

func (core *Uicore) ShowWindow() {
	core.mainwin.Show()
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
			brush := ui.DrawBrush{
				R: cellComposer.BackColor.R,
				G: cellComposer.BackColor.G,
				B: cellComposer.BackColor.B,
				A: cellComposer.BackColor.A,
			}
			path := drawRect(Point{cellWidth * float64(i), cellHeight * float64(j)}, cellWidth, cellHeight)
			params.Context.Fill(path, &brush)
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

type AreaHandler struct{}

func (AreaHandler) Draw(a *ui.Area, p *ui.AreaDrawParams) {
	Log.Println("Draw call.")

	w := 40
	h := 30
	composer := utils.DefaultFieldComposer(w, h)
	handleComposer(composer, p)

	/*path := drawCircle(Point{300, 300}, 100)
	brush := ui.DrawBrush{A:1,R:1,G:1,B:1}
	p.Context.Fill(path, &brush)*/
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
