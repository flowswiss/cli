package commands

import (
	"encoding/json"
	"github.com/flowswiss/cli/pkg/output"
	"github.com/spf13/viper"
)

const flagFormat = "format"

func init() {
	root.PersistentFlags().String(flagFormat, "table", "Format to output the data in. Allowed values are table, csv or json")
	handleError(viper.BindPFlag("format", root.PersistentFlags().Lookup(flagFormat)))
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

	table.Format(stdout, separator, pretty)

	stderr.Printf("Found a total of %d items\n", len(table.Rows))
	return nil
}
