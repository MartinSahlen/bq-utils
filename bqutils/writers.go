package bqutils

func QueryToCsv(project, query, filename string) error {
	queryData, err := GetQueryData(project, query)

	if err != nil {
		return err
	}

	return WriteCsvFile(filename, queryData.Rows, queryData.Schema)
}

func TableToCsv(project, tablename, filename string) error {

	table, dataset, err := ParseTableName(tablename)

	if err != nil {
		return err
	}

	tableData, err := GetTableData(project, *dataset, *table)

	if err != nil {
		return err
	}

	return WriteCsvFile(filename, tableData.Rows, tableData.Schema)
}

func QueriesToExcel(project string, queries []string, filename string) error {

	return nil
}

func TablesToExcel(project string, tables []string, filename string) error {
	return nil
}
