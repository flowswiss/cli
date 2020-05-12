package commands

import (
	"encoding/json"
	"github.com/flowswiss/cli/pkg/output"
	"github.com/spf13/viper"
)

func init() {
	root.PersistentFlags().String("format", "table", "Format to output the data in. Allowed values are table, csv or json")
	handleError(viper.BindPFlag("format", root.PersistentFlags().Lookup("format")))
}

func display(val interface{}) error {
	if config.Format == "json" {
		return json.NewEncoder(stdout.Writer).Encode(val)
	}

	separator := "   "
	pretty := true

	if config.Format == "csv" {
		separator = ","
		pretty = false
	}

	table := output.Table{}

	err := table.Insert(val)
	if err != nil {
		return err
	}

	err = table.Format(stdout.Writer, separator, pretty)
	if err != nil {
		return err
	}

	stderr.Printf("Found a total of %d items\n", len(table.Rows))
	return nil
}
