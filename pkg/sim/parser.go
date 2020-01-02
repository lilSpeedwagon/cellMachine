package sim

import (
	"cellMachine/pkg/Cell"
	"encoding/json"
	"io/ioutil"
	"os"
)

type cellDrop struct {
	name    string
	x, y, r int
}

type entityDrop struct {
	name    string
	x, y, r int
}

func parseJson(jsonBytes []byte) (*Cell.CellField, error) {
	// unmarshalling
	var unmarshalledObjects map[string]interface{}
	err := json.Unmarshal(jsonBytes, &unmarshalledObjects)
	if err != nil {
		Error.Printf("Marshalling error: %s", err.Error())
	}

	Log.Printf("%s", unmarshalledObjects)

	var field *Cell.CellField

	// base values
	w, h := baseWidth, baseHeight
	cellTypes := map[string]Cell.CellType{}
	entityTypes := map[string]Cell.EntityType{}

	// parsing
	if cellTypesRaw, ok := unmarshalledObjects["cellTypes"]; ok {
		cellTypes = cellTypesRaw.(map[string]Cell.CellType)
	} else {
		Warning.Printf("cellTypes not found")
	}

	if entityTypesRaw, ok := unmarshalledObjects["cellTypes"]; ok {
		entityTypes = entityTypesRaw.(map[string]Cell.EntityType)
	} else {
		Warning.Printf("entityTypes not found")
	}

	if wRaw, ok := unmarshalledObjects["width"]; ok {
		w = wRaw.(int)
	} else {
		Warning.Printf("width not found")
	}

	if hRaw, ok := unmarshalledObjects["width"]; ok {
		h = hRaw.(int)
	} else {
		Warning.Printf("height not found")
	}

	baseCellTypeName, hasBaseType := unmarshalledObjects["baseCellType"].(string)
	if baseCellType, ok := cellTypes[baseCellTypeName]; ok && hasBaseType {
		field = Cell.NewFieldWithBaseCell(w, h, baseCellType)
	} else {
		Warning.Printf("baseCellType not found")
		field = Cell.NewField(w, h)
	}

	cellDrops, hasCellDrops := unmarshalledObjects["cellDrops"].([]cellDrop)
	if hasCellDrops {
		for i := range cellDrops {
			drop := cellDrops[i]
			if cellType, ok := cellTypes[drop.name]; ok {
				dropErr := field.DropCell(drop.x, drop.y, drop.r, cellType)
				if dropErr != nil {
					Error.Printf(dropErr.Error())
					return nil, dropErr
				}
			}
		}
	}

	entityDrops, hasEntityDrops := unmarshalledObjects["entityDrops"].([]entityDrop)
	if hasEntityDrops {
		for i := range entityDrops {
			drop := entityDrops[i]
			if entityType, ok := entityTypes[drop.name]; ok {
				dropErr := field.DropEntity(drop.x, drop.y, drop.r, entityType)
				if dropErr != nil {
					Error.Printf(dropErr.Error())
					return nil, dropErr
				}
			}
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
		file.Close()
	}()

	jsonBytes, err := ioutil.ReadAll(file)
	if err != nil {
		Error.Printf("Cannot read file %s: %s", fileName, err.Error())
		return nil, err
	}

	Log.Printf("Success. Parsing json...")
	return parseJson(jsonBytes)
}
