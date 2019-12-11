package sim

import "cellMachine/pkg/utils"

type CellField struct {
	Cells [][]Cell
	W, H  int
}

func (field *CellField) makeComposer() utils.FieldComposer {
	composer := utils.DefaultFieldComposer(field.W, field.H)

	for i := 0; i < field.W; i++ {
		for j := 0; j < field.H; j++ {
			entityComposer := utils.EmptyEntityComposer()
			if field.Cells[i][j].Entity != nil {
				entityComposer.Size = field.Cells[i][j].Entity.Size
				entityComposer.Color = field.Cells[i][j].Entity.Color
			}
			cellComposer := utils.CellComposer{
				BackColor: field.Cells[i][j].Color,
				Composer:  entityComposer,
			}
			composer.Cells[i][j] = cellComposer
		}
	}

	return composer
}

func NewField(w, h int) *CellField {
	field := new(CellField)
	field.W = w
	field.H = h
	field.Cells = make([][]Cell, w)
	for i := 0; i < w; i++ {
		field.Cells[i] = make([]Cell, h)
		for j := 0; j < h; j++ {
			field.Cells[i][j].Color = utils.DefaultColor()
		}
	}
	return field
}

type Cell struct {
	Entity *Entity
	Color  utils.Color
}

type Entity struct {
	Color utils.Color
	Size  utils.Size
}

func NewEntity(color utils.Color, size utils.Size) *Entity {
	entity := new(Entity)
	entity.Size = size
	entity.Color = color
	return entity
}
