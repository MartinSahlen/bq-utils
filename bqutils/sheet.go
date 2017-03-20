package bqutils

type SheetWriterConfig struct {
	IsQuery   bool
	Project   string
	Query     string
	Table     string
	SheetName string
}

type SheetConfig struct {
	RowData   *RowData
	SheetName string
}

//StitchSheetNames : Since we get two arrays for queries and their corresponding sheet names,
//We need to "stitch" them together. Docopt guarantees that this will not blast
//Because the slices will have the same length
func StitchSheetNames(queriesOrTables, sheetNames []string, project string, isQuery bool) []SheetWriterConfig {
	writeConfigs := []SheetWriterConfig{}

	for i, queryOrTable := range queriesOrTables {

		writeConfig := SheetWriterConfig{
			SheetName: sheetNames[i],
			Project:   project,
		}

		if isQuery {
			writeConfig.Query = queryOrTable
			writeConfig.IsQuery = true
		} else {
			writeConfig.Table = queryOrTable
			writeConfig.IsQuery = false
		}
		writeConfigs = append(writeConfigs, writeConfig)
	}
	return writeConfigs
}

func sheetConfigToWriterConfig(config []SheetWriterConfig) ([]SheetConfig, error) {
	sheets := []SheetConfig{}

	for _, c := range config {
		rowData, err := c.Execute()

		if err != nil {
			return nil, err
		}

		sheets = append(sheets, SheetConfig{
			SheetName: c.SheetName,
			RowData:   rowData,
		})
	}
	return sheets, nil
}
