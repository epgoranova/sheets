package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"golang.org/x/net/context"

	"github.com/catiepg/sheets/components"
)

// NOTE: If modifying these scopes, delete your previously saved credentials.
const scope = "https://www.googleapis.com/auth/spreadsheets.readonly"

// commandHandler executes a command line subcommand.
type commandHandler func(args []string) error

// subcommand specifies the names of the subcommands.
var subcommands = map[string]commandHandler{
	"auth": handleAuthCommand,
	"get":  handleGetCommand,
}

// handleAuthCommand handles the authentication subcommand.
func handleAuthCommand(args []string) error {
	auth := flag.NewFlagSet("auth", flag.ExitOnError)
	auth.Parse(args)

	config, err := components.ClientConfig(scope)
	if err != nil {
		return err
	}

	return components.CacheToken(config)
}

// handleGetCommand handles the subcommand for getting data from sheets.
func handleGetCommand(args []string) error {
	get := flag.NewFlagSet("get", flag.ExitOnError)
	id := get.String("spreadsheet", "", "ID of the spreadsheet. (Required)")
	sheet := get.String("sheet", "", "Name of the sheet. (Required)")
	column := get.String("column", "", "Column letter. (Required)")
	output := get.String("output", "", "Output path.")
	get.Parse(args)

	required := []string{"spreadsheet", "sheet", "column"}
	if err := validateRequiredFlags(get, required); err != nil {
		return err
	}

	client, err := components.Client(context.Background(), scope)
	if err != nil {
		log.Fatalf("Unable to get HTTP client: %v", err)
	}

	sheets, err := components.NewSheets(client, *id)
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

	return nil
}

// validateRequiredFlags checks if the flags are set.
func validateRequiredFlags(command *flag.FlagSet, flags []string) error {
	seen := map[string]bool{}
	command.Visit(func(flag *flag.Flag) {
		seen[flag.Name] = true
	})

	for _, flag := range flags {
		if !seen[flag] {
			return fmt.Errorf("missing required flag: -%s\n", flag)
		}
	}

	return nil
}

func main() {
	// Verify that a subcommand has been provided.
	if len(os.Args) < 2 {
		log.Fatalf("subcommand is required")
	}

	handler, ok := subcommands[os.Args[1]]
	if !ok {
		log.Fatalf("unexpected subcommand")
	}

	if err := handler(os.Args[2:]); err != nil {
		log.Fatalf("error on executing command: %s", err)
	}
}
