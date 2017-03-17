package main

import (
	"github.com/MartinSahlen/bq-utils/bqutils"
	"github.com/docopt/docopt-go"
)

func main() {
	usage := `BigQuery Utilities

Usage:
  bq-utils --project=<project> (--csv|--ndjson) --output=<file> (--query=<query>|--table=<table>)
  bq-utils --project=<project> --excel --output=<file> (--query=<query> <query-sheet-name>|--table=<table> <table-sheet-name>)...

Options:
  -h --help                     Show this screen
  -p project --project=project  The GCP project you are working with.
  -q query --query=query        The query to use as input to the csv writer
  -t table --table=table        The table to use as input to the csv writer
	-c --csv                      Use CSV as output for the writer
	-n --ndjson                   Use Newline delimited JSON as output for the writer
	-e --excel                    Use Excel as the output for the writer
	-o file --output=file         The path of the output file, i.e ~/Desktop/file.csv
	-l --legacy										Use legacy SQL for the queries
  -v --version                  Show version`

	arguments, err := docopt.Parse(usage, nil, true, "BigQuery Utilities 0.0 Pre-Alpha", false)

	if err != nil {
		panic(err)
	}

	err = run(arguments)

	if err != nil {
		panic(err)
	}
}

func run(arguments map[string]interface{}) error {
	csv := arguments["--csv"].(bool)
	excel := arguments["--excel"].(bool)
	ndjson := arguments["--ndjson"].(bool)
	//useLegacySQL := arguments["--legacysql"].(bool)

	filename := arguments["--output"].(string)
	project := arguments["--project"].(string)

	queries := arguments["--query"].([]string)
	querySheetNames := arguments["<query-sheet-name>"].([]string)
	tables := arguments["--table"].([]string)
	tableSheetNames := arguments["<table-sheet-name>"].([]string)

	//Docopt should make sure that the below combinations are the only legal ones.

	if csv && len(queries) == 1 {
		return bqutils.QueryToCsv(project, queries[0], filename)
	}

	if csv && len(tables) == 1 {
		return bqutils.TableToCsv(project, tables[0], filename)
	}

	if ndjson && len(queries) == 1 {
		return bqutils.QueryToNdJSON(project, queries[0], filename)
	}

	if ndjson && len(tables) == 1 {
		return bqutils.QueryToNdJSON(project, tables[0], filename)
	}

	if excel {
		writeConfigs := []bqutils.ExcelWriterConfig{}

		//We just be puttin' them queries first in the sheets
		for i, q := range queries {
			writeConfigs = append(writeConfigs, bqutils.ExcelWriterConfig{
				SheetName: querySheetNames[i],
				Query:     q,
				IsQuery:   true,
				Project:   project,
			})
		}

		for i, t := range tables {
			writeConfigs = append(writeConfigs, bqutils.ExcelWriterConfig{
				SheetName: tableSheetNames[i],
				Table:     t,
				IsQuery:   false,
				Project:   project,
			})
		}
		return bqutils.WriteToExcel(project, writeConfigs, filename)
	}
	return nil
}
