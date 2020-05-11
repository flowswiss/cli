package output

import (
	"fmt"
	"io"
	"reflect"
	"strings"
)

type column struct {
	index int
	name  string
	width int
}

func (c *column) format() string {
	return fmt.Sprintf("%%%ds", -c.width)
}

type table struct {
	columns []*column
	rows    [][]string
}

func (t *table) findColumn(name string) *column {
	for _, col := range t.columns {
		if col.name == name {
			return col
		}
	}
	return nil
}

func (t *table) insertColumns(cols []string) {
	for idx, col := range cols {
		t.columns = append(t.columns, &column{
			index: idx,
			name:  col,
			width: len(col),
		})
	}
}

func (t *table) insertRow(data map[string]interface{}) {
	row := make([]string, len(t.columns))

	for key, val := range data {
		col := t.findColumn(key)
		if col == nil {
			continue
		}

		str := fmt.Sprintf("%+v", val)
		row[col.index] = str

		if len(str) > col.width {
			col.width = len(str)
		}
	}

	t.rows = append(t.rows, row)
}

func (t *table) format(writer io.Writer, separator string, pretty bool) error {
	format := "%s"
	for idx, col := range t.columns {
		if pretty {
			format = col.format()
		}

		_, err := fmt.Fprintf(writer, format, strings.ToUpper(col.name))
		if err != nil {
			return err
		}

		if (idx + 1) < len(t.columns) {
			_, err := fmt.Fprintf(writer, separator)
			if err != nil {
				return err
			}
		}
	}

	_, err := fmt.Fprintln(writer)
	if err != nil {
		return err
	}

	for _, row := range t.rows {
		for idx, val := range row {
			if pretty {
				format = t.columns[idx].format()
			}

			if strings.Contains(val, separator) {
				val = fmt.Sprintf("%q", val)
			}

			_, err := fmt.Fprintf(writer, format, val)
			if err != nil {
				return err
			}

			if (idx + 1) < len(row) {
				_, err := fmt.Fprint(writer, separator)
				if err != nil {
					return err
				}
			}
		}

		_, err := fmt.Fprintln(writer)
		if err != nil {
			return err
		}
	}

	return nil
}

func (t *table) insertMap(value reflect.Value) error {
	if t.columns == nil {
		var cols []string
		for _, key := range value.MapKeys() {
			cols = append(cols, fmt.Sprintf("%v", key))
		}
		t.insertColumns(cols)
	}

	row := make(map[string]interface{})

	iter := value.MapRange()
	for iter.Next() {
		row[fmt.Sprintf("%v", iter.Key().Interface())] = iter.Value().Interface()
	}

	t.insertRow(row)
	return nil
}

func (t *table) insertStruct(value reflect.Value) error {
	if !value.Type().AssignableTo(reflect.TypeOf((*Displayable)(nil)).Elem()) {
		return fmt.Errorf("unable to serialize non `Displayable` struct of type %q", value.Type().String())
	}

	if t.columns == nil {
		columnsFunc := value.MethodByName("Columns")
		columns := columnsFunc.Call([]reflect.Value{})[0]

		var cols []string
		for i := 0; i < columns.Len(); i++ {
			cols = append(cols, columns.Index(i).String())
		}
		t.insertColumns(cols)
	}

	valuesFunc := value.MethodByName("Values")
	values := valuesFunc.Call([]reflect.Value{})[0]
	return t.insertMap(values)
}

func (t *table) insertValue(value reflect.Value) error {
	switch value.Kind() {
	case reflect.Array:
		fallthrough
	case reflect.Slice:
		for i := 0; i < value.Len(); i++ {
			err := t.insertValue(value.Index(i))
			if err != nil {
				return err
			}
		}
		return nil
	case reflect.Map:
		return t.insertMap(value)
	case reflect.Ptr:
		fallthrough
	case reflect.Struct:
		return t.insertStruct(value)
	}

	return fmt.Errorf("unable to serialize value of type %q (%q)", value.Type().String(), value.Kind().String())
}
