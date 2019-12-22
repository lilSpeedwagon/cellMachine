package Cell

import (
	"cellMachine/pkg/utils"
	"math/rand"
)

const (
	baseFood = 1000
)

type CellField struct {
	Cells [][]Cell
	W, H  int
}

func (field *CellField) Divide(e Entity, x, y int) {
	field.PutEntity(e, x, y)
	// make an array with free cells and iterate through them

	emptyCells := make([]utils.Position, 0)
	for i := x - 1; i <= x+1; i++ {
		for j := y - 1; j <= y+1; j++ {
			if i != x || j != y {
				posX := (i + field.W) % field.W
				posY := (j + field.W) % field.H
				if field.Cells[posX][posY].entity == nil {
					emptyCells = append(emptyCells, utils.Position{posX, posY})
				}
			}
		}
	}
	emptyCount := len(emptyCells)
	if emptyCount != 0 {
		pos := rand.Intn(emptyCount)
		field.PutEntity(e, emptyCells[pos].X, emptyCells[pos].Y)
	}
}

func (field *CellField) MakeComposer() utils.FieldComposer {
	composer := utils.MakeFieldComposer(field.W, field.H)

	for i := 0; i < field.W; i++ {
		for j := 0; j < field.H; j++ {
			entityComposer := utils.EmptyEntityComposer()
			if field.Cells[i][j].entity != nil {
				entityComposer.Size = field.Cells[i][j].entity.Size()
				entityComposer.Color = field.Cells[i][j].entity.Color()
			}
			cellComposer := utils.CellComposer{
				BackColor: field.Cells[i][j].color,
				Composer:  entityComposer,
			}
			composer.Cells[i][j] = cellComposer
		}
	}

	return composer
}

func (field *CellField) PutEntity(e Entity, x, y int) {
	field.Cells[x][y].entity = NewEntityFrom(e)
	field.Cells[x][y].entity.SetParent(&field.Cells[x][y])
}

func (field *CellField) PutDefaultEntity(x, y int) {

}

func (field *CellField) Update() {
	for i := 0; i < field.W; i++ {
		for j := 0; j < field.H; j++ {
			if field.Cells[i][j].entity != nil {
				field.Cells[i][j].entity.Update()
				if field.Cells[i][j].entity.IsReadyToDeath() {
					field.Cells[i][j].Kill()
				} else if field.Cells[i][j].entity.IsReadyToDivide() {
					field.Cells[i][j].Divide()
				}
				field.Cells[i][j].updateColor()
			}
		}
	}
}

func NewField(w, h int) *CellField {
	field := new(CellField)
	field.W = w
	field.H = h
	field.Cells = make([][]Cell, w)
	for i := 0; i < w; i++ {
		field.Cells[i] = make([]Cell, h)
		for j := 0; j < h; j++ {
			field.Cells[i][j].color = utils.DefaultColor()
			field.Cells[i][j].field = field
			field.Cells[i][j].x = i
			field.Cells[i][j].y = j
			field.Cells[i][j].foodStorage = baseFood
		}
	}
	return field
}

type Cell struct {
	entity      *Entity
	field       *CellField
	color       utils.Color
	x, y        int
	foodStorage float64
	// to split in several
	badConditions float64
}

func (c *Cell) updateColor() {
	c.color.A = c.foodStorage / baseFood
}

func (c *Cell) Feed(foodVolume float64) float64 {
	if c.foodStorage-foodVolume < 0 {
		volume := c.foodStorage
		c.foodStorage = 0
		return volume
	} else {
		c.foodStorage -= foodVolume
		return foodVolume
	}
}

func (c *Cell) Kill() {
	if c.entity != nil {
		c.entity.parent = nil
		c.entity = nil
	}
}

func (c *Cell) Divide() {
	if c.entity != nil {
		e := *c.entity
		c.Kill()
		c.field.Divide(e, c.x, c.y)
	}
}

// getters
func (c *Cell) FoodStorage() float64 {
	return c.foodStorage
}
func (c *Cell) BadConditions() float64 {
	return c.badConditions
}
