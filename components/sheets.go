package components

import (
	"fmt"
	"net/http"
	"strings"

	sheets "google.golang.org/api/sheets/v4"
)

// Sheets provides access to a Google Sheets instance which works with a
// particular spreadsheet.
type Sheets struct {
	service     *sheets.Service
	spreadsheet string
}

// NewSheets creates a new Sheets instance for working with the given spreadsheet.
func NewSheets(client *http.Client, spreadsheet string) (*Sheets, error) {
	service, err := sheets.New(client)
	if err != nil {
		return nil, err
	}

	return &Sheets{
		service:     service,
		spreadsheet: spreadsheet,
	}, nil
}

// GetColumn retrieves a column from a given sheet. The column must be specified
// as a single letter.
func (s *Sheets) GetColumn(sheet, column string) ([]string, error) {
	if !isLetter(column) {
		return nil, fmt.Errorf("invalid column '%s', expected a letter", column)
	}

	upperColumn := strings.ToUpper(column)
	requestRange := fmt.Sprintf("%s!%s:%s", sheet, upperColumn, upperColumn)
	resp, err := s.service.Spreadsheets.Values.Get(s.spreadsheet, requestRange).Do()
	if err != nil {
		return nil, fmt.Errorf("Unable to retrieve data from sheet. %v", err)
	}

	var values []string
	for _, value := range resp.Values {
		values = append(values, fmt.Sprint(value[0]))
	}

	return values, nil
}

const alpha = "abcdefghijklmnopqrstuvwxyz"

func isLetter(s string) bool {
	if len(s) != 1 {
		return false
	}

	return strings.Contains(alpha, strings.ToLower(s))
}
