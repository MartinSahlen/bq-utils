package main

import (
	"os"

	"github.com/MartinSahlen/bq-utils/bqutils"
	"github.com/docopt/docopt-go"
)

func main() {
	usage := `BigQuery Utilities.

Usage:
  bq-utils csv query <csv-query> <output-file>
  bq-utils csv table <csv-table> <output-file>
  bq-utils excel query <output-file> (<excel-query> <sheetname>)...
  bq-utils excel table <output-file> (<excel-table> <sheetname>)...
  bq-utils excel mixed <output-file> (q <excel-query> <query-sheetname>)|(t <excel-table> <table-sheetname>)...

Options:
  -h --help     Show this screen.
  --version     Show version.`

	project := os.Getenv("PROJECT")

	if project == "" {
		panic("PROJECT env var is not set")
	}

	arguments, _ := docopt.Parse(usage, nil, true, "BigQuery Utilities 0.0 Pre-Alpha", false)

	csv := arguments["csv"].(bool)
	excel := arguments["excel"].(bool)
	query := arguments["query"].(bool)
	table := arguments["table"].(bool)
	mixed := arguments["mixed"].(bool)
	filename := arguments["<output-file>"].(string)

	if csv {
		if query {
			q := arguments["<csv-query>"].(string)
			err := bqutils.QueryToCsv(project, q, filename)
			if err != nil {
				panic(err)
			}
		} else if table {
			t := arguments["<csv-table>"].(string)
			err := bqutils.TableToCsv(project, t, filename)
			if err != nil {
				panic(err)
			}
		}
	} else if excel {
		if query {
			sheets := arguments["<sheetname>"].([]string)
			es := []bqutils.ExcelWriterConfig{}
			for i, q := range arguments["<excel-query>"].([]string) {
				es = append(es, bqutils.ExcelWriterConfig{
					SheetName: sheets[i],
					Query:     &q,
					IsQuery:   true,
					Project:   project,
				})
			}
			err := bqutils.WriteToExcel(project, es, filename)
			if err != nil {
				panic(err)
			}
		} else if table {

		} else if mixed {

		}
	}
}
