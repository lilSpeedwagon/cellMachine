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
	var posX, posY int
	for isPosEmpty := false; !isPosEmpty; {
		posX = rand.Intn(3) + x - 1
		posY = rand.Intn(3) + y - 1
		posX = (posX + field.W) % field.W // bounds handling
		posY = (posY + field.H) % field.H
		isPosEmpty = (posX != x || posY != y) && (field.Cells[posX][posY].entity == nil)
	}
	field.PutEntity(e, posX, posY)
}

func (field *CellField) MakeComposer() utils.FieldComposer {
	composer := utils.DefaultFieldComposer(field.W, field.H)

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
			// to delete
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
	c.entity = nil
}

func (c *Cell) Divide() {
	e := *c.entity
	c.Kill()
	c.field.Divide(e, c.x, c.y)
}

// getters
func (c *Cell) FoodStorage() float64 {
	return c.foodStorage
}
func (c *Cell) BadConditions() float64 {
	return c.badConditions
}
