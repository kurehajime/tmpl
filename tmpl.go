// Package tmpl is CSV Generator
package tmpl

import (
	"encoding/csv"
	"strings"
)

type Result struct {
	No    int
	Str   string
	Name  string
	Err   error
	Total int
}

//Generate returns [](template x csv)
func Generate(templ string, csvStr string, nameCol int, ch chan Result, tsv bool) {
	defer close(ch)
	records, err := getRecords(csvStr, tsv)
	if err != nil {
		ch <- Result{No: 0, Str: "", Name: "", Err: err, Total: 0}
		return
	}
	head := records[0]
	total := len(records)
	for row := 1; row < len(records); row++ {
		str := templ
		name := ""
		for col := 0; col < len(head); col++ {
			if col < len(head) && col < len(records[row]) {
				str = strings.Replace(str, head[col], records[row][col], -1)
			}
		}
		if nameCol >= 0 && nameCol < len(records[row]) {
			name = records[row][nameCol]
		}
		ch <- Result{No: row, Str: str, Name: name, Err: nil, Total: total}
	}
}
func getRecords(csvStr string, tsv bool) ([][]string, error) {
	r := csv.NewReader(strings.NewReader(csvStr))
	if tsv {
		r.Comma = []rune("\t")[0]
	}
	records, err := r.ReadAll()
	if err != nil {
		return nil, err
	}
	return records, nil
}
