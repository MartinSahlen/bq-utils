package main

import (
	"fmt"

	"github.com/docopt/docopt-go"
)

func main() {
	usage := `BigQuery Utilities.

Usage:
  bq-utils csv query <csv-query> <output-file>
  bq-utils csv table <csv-table> <output-file>
  bq-utils excel query <output-file> (<query> <sheetname>)...
  bq-utils excel table <output-file> (<table> <sheetname>)...
  bq-utils excel mixed <output-file> (<table>|<query> <sheetname>)...

Options:
  -h --help     Show this screen.
  --version     Show version.`

	arguments, _ := docopt.Parse(usage, nil, true, "BigQuery Utilities 0.0 Pre-Alpha", false)

	fmt.Println(arguments)
}
