package commands

var format string

func init() {
	root.PersistentFlags().StringVar(&format, "format", "table", "Format to output the data in. Allowed values are table, csv or json")
}

func display(val interface{}) error {
	switch format {
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
