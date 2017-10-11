package main

import (
	"flag"
	"log"

	"golang.org/x/net/context"

	"github.com/catiepg/sheets/components"
)

// NOTE: If modifying these scopes, delete your previously saved credentials.
const scope = "https://www.googleapis.com/auth/spreadsheets.readonly"

func main() {
	spreadsheetID := flag.String("spreadsheet", "", "ID of the spreadsheet. (Required)")
	sheet := flag.String("sheet", "", "Name of the sheet. (Required)")
	column := flag.String("column", "", "Column letter. (Required)")
	output := flag.String("output", "", "Output path.")
	flag.Parse()

	required := []string{"spreadsheet", "sheet", "column"}

	seen := map[string]bool{}
	flag.Visit(func(f *flag.Flag) { seen[f.Name] = true })

	for _, req := range required {
		if !seen[req] {
			log.Fatalf("missing required flag: -%s\n", req)
		}
	}

	client, err := components.Client(context.Background(), scope)
	if err != nil {
		log.Fatalf("Unable to get HTTP client: %v", err)
	}

	sheets, err := components.NewSheets(client, *spreadsheetID)
	if err != nil {
		log.Fatalf("Unable to retrieve Google Sheets service %v", err)
	}

	values, err := sheets.GetColumn(*sheet, *column)
	if err != nil {
		log.Fatalf("Unable to get column values %v", err)
	}

	if len(*output) > 0 {
		if err := components.WriteSliceToFile(*output, values); err != nil {
			log.Fatalf("Unable to write output to file %v", err)
		}

	} else {
		components.WriteSliceToStdout(values)
	}
}
