package output

import (
	"encoding/json"
	"reflect"
)

type Displayable interface {
	Columns() []string
	Values() map[string]interface{}
}

func (o *Output) DisplayJson(object interface{}) error {
	return json.NewEncoder(o.Writer).Encode(object)
}

func (o *Output) DisplayTable(object interface{}, separator string, pretty bool) error {
	tbl := &table{}

	err := tbl.insertValue(reflect.ValueOf(object))
	if err != nil {
		return err
	}

	return tbl.format(o.Writer, separator, pretty)
}
