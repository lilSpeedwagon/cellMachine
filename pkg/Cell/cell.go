package Cell

import (
	"cellMachine/pkg/utils"
	"errors"
	"math"
	"math/rand"
)

const (
	baseFood       = 10000
	baseFoodDelta  = 1
	baseConditions = 5
	maxCellAlpha   = 0.6
)

var (
	MaxAntibiotic float64
	MinAntibiotic float64
)

// CellField

type CellField struct {
	cells       [][]Cell
	newCells    [][]Cell
	W, H        int
	entityCount uint64
}

func (field *CellField) EntityCount() uint64 {
	return field.entityCount
}

func (field *CellField) Divide(e Entity, x, y int) {
	// make an array with free cells and iterate through them

	emptyCells := make([]utils.Position, 0)
	for i := x - 1; i <= x+1; i++ {
		for j := y - 1; j <= y+1; j++ {
			if i != x || j != y {
				posX := (i + field.W) % field.W
				posY := (j + field.H) % field.H
				if field.newCells[posX][posY].entity == nil {
					emptyCells = append(emptyCells, utils.Position{posX, posY})
				}
			}
		}
	}
	emptyCount := len(emptyCells)
	if emptyCount > 3 {
		pos := rand.Intn(emptyCount)
		field.PutEntity(e, emptyCells[pos].X, emptyCells[pos].Y)
		if emptyCount > 4 {
			field.PutEntity(e, x, y)
		}
	}
}

func (field *CellField) MakeComposer() utils.FieldComposer {
	composer := utils.MakeFieldComposer(field.W, field.H)

	for i := 0; i < field.W; i++ {
		for j := 0; j < field.H; j++ {
			entityComposer := utils.EmptyEntityComposer()
			if field.cells[i][j].entity != nil {
				entityComposer.Size = field.cells[i][j].entity.Size()
				entityComposer.Color = field.cells[i][j].entity.Color()
			}
			cellComposer := utils.CellComposer{
				BackColor: field.cells[i][j].color,
				Composer:  entityComposer,
			}
			composer.Cells[i][j] = cellComposer
		}
	}

	return composer
}

func (field *CellField) PutEntity(e Entity, x, y int) {
	field.newCells[x][y].entity = NewEntityFromEntity(e)
	field.newCells[x][y].entity.SetParent(&field.newCells[x][y])
	field.entityCount++
}

func (field *CellField) copyCellsToNew() {
	for i := 0; i < field.W; i++ {
		for j := 0; j < field.H; j++ {
			field.newCells[i][j] = field.cells[i][j]
			// make copy of entity to avoid pointers which point to same memory
			if field.cells[i][j].entity != nil {
				var e = new(Entity)
				e = field.cells[i][j].entity
				e.SetParent(&field.newCells[i][j])
				field.newCells[i][j].entity = e
			}
		}
	}
}

func (field *CellField) copyCellsFromNew() {
	for i := 0; i < field.W; i++ {
		for j := 0; j < field.H; j++ {
			field.cells[i][j] = field.newCells[i][j]
		}
	}
}

func (field *CellField) Update() {
	field.copyCellsToNew()
	for i := 0; i < field.W; i++ {
		for j := 0; j < field.H; j++ {
			cell := &field.newCells[i][j]

			if cell.foodStorage < cell.maxFood {
				field.newCells[i][j].foodStorage += baseFoodDelta
			}

			if cell.entity != nil {
				cell.entity.Update()
				if cell.entity.IsReadyToDeath() {
					cell.Kill()
					field.entityCount--
				} else if cell.entity.IsReadyToDivide() {
					cell.Divide()
					field.entityCount--
				}
			}

			cell.updateColor()
		}
	}
	field.copyCellsFromNew()
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
		field.cells[posX][posY].badConditions = cellType.Antibiotic
		field.cells[posX][posY].foodStorage = cellType.FoodStorage
		field.cells[posX][posY].maxFood = cellType.FoodStorage
		field.cells[posX][posY].updateColor()
	})
}

func (field *CellField) DropEntity(x, y, r int, entityType EntityType) error {
	return field.drop(x, y, r, func(posX, posY int) {
		field.cells[x][y].entity = NewEntityFromEntityType(entityType)
		field.cells[x][y].entity.SetParent(&field.cells[x][y])
		field.entityCount++
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
		field.cells[posX][posY].badConditions = cellType.Antibiotic
		field.cells[posX][posY].foodStorage = cellType.FoodStorage
		field.cells[posX][posY].maxFood = cellType.FoodStorage
		field.cells[posX][posY].updateColor()
	})
}

func (field *CellField) DropEntityRect(x, y, w, h int, entityType EntityType) error {
	return field.dropRect(x, y, w, h, func(posX, posY int) {
		field.cells[x][y].entity = NewEntityFromEntityType(entityType)
		field.cells[x][y].entity.SetParent(&field.cells[x][y])
		field.entityCount++
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
	field.cells = make([][]Cell, w)
	field.newCells = make([][]Cell, w)
	for i := 0; i < w; i++ {
		field.cells[i] = make([]Cell, h)
		field.newCells[i] = make([]Cell, h)
		for j := 0; j < h; j++ {
			field.cells[i][j].field = field
			field.cells[i][j].x = i
			field.cells[i][j].y = j
			field.cells[i][j].foodStorage = base.FoodStorage
			field.cells[i][j].maxFood = base.FoodStorage
			field.cells[i][j].badConditions = base.Antibiotic
			field.cells[i][j].updateColor()
		}
	}
	return field
}

// end of CellField

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
	maxFood     float64
	// to split in several
	badConditions float64
}

func (c *Cell) updateColor() {
	c.color.A = c.foodStorage / c.maxFood * maxCellAlpha
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

// end of Cell
