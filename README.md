# bq-utils

**TLDR**: Export BigQuery tables and queries to csv or Excel sheet, supporting single excel file with multiple sheets. Supports basic CLI tasks + build your own using the code. The CLI is just composed of the functions contained in the `bqutils` package.

## Background and motivation
So, you know the feeling when the sales people ask you to "just grab some data"? Assuming your company is using **BigQuery** for you data and analytics needs, this is an easy task. Just launch the web UI and do the query. You can then download it as csv / ndjson file as well as exporting it to a Google sheet.

There are however, some limitations to this. Let's say, for instance that the query result is too big to save or download. Then the query must be saved (or you can use the temporary table) and exported to a GCS bucket, after which you must download and upload it to a sheet.... You get it.

What if I told you that you that there is a tool that allows you to do run a query / scan a table and dump it to your local file system. You can even run multiple queries / table scans and have them appear in different sheets in a n excel file. I present: `bq-utils`, a package for working with ad-hoc biguquery data in a friendlier way.

## Overview
`bq-utils` is a CLI **as well as** a set of functions that allow you to do stuff with bigquery such as

- Write a table or query to CSV
- Write a set of tables/queries to an excel spreadsheet with a sheet per query result
- Due to composability, most functions and features are exported. I give no guarantees to maintain compatibility, so please use a dependency manager to pin a version to make sure I don't break your builds! If this package matures I am sure there will be some conventions, semver and pull requests from the community that will stabilize the feature set.

## Install
`go get github.com/MartinSahlen/bq-utils`

## Usage

```bash
ls -la $YOLO
```
