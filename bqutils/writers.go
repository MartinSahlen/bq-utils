package bqutils

func QueryToCsv(project, query, filename string) error {
	queryData, err := GetQueryData(project, query)

	if err != nil {
		return err
	}

	return WriteCsvFile(filename, queryData.Rows, queryData.Schema)
}

func TableToCsv(project, tablename, filename string) error {

	dataset, table, err := ParseTableName(tablename)

	if err != nil {
		return err
	}

	tableData, err := GetTableData(project, *dataset, *table)

	if err != nil {
		return err
	}

	return WriteCsvFile(filename, tableData.Rows, tableData.Schema)
}

type ExcelWriterConfig struct {
	IsQuery   bool
	Project   string
	Query     string
	Table     string
	SheetName string
}

func (e ExcelWriterConfig) Exeute() (*RowData, error) {
	if e.IsQuery {
		return GetQueryData(e.Project, e.Query)
	}

	dataset, table, err := ParseTableName(e.Table)

	if err != nil {
		return nil, err
	}

	return GetTableData(e.Project, *dataset, *table)
}

func WriteToExcel(project string, ss []ExcelWriterConfig, filename string) error {

	sheets := []SheetConfig{}

	for _, s := range ss {
		rowData, err := s.Exeute()

		if err != nil {
			return err
		}

		sheets = append(sheets, SheetConfig{
			SheetName: s.SheetName,
			Schema:    rowData.Schema,
			Rows:      rowData.Rows,
		})
	}
	return WriteExcelFile(filename, sheets)
}
