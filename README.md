# bq-utils
**TLDR**: Export BigQuery tables and queries to csv, ndjson or Excel sheet, supporting excel files with multiple sheets. Supports some basic CLI tasks + build your own "whatever" using the code. The CLI is just composed of the functions and components contained in the `bqutils` package.

## Install
`go get github.com/MartinSahlen/bq-utils`

## Background and motivation
So, you know the feeling when the sales people ask you to "just grab some data"? Assuming your company is using **BigQuery** for yo' data and analytics needs, this is an easy task. Just launch the web UI and do the query. You can then download it as csv / ndjson file as well as exporting it to a Google sheet. Or you can use the [bq command line tool from](https://cloud.google.com/bigquery/bq-command-line-tool)

There are however, some limitations to this. Let's say, for instance that the query result is too big to save or download. Then the query must be saved (or you can use the temporary table) and exported to a GCS bucket, after which you must download and upload it to a sheet.... You get it. Also, What if somebody wants an excel file with multiple sheets?

What if I told you that you that there is a tool that allows you to do run a query / scan a table and dump it to your local file system. You can even run multiple queries / table scans and have them appear in different sheets in an excel file. I present: `bq-utils`, a package for working with ad-hoc biguquery data in a friendlier way.

*DISCLAIMER*: Of course there are not always people nagging you for data, many times I have simple wanted to grab a set of
data to play around with in a frontend or some other analytics tool like a graph database.

## Overview
`bq-utils` is a CLI **as well as** a set of functions that allow you to do stuff with bigquery such as

- Write a table or query to CSV
- Write a set of tables/queries to an excel spreadsheet with a sheet per query result
- Due to composability, most functions and features are exported. I give no guarantees to maintain compatibility, so please use a dependency manager to pin a version to make sure I don't break your builds! If this package matures I am sure there will be some conventions, semver and pull requests from the community that will stabilize the feature set.
- Using [dep](https:///www.github.com/golang/dep) for dependency management, and ignoring vendor folder.
- Using STANDARD SQL only
- It will BLAST (probably) if the provided project is not the one you are currently authenticated against.
- Everything (all data) runs through this computer so you might want to have a strong connection / run this in a box within google cloud.

## Usage

#### Exporting a table to CSV
`bq-utils -p my-project -o file.csv -c -t dataset.table`

#### Exporting a query to CSV
`bq-utils -p my-project -o file.csv -c -q 'SELECT * FROM dataset.table'`

#### Exporting a complex query to CSV
In some cases, we might have complex queries that would be bad to write out in the CLI. Currently I haven't bothered to support anything here, it can be solved pretty easily like this:

`bq-utils -p my-project -o file.csv -c -q "$(cat query.sql)"`

#### Exporting a table to NDJSON
`bq-utils -p my-project -o file.ndjson -n -t dataset.table`

#### Exporting a query to NDJSON
`bq-utils -p my-project -o file.ndjson -n -q 'SELECT * FROM dataset.table'`

#### Exporting a complex query to NDJSON
`bq-utils -p my-project -o file.ndjson -n -q "$(cat query.sql)"`

#### Exporting a table to Excel
`bq-utils -p my-project -o file.xslx -e -t dataset.table`

#### Exporting a query to Excel
`bq-utils -p my-project -o file.xslx -e -q 'SELECT * FROM dataset.table'`

#### Exporting a complex query to Excel
`bq-utils -p my-project -o file.csv -e -q "$(cat query.sql)"`

### Exporting a mix of queries and tables to Excel

When exporting to Excel, remember to include the sheet name
(even though you only have one sheet)

`bq-utils -p my-project -o file.csv -e -q "$(cat query.sql)" complex-query-sheet -q 'SELECT * FROM dataset.table' query-sheet -t dataset.table table-1-sheet`

## Roadmap
- Loading files from file system / GCS (maybe already supported in some bqutil tool?)
- Uploading files to GCS
- Exporting to Google sheets
- Workerpool for running parallel queries and table scans, and populating excel sheets in parallel.
- Support legacy SQL
- Support auto-detecting if you send a file for query parameter
- Add sending spreadsheets to a specified email if your boss wants some data TOMORROW
- Wrap errors to get origin of error
- Write tests
- DOCUMENTATION AND COMMENTS PLEASE
