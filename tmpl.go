// Package tmpl is CSV Generator
package tmpl

import (
	"encoding/csv"
	"strings"
)

//Generate returns [](template x csv)
func Generate(templ string, csvStr string) ([]string, error) {
	records, err := getRecords(csvStr)
	if err != nil {
		return nil, err
	}
	head := records[0]
	results := make([]string, 0, len(records)-1)
	for row := 1; row < len(records); row++ {
		str := templ
		for col := 0; col < len(head[0]); col++ {
			if col < len(head) && col < len(records[row]) {
				str = strings.Replace(str, head[col], records[row][col], -1)
			}
		}
		results = append(results, str)
	}
	return results, nil
}
func getRecords(csvStr string) ([][]string, error) {
	r := csv.NewReader(strings.NewReader(csvStr))
	records, err := r.ReadAll()
	if err != nil {
		return nil, err
	}
	return records, nil
}
