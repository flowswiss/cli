package output

import (
	"encoding/json"
	"reflect"
)

func (o *Output) DisplayJson(object interface{}) error {
	return json.NewEncoder(o.Writer).Encode(object)
}

func (o *Output) DisplayTable(object interface{}, separator string, pretty bool) error {
	tbl := &Table{}

	err := tbl.insertValue(reflect.ValueOf(object))
	if err != nil {
		return err
	}

	tbl.Format(o, separator, pretty)
	return nil
}
