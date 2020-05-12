package commands

import "github.com/spf13/viper"

func init() {
	root.PersistentFlags().String("format", "table", "Format to output the data in. Allowed values are table, csv or json")
	handleError(viper.BindPFlag("format", root.PersistentFlags().Lookup("format")))
}

func display(val interface{}) error {
	switch config.Format {
	case "csv":
		return stdout.DisplayTable(val, ",", false)
	case "json":
		return stdout.DisplayJson(val)
	case "table":
		fallthrough
	default:
		return stdout.DisplayTable(val, "   ", true)
	}
}
