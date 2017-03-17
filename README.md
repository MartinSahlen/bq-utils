# bq-utils
**TLDR**: Export BigQuery tables and queries to csv or Excel sheet, supporting excel files with multiple sheets. Supports some basic CLI tasks + build your own "whatever" using the code. The CLI is just composed of the functions and components contained in the `bqutils` package.

## Install
`go get github.com/MartinSahlen/bq-utils`

## Background and motivation
So, you know the feeling when the sales people ask you to "just grab some data"? Assuming your company is using **BigQuery** for yo' data and analytics needs, this is an easy task. Just launch the web UI and do the query. You can then download it as csv / ndjson file as well as exporting it to a Google sheet.

There are however, some limitations to this. Let's say, for instance that the query result is too big to save or download. Then the query must be saved (or you can use the temporary table) and exported to a GCS bucket, after which you must download and upload it to a sheet.... You get it.

What if I told you that you that there is a tool that allows you to do run a query / scan a table and dump it to your local file system. You can even run multiple queries / table scans and have them appear in different sheets in an excel file. I present: `bq-utils`, a package for working with ad-hoc biguquery data in a friendlier way.

## Overview
`bq-utils` is a CLI **as well as** a set of functions that allow you to do stuff with bigquery such as

- Write a table or query to CSV
- Write a set of tables/queries to an excel spreadsheet with a sheet per query result
- Due to composability, most functions and features are exported. I give no guarantees to maintain compatibility, so please use a dependency manager to pin a version to make sure I don't break your builds! If this package matures I am sure there will be some conventions, semver and pull requests from the community that will stabilize the feature set.
- Using [dep](https:///www.github.com/golang/dep) for dependency management, and ignoring vendor folder.
- Using STANDARD SQL!
- It will BLAST (probably) if the provided project is not the one you are currently authenticated against.

## Usage
Since we are using [docopt](https://github.com/docopt/docopt.go), I'm just pasting the Usage doc for that.

```
BigQuery Utilities

Usage:
bq-utils --project=<project> --csv --output=<file> (--query=<query>|--table=<table>)
bq-utils --project=<project> --excel --output=<file> (--query=<query> <query-sheet-name>|--table=<table> <table-sheet-name>)...

Options:
-h --help                     Show this screen
-p project --project=project  The GCP project you are working with.
-q query --query=query        The query to use as input to the csv writer
-t table --table=table        The table to use as input to the csv writer
-c --csv                      Use CSV as output for the writer
-e --excel                    Use Excel as the output for the writer
-o file --output=file         The path of the output file, i.e ~/Desktop/file.csv
-v --version                  Show version
```

# NB!
In the examples, I have not included the backticks around bigquery table names because
they ruined the markdown formatting.

#### Exporting a table to CSV
`bq-utils -p my-project -o file.csv -c -t dataset.table`

#### Exporting a query to CSV
`bq-utils -p my-project -o file.csv -c -q 'SELECT * FROM dataset.table'`

#### Exporting a complex query to CSV
In some cases, we might have complex queries that would be bad to write out in the CLI. Currently I haven't bothered to support anything here, it can be solved pretty easily like this:

`bq-utils -p my-project -o file.xlsx -c -q "$(cat query.sql)"`

#### Exporting a table to Excel
`bq-utils -p my-project -o file.xlsx -e -t dataset.table`

#### Exporting a query to Excel
`bq-utils -p my-project -o file.xslx -e -q 'SELECT * FROM dataset.table'`

#### Exporting a complex query to Excel
`bq-utils -p my-project -o file.csv -e -q "$(cat query.sql)"`

### Exporting a mix of queries and tables to Excel
`bq-utils -p my-project -o file.csv -e -q "$(cat query.sql)" complex-query-sheet -q 'SELECT * FROM dataset.table' query-sheet -t dataset.table table-1-sheet`

## Roadmap
- Loading files from file system / GCS (maybe already supported in some bqutil tool?)
- Exporting to Google sheets
- Workerpool for running parallel queries and table scans, and populating excel sheet.
- Support legacy SQL
- Support auto-detecting if you send a file for query parameter
