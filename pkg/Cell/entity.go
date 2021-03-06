package Cell

import (
	"cellMachine/pkg/utils"
	"math/rand"
)

const (
	baseSize            = utils.Size(0.1)
	maxSize             = utils.Size(0.95)
	baseResistance      = 10
	baseConsumptionBase = 5
	baseMutationChance  = 0.01
	baseGrownRateBase   = 0.2

	// for color calculation
	borderResistance  = 15
	borderConsumption = 10
	borderGrownRate   = 1.0
)

var MutationCounter uint64

type Mutator struct {
	mutationChance float64
}

func newMutator() Mutator {
	return Mutator{mutationChance: baseMutationChance}
}

func (m *Mutator) MutateFloat64(num float64) float64 {
	dice := rand.Float64()
	if dice <= m.mutationChance {
		factor := rand.Float64()/10.0 + 0.95 // from 0.95 to 1.05
		num *= factor
		MutationCounter++
	}
	return num
}

type EntityState struct {
	isReadyToDivide bool
	isReadyToDeath  bool
}

// for json unmarshalling
type EntityType struct {
	Name            string
	ConsumptionBase float64
	Resistance      float64
	GrownRateBase   float64
	MutationChance  float64
}

type Entity struct {
	// basic
	consumptionBase float64 // expected not more than 100
	resistance      float64 // expected not more than 100
	grownRateBase   float64 // less than 1.0
	mutator         Mutator
	// volatile
	color  utils.Color
	size   utils.Size
	parent *Cell
	state  EntityState
}

func (e *Entity) calculateColor() {
	e.color.A = 0.8
	e.color.R = e.resistance / borderResistance
	e.color.G = e.grownRateBase / borderGrownRate
	e.color.B = e.consumptionBase / borderConsumption
}

func (e *Entity) Update() {
	vitality := (e.resistance - e.parent.BadConditions()) / e.resistance
	if vitality <= 0 {
		e.state.isReadyToDeath = true
		return
	}

	grownRate := e.grownRateBase*vitality + 1
	consumptionVolume := grownRate * e.consumptionBase
	// isn't enough food in the cell
	if e.parent.Feed(consumptionVolume) < consumptionVolume {
		e.state.isReadyToDeath = true
		return
	}

	e.size *= utils.Size(grownRate)
	if e.size >= maxSize {
		e.state.isReadyToDivide = true
		return
	}

	e.calculateColor()
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
func (e *Entity) IsReadyToDivide() bool {
	return e.state.isReadyToDivide
}
func (e *Entity) IsReadyToDeath() bool {
	return e.state.isReadyToDeath
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
	entity.mutator = newMutator()
	entity.calculateColor()
	entity.state = EntityState{false, false}
	return entity
}

func NewEntityFromEntity(entity Entity) *Entity {
	e := new(Entity)
	e.mutator = entity.mutator
	e.size = baseSize
	e.grownRateBase = e.mutator.MutateFloat64(entity.grownRateBase)
	e.resistance = e.mutator.MutateFloat64(entity.resistance)
	e.consumptionBase = e.mutator.MutateFloat64(entity.consumptionBase)
	e.calculateColor()
	e.state = EntityState{false, false}
	return e
}

func NewEntityFromEntityType(base EntityType) *Entity {
	e := new(Entity)
	e.mutator = Mutator{mutationChance: base.MutationChance}
	e.size = baseSize
	e.grownRateBase = base.GrownRateBase
	e.resistance = base.Resistance
	e.consumptionBase = base.ConsumptionBase
	e.calculateColor()
	e.state = EntityState{false, false}
	return e
}
