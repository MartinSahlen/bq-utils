package bqutils

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/user"
	"path/filepath"

	"cloud.google.com/go/bigquery"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/sheets/v4"
)

const maxCells = 2000000

func writeGoogleSheetRow(row []string) *sheets.RowData {
	rowData := &sheets.RowData{
		Values: []*sheets.CellData{},
	}
	for _, cell := range row {
		rowData.Values = append(rowData.Values, &sheets.CellData{
			UserEnteredValue: &sheets.ExtendedValue{
				StringValue: cell,
			},
		})
	}
	return rowData
}

func calculateTotalCells(configs []SheetConfig) (int64, error) {
	var totalCells int64
	for _, config := range configs {
		totalCells += int64(config.RowData.NumRows) * int64(len(config.RowData.Schema))
	}
	if totalCells > maxCells {
		return totalCells, errors.New(fmt.Sprintf("The total number of cells (%d) in the ouptput will exceed the max limit (%d)", totalCells, maxCells))
	}
	return totalCells, nil
}

func WriteToGoogleSheet(config []SheetWriterConfig, name string) error {
	client, err := sheetClient()

	if err != nil {
		return err
	}

	sheetConfigs, err := sheetConfigToWriterConfig(config)

	if err != nil {
		return err
	}

	_, err = calculateTotalCells(sheetConfigs)

	if err != nil {
		return err
	}

	outputSheet := &sheets.Spreadsheet{
		Properties: &sheets.SpreadsheetProperties{
			Title: name,
		},
		Sheets: []*sheets.Sheet{},
	}

	for _, c := range sheetConfigs {
		sheet := &sheets.Sheet{
			Properties: &sheets.SheetProperties{
				Title: c.SheetName,
				GridProperties: &sheets.GridProperties{
					RowCount:    int64(c.RowData.NumRows) + 1, //+1 for the header row
					ColumnCount: int64(len(c.RowData.Schema)),
				},
			},
			Data: []*sheets.GridData{
				&sheets.GridData{
					RowData: []*sheets.RowData{},
				},
			},
		}

		header := []string{}

		for _, f := range c.RowData.Schema {
			header = append(header, f.Name)
		}

		sheet.Data[0].RowData = append(sheet.Data[0].RowData, writeGoogleSheetRow(header))

		mapper := func(row map[string]bigquery.Value, schema *bigquery.Schema) (map[string]bigquery.Value, error) {
			if schema == nil {
				return nil, errors.New("Schema is nil")
			}
			sheet.Data[0].RowData = append(sheet.Data[0].RowData, writeGoogleSheetRow(mapToStringSlice(row, *schema)))
			return nil, nil
		}

		err = MapRows(c.RowData.Rows, &c.RowData.Schema, mapper)
		if err != nil {
			return err
		}
		outputSheet.Sheets = append(outputSheet.Sheets, sheet)
	}
	_, err = client.Spreadsheets.Create(outputSheet).Do()
	return err
}

func sheetClient() (*sheets.Service, error) {
	ctx := context.Background()

	b, err := ioutil.ReadFile(os.Getenv("CLIENT_SECRET"))
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}

	config, err := google.ConfigFromJSON(b, sheets.SpreadsheetsScope)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}
	client := getClient(ctx, config)
	return sheets.New(client)
}

// getClient uses a Context and Config to retrieve a Token
// then generate a Client. It returns the generated Client.
func getClient(ctx context.Context, config *oauth2.Config) *http.Client {
	cacheFile, err := tokenCacheFile()
	if err != nil {
		log.Fatalf("Unable to get path to cached credential file. %v", err)
	}
	tok, err := tokenFromFile(cacheFile)
	if err != nil {
		tok = getTokenFromWeb(config)
		saveToken(cacheFile, tok)
	}
	return config.Client(ctx, tok)
}

// getTokenFromWeb uses Config to request a Token.
// It returns the retrieved Token.
func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser then type the "+
		"authorization code: \n%v\n", authURL)

	var code string
	if _, err := fmt.Scan(&code); err != nil {
		log.Fatalf("Unable to read authorization code %v", err)
	}

	tok, err := config.Exchange(oauth2.NoContext, code)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web %v", err)
	}
	return tok
}

// tokenCacheFile generates credential file path/filename.
// It returns the generated credential path/filename.
func tokenCacheFile() (string, error) {
	usr, err := user.Current()
	if err != nil {
		return "", err
	}
	tokenCacheDir := filepath.Join(usr.HomeDir, ".credentials")
	os.MkdirAll(tokenCacheDir, 0700)
	return filepath.Join(tokenCacheDir,
		url.QueryEscape("sheets.googleapis.com-go-quickstart.json")), err
}

// tokenFromFile retrieves a Token from a given file path.
// It returns the retrieved Token and any read error encountered.
func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	t := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(t)
	defer f.Close()
	return t, err
}

// saveToken uses a file path to create a file and store the
// token in it.
func saveToken(file string, token *oauth2.Token) {
	fmt.Printf("Saving credential file to: %s\n", file)
	f, err := os.OpenFile(file, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}
