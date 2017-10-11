# Sheets

[![GoDoc](https://godoc.org/github.com/golang/gddo?status.svg)](http://godoc.org/github.com/catiepg/sheets/components)
[![Go Report Card](https://goreportcard.com/badge/github.com/catiepg/sheets)](https://goreportcard.com/report/github.com/catiepg/sheets)

Sheets is a command line tool that provides easy reading access to Google Sheets.

## Usage

```
./sheets.exe -spreadsheet 1qpyC0XzvTcKT6EISywvqESX3A0MwQoFDE8p-Bll4hps -sheet Sheet1 -column A -output ./my-column.txt
```

#### Spreadsheet ID

The tool requires an `spreadsheet` parameter which is used to identify which
spreadsheet is to be accessed or altered. This ID is the value between the "/d/"
and the "/edit" in the URL of your spreadsheet. For example:

```
https://docs.google.com/spreadsheets/d/1qpyC0XzvTcKT6EISywvqESX3A0MwQoFDE8p-Bll4hps/edit#gid=0
```

The ID of this spreadsheet is `1qpyC0XzvTcKT6EISywvqESX3A0MwQoFDE8p-Bll4hps`.

#### Sheet ID

Individual sheets in a spreadsheet have titles and IDs. In the Sheets UI, you
can find the sheet ID of the open sheet in the spreadsheet URL, as the value of
the `gid` parameter. The following shows where the sheet ID can be found:

```
https://docs.google.com/spreadsheets/d/<spreadsheetId>/edit#gid=<sheetId>
```

In the example above, the sheet ID is `0`. Alternatively you can use the sheet
title which can be found at the bottom of the page in the browser UI. The ID of
the sheet is a required parameter and must be specified as the `sheet` flag.

#### Column

The `column` parameter is the column to be read from the sheet. It must be a single letter.

#### Output

The `output` parameter is the path to the file where the column data will be written.
Each column value will be written on a new line.

## Authentication

On first run of the tool a URL will be generated. Go to the link in your browser
and then type the authorization code in the command line. Subsequent runs will
not require authorization.
