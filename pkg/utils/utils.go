package utils

import "math/rand"

// less or equal to 1.0
type Size float32

type Color struct {
	A, R, G, B float64
}

func DefaultColor() Color {
	return Color{1.0, 0.7, 1.0, 0.7}
}

// composer is an entry for storing visual representation of cells
// and translating them from simulation main code to ui code
type EntityComposer struct {
	Color Color
	Size  Size
}

func DefaultEntityComposer() EntityComposer {
	return EntityComposer{
		DefaultColor(),
		0.8,
	}
}

func EmptyEntityComposer() EntityComposer {
	return EntityComposer{
		Color: DefaultColor(),
		Size:  0,
	}
}

type CellComposer struct {
	BackColor Color
	Composer  EntityComposer
}

func DefaultCellComposer() CellComposer {
	return CellComposer{DefaultColor(), EmptyEntityComposer()}
}

type FieldComposer struct {
	Cells [][]CellComposer
	W, H  int
}

func DefaultFieldComposer(w, h int) FieldComposer {
	composer := FieldComposer{W: w, H: h}
	composer.Cells = make([][]CellComposer, w)
	for i := range composer.Cells {
		composer.Cells[i] = make([]CellComposer, h)
		for j := range composer.Cells[i] {
			composer.Cells[i][j] = DefaultCellComposer()
		}
	}
	return composer
}

func MutateFloat64(num, chance float64) float64 {
	dice := rand.Float64()
	if dice <= chance {
		factor := rand.Float64()/10 + 0.95 // from 0.95 to 1.05
		num *= factor
	}
	return num
}
