package sim

import (
	"cellMachine/pkg/Cell"
	"encoding/json"
	"io/ioutil"
	"os"
)

type cellDrop struct {
	TypeName string
	X, Y, R  int
}

type entityDrop struct {
	TypeName string
	X, Y, R  int
}

type cellDropRect struct {
	TypeName   string
	X, Y, W, H int
}

type entityDropRect struct {
	TypeName   string
	X, Y, W, H int
}

type parsingStruct struct {
	CellTypes    []Cell.CellType
	EntityTypes  []Cell.EntityType
	Width        int
	Height       int
	BaseCellType string
	DropFood     bool
	CellDrops    []cellDrop
	EntityDrops  []entityDrop
	CellRects    []cellDropRect
	EntityRects  []entityDropRect
}

func parseJson(jsonBytes []byte) (*Cell.CellField, error) {

	// unmarshalling
	var unmarshalledObjects parsingStruct
	err := json.Unmarshal(jsonBytes, &unmarshalledObjects)
	if err != nil {
		Error.Printf("Marshalling error: %s", err.Error())
		return nil, err
	}

	// definition of cellTypes
	Cell.MinAntibiotic = 100000
	cellTypes := make(map[string]Cell.CellType, 0)
	for i := range unmarshalledObjects.CellTypes {
		t := unmarshalledObjects.CellTypes[i]
		cellTypes[t.Name] = Cell.CellType{Name: t.Name, Antibiotic: t.Antibiotic, FoodStorage: t.FoodStorage}
		if t.Antibiotic > Cell.MaxAntibiotic {
			Cell.MaxAntibiotic = t.Antibiotic
		}
		if t.Antibiotic < Cell.MinAntibiotic {
			Cell.MinAntibiotic = t.Antibiotic
		}
	}

	// definition of entityTypes
	entityTypes := make(map[string]Cell.EntityType, 0)
	for i := range unmarshalledObjects.EntityTypes {
		e := unmarshalledObjects.EntityTypes[i]
		entityTypes[e.Name] = Cell.EntityType{
			Name:            e.Name,
			ConsumptionBase: e.ConsumptionBase,
			Resistance:      e.Resistance,
			GrownRateBase:   e.GrownRateBase,
			MutationChance:  e.MutationChance,
		}
	}

	// definition a base type for whole field
	baseType := Cell.BaseCellType()
	if t, ok := cellTypes[unmarshalledObjects.BaseCellType]; ok {
		baseType = t
	} else {
		Warning.Printf("Base type %s not found", unmarshalledObjects.BaseCellType)
	}

	var dropFood bool = unmarshalledObjects.DropFood

	// field creation
	var field *Cell.CellField
	field = Cell.NewFieldWithBaseCell(unmarshalledObjects.Width, unmarshalledObjects.Height, baseType)
	field.DropFood(dropFood)

	// cell drops
	for i := range unmarshalledObjects.CellDrops {
		d := unmarshalledObjects.CellDrops[i]
		Log.Printf("Dropping cell of type %s in point %d : %d with radius %r", d.TypeName, d.X, d.Y, d.R)
		if t, ok := cellTypes[d.TypeName]; ok {
			err := field.DropCell(d.X, d.Y, d.R, t)
			if err != nil {
				Warning.Printf(err.Error())
			}
		} else {
			Warning.Printf("Type %s not found", d.TypeName)
		}
	}

	// entity drops
	for i := range unmarshalledObjects.EntityDrops {
		d := unmarshalledObjects.EntityDrops[i]
		Log.Printf("Dropping entity of type %s in point %d : %d with radius %r", d.TypeName, d.X, d.Y, d.R)
		if e, ok := entityTypes[d.TypeName]; ok {
			err := field.DropEntity(d.X, d.Y, d.R, e)
			if err != nil {
				Warning.Printf(err.Error())
			}
		} else {
			Warning.Printf("Type %s not found", d.TypeName)
		}
	}

	// cell rects
	for i := range unmarshalledObjects.CellRects {
		r := unmarshalledObjects.CellRects[i]
		Log.Printf("Dropping rectangle of cell of type %s in point %d : %d with size %d : %d", r.TypeName, r.X, r.Y, r.W, r.H)
		if c, ok := cellTypes[r.TypeName]; ok {
			err := field.DropCellRect(r.X, r.Y, r.W, r.H, c)
			if err != nil {
				Warning.Printf(err.Error())
			}
		} else {
			Warning.Printf("Type %s not found", r.TypeName)
		}
	}

	// entity rects
	for i := range unmarshalledObjects.EntityRects {
		r := unmarshalledObjects.EntityRects[i]
		Log.Printf("Dropping rectangle of entity of type %s in point %d : %d with size %d : %d", r.TypeName, r.X, r.Y, r.W, r.H)
		if c, ok := entityTypes[r.TypeName]; ok {
			err := field.DropEntityRect(r.X, r.Y, r.W, r.H, c)
			if err != nil {
				Warning.Printf(err.Error())
			}
		} else {
			Warning.Printf("Type %s not found", r.TypeName)
		}
	}

	return field, nil
}

func initFieldByJSON(fileName string) (*Cell.CellField, error) {
	Log.Printf("Opening file %s...", fileName)
	file, err := os.Open(fileName)
	if err != nil {
		Error.Printf("Cannot open file %s: %s", fileName, err.Error())
		return nil, err
	}
	defer func() {
		Warning.Printf("Closing file %s...", fileName)
		err := file.Close()
		if err != nil {
			Error.Printf(err.Error())
		}
	}()

	jsonBytes, err := ioutil.ReadAll(file)
	if err != nil {
		Error.Printf("Cannot read file %s: %s", fileName, err.Error())
		return nil, err
	}

	Log.Printf("Success. Parsing json...")
	return parseJson(jsonBytes)
}
