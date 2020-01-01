package sim

import (
	"cellMachine/pkg/Cell"
	"io/ioutil"
	"os"
)

/*
{
	"cellTypes":
	[
		{
			"name": "ctype1",
			"param1": "val1",
			"param2": "val2"
		}
	],
	"entityTypes":
	[
		{
			"name": "etype1",
			"param1": "val1",
			"param2": "val2"
		}
	],
	""
}
*/

type cellTypeName string
type entityTypeName string

type cellDrop struct {
	name    cellTypeName
	x, y, r int
}

type cellRect struct {
	name       cellTypeName
	x, y, w, h int
}

// for json unmarshalling
type initConditions struct {
	cellTypes    []Cell.CellType
	entityTypes  []Cell.EntityType
	w, h         int
	baseCellType cellTypeName // to fill whole field

}

func parseJson(json []byte, field *Cell.CellField) {

}

func initFieldByJSON(fileName string, field *Cell.CellField) error {
	Log.Printf("Opening file %s...", fileName)
	file, err := os.Open(fileName)
	if err != nil {
		Error.Printf("Cannot open file %s.", fileName)
		return err
	}
	defer func() {
		Warning.Printf("Closing file %s...", fileName)
		file.Close()
	}()

	json, err := ioutil.ReadAll(file)
	if err != nil {
		Error.Printf("Cannot read file %s.", fileName)
		return err
	}

	parseJson(json, field)

	return nil
}
