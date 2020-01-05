package Cell

import (
	"cellMachine/pkg/utils"
	"errors"
	"math"
	"math/rand"
)

const (
	baseFood       = 1000
	baseConditions = 5
)

var (
	MaxAntibiotic float64
	MinAntibiotic float64
)

// CellField

type CellField struct {
	Cells       [][]Cell
	W, H        int
	entityCount uint64
}

func (field *CellField) EntityCount() uint64 {
	return field.entityCount
}

func (field *CellField) Divide(e Entity, x, y int) {
	field.PutEntity(e, x, y)
	// make an array with free cells and iterate through them

	emptyCells := make([]utils.Position, 0)
	for i := x - 1; i <= x+1; i++ {
		for j := y - 1; j <= y+1; j++ {
			if i != x || j != y {
				posX := (i + field.W) % field.W
				posY := (j + field.H) % field.H
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
	field.Cells[x][y].entity = NewEntityFromEntity(e)
	field.Cells[x][y].entity.SetParent(&field.Cells[x][y])
	field.entityCount++
}

func (field *CellField) Update() {
	for i := 0; i < field.W; i++ {
		for j := 0; j < field.H; j++ {
			if field.Cells[i][j].entity != nil {
				field.Cells[i][j].entity.Update()
				if field.Cells[i][j].entity.IsReadyToDeath() {
					field.Cells[i][j].Kill()
					field.entityCount--
				} else if field.Cells[i][j].entity.IsReadyToDivide() {
					field.Cells[i][j].Divide()
					field.entityCount--
				}
				field.Cells[i][j].updateColor()
			}
		}
	}
}

func (field *CellField) drop(x, y, r int, operation func(int, int)) error {
	if x >= field.W || x < 0 || y >= field.H || y < 0 {
		return errors.New("invalid index")
	}

	for i := x - r; i <= x+r; i++ {
		for j := y - r; j <= y+r; j++ {
			dist := math.Sqrt(math.Pow(float64(x-i), 2) + math.Pow(float64(y-j), 2))
			if dist <= float64(r) {
				posX := (i + field.W) % field.W
				posY := (j + field.H) % field.H
				operation(posX, posY)
			}
		}
	}

	return nil
}

func (field *CellField) DropCell(x, y, r int, cellType CellType) error {
	return field.drop(x, y, r, func(posX, posY int) {
		field.Cells[posX][posY].badConditions = cellType.Antibiotic
		field.Cells[posX][posY].foodStorage = cellType.FoodStorage
		field.Cells[posX][posY].updateColor()
	})
}

func (field *CellField) DropEntity(x, y, r int, entityType EntityType) error {
	return field.drop(x, y, r, func(posX, posY int) {
		field.PutEntity(*NewEntityFromEntityType(entityType), posX, posY)
	})
}

func (field *CellField) dropRect(x, y, w, h int, operation func(int, int)) error {
	if x >= field.W || x < 0 || y >= field.H || y < 0 {
		return errors.New("invalid index")
	}

	x2 := x + w
	if x2 > field.W {
		x2 = field.W
	}
	y2 := y + h
	if y2 > field.H {
		y2 = field.H
	}

	for i := x; i < x2; i++ {
		for j := y; j < y2; j++ {
			operation(i, j)
		}
	}

	return nil
}

func (field *CellField) DropCellRect(x, y, w, h int, cellType CellType) error {
	return field.dropRect(x, y, w, h, func(posX, posY int) {
		field.Cells[posX][posY].badConditions = cellType.Antibiotic
		field.Cells[posX][posY].foodStorage = cellType.FoodStorage
		field.Cells[posX][posY].updateColor()
	})
}

func (field *CellField) DropEntityRect(x, y, w, h int, entityType EntityType) error {
	return field.dropRect(x, y, w, h, func(posX, posY int) {
		field.PutEntity(*NewEntityFromEntityType(entityType), posX, posY)
	})
}

func NewField(w, h int) *CellField {
	return NewFieldWithBaseCell(w, h, CellType{
		FoodStorage: baseFood,
		Antibiotic:  baseConditions,
	})
}

func NewFieldWithBaseCell(w, h int, base CellType) *CellField {
	field := new(CellField)
	field.W = w
	field.H = h
	field.Cells = make([][]Cell, w)
	for i := 0; i < w; i++ {
		field.Cells[i] = make([]Cell, h)
		for j := 0; j < h; j++ {
			field.Cells[i][j].field = field
			field.Cells[i][j].x = i
			field.Cells[i][j].y = j
			field.Cells[i][j].foodStorage = base.FoodStorage
			field.Cells[i][j].badConditions = base.Antibiotic
			field.Cells[i][j].updateColor()
		}
	}
	return field
}

// Cell

// for json unmarshalling
type CellType struct {
	Name        string
	FoodStorage float64
	Antibiotic  float64
}

func BaseCellType() CellType {
	return CellType{Name: "base", FoodStorage: baseFood, Antibiotic: baseConditions}
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
	if c.color.A > 0.5 {
		c.color.A = 0.5
	}
	c.color.R = (c.badConditions - MinAntibiotic) / (MaxAntibiotic - MinAntibiotic)
	c.color.G = 0.3
	c.color.B = 0.3
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

func (c *Cell) countNeighbors() int {
	count := 0
	for i := c.x - 1; i <= c.x+1; i++ {
		for j := c.y - 1; j <= c.y+1; j++ {
			if i != c.x || j != c.y {
				posX := (i + c.field.W) % c.field.W
				posY := (j + c.field.H) % c.field.H
				if c.field.Cells[posX][posY].entity != nil {
					count++
				}
			}
		}
	}
	return count
}
