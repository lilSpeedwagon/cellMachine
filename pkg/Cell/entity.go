package Cell

import (
	"cellMachine/pkg/utils"
)

const (
	baseSize            = utils.Size(0.1)
	maxSize             = utils.Size(0.95)
	baseResistance      = 10
	baseConsumptionBase = 5
	baseMutationChance  = 0.05
	baseGrownRateBase   = 0.2
)

type Entity struct {
	// basic
	consumptionBase float64 // expected not more than 100
	resistance      float64 // expected not more than 100
	grownRateBase   float64 // less than 1.0
	mutationChance  float64
	// volatile
	color  utils.Color
	size   utils.Size
	parent *Cell
}

func (e *Entity) calculateColor() {
	e.color.A = 1.0
	e.color.R = e.consumptionBase / 100
	e.color.G = e.grownRateBase
	e.color.B = e.resistance / 100
}

func (e *Entity) Update() {
	vitality := (e.resistance - e.parent.BadConditions()) / e.resistance
	if vitality <= 0 {
		go e.die()
		return
	}

	grownRate := e.grownRateBase*vitality + 1
	consumptionVolume := grownRate * e.consumptionBase
	// isn't enough food in the cell
	if e.parent.Feed(consumptionVolume) < consumptionVolume {
		go e.die()
		return
	}

	e.size *= utils.Size(grownRate)
	if e.size >= maxSize {
		go e.divide()
	}
}

func (e *Entity) die() {
	if e.parent != nil {
		e.parent.Kill()
	}
}

func (e *Entity) divide() {
	e.parent.Divide()
}

// getters
func (e *Entity) Color() utils.Color {
	return e.color
}
func (e *Entity) Size() utils.Size {
	return e.size
}
func (e *Entity) Parent() *Cell {
	return e.parent
}

func (e *Entity) SetParent(c *Cell) {
	e.parent = c
}

func NewEntity() *Entity {
	entity := new(Entity)
	entity.size = baseSize
	entity.resistance = baseResistance
	entity.grownRateBase = baseGrownRateBase
	entity.consumptionBase = baseConsumptionBase
	entity.mutationChance = baseMutationChance
	entity.calculateColor()
	return entity
}

func NewEntityFrom(entity Entity) *Entity {
	e := new(Entity)
	mutationChance := entity.mutationChance
	e.size = baseSize
	e.mutationChance = mutationChance
	e.grownRateBase = utils.MutateFloat64(entity.grownRateBase, mutationChance)
	e.resistance = utils.MutateFloat64(entity.resistance, mutationChance)
	e.consumptionBase = utils.MutateFloat64(entity.consumptionBase, mutationChance)
	e.calculateColor()
	return e
}
