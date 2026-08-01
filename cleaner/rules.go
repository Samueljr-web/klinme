package cleaner

import (
	"encoding/csv"
	"io"
	"strings"
	"unicode"
)

// Record represents a single row of data
type Record map[string]string

// RuleBasedClean applies all rule-based cleaning to the records
func RuleBasedClean(records []Record) ([]Record, int) {
	cleaned := []Record{}
	cleanedCount := 0

	seen := make(map[string]bool) // for duplicate detection

	for _, record := range records {
		original := copyRecord(record)

		// Apply all rules
		record = trimWhitespace(record)
		record = fixCapitalisation(record)
		record = cleanNumericFields(record)

		// Skip empty rows
		if isEmptyRow(record) {
			continue
		}

		// Skip duplicate rows
		key := recordKey(record)
		if seen[key] {
			continue
		}
		seen[key] = true

		// Check if anything changed
		if !recordsEqual(original, record) {
			cleanedCount++
		}

		cleaned = append(cleaned, record)
	}

	return cleaned, cleanedCount
}

func trimWhitespace(record Record) Record {
	for key, value := range record {
		record[key] = strings.TrimSpace(value)
	}
	return record
}

func fixCapitalisation(record Record) Record {
	for key, value := range record {
		record[key] = toTitleCase(value)
	}
	return record
}

func cleanNumericFields(record Record) Record {
	for key, value := range record {
		if looksNumeric(value) {
			record[key] = cleanNumeric(value)
		}
	}
	return record
}

// isEmptyRow returns true if all fields in the record are empty
func isEmptyRow(record Record) bool {
	for _, value := range record {
		if value != "" {
			return false
		}
	}
	return true
}

func recordKey(record Record) string {
	var parts []string
	for _, value := range record {
		parts = append(parts, value)
	}
	return strings.Join(parts, "|")
}

func recordsEqual(a, b Record) bool {
	for key, value := range a {
		if b[key] != value {
			return false
		}
	}
	return true
}

func copyRecord(record Record) Record {
	copy := Record{}
	for key, value := range record {
		copy[key] = value
	}
	return copy
}

func toTitleCase(s string) string {
	words := strings.Fields(s)
	for i, word := range words {
		if len(word) > 0 {
			runes := []rune(word)
			runes[0] = unicode.ToUpper(runes[0])
			words[i] = string(runes)
		}
	}
	return strings.Join(words, " ")
}

// looksNumeric returns true if a string looks like it should be a number
func looksNumeric(s string) bool {
	hasDigit := false
	for _, r := range s {
		if unicode.IsDigit(r) {
			hasDigit = true
		} else if !unicode.IsDigit(r) && r != '.' && r != ',' && r != '$' && r != '%' && r != '-' && r != ' ' {
			return false
		}
	}
	return hasDigit
}

// cleanNumeric removes non-numeric characters except . and -
func cleanNumeric(s string) string {
	var result strings.Builder
	for _, r := range s {
		if unicode.IsDigit(r) || r == '.' || r == '-' {
			result.WriteRune(r)
		}
	}
	return result.String()
}

// ParseCSV reads a CSV file and returns records
func ParseCSV(r io.Reader) ([]Record, []string, error) {
	reader := csv.NewReader(r)

	// Read headers
	headers, err := reader.Read()
	if err != nil {
		return nil, nil, err
	}

	// Trim header whitespace
	for i, h := range headers {
		headers[i] = strings.TrimSpace(h)
	}

	var records []Record
	for {
		row, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, nil, err
		}

		record := Record{}
		for i, value := range row {
			if i < len(headers) {
				record[headers[i]] = value
			}
		}
		records = append(records, record)
	}

	return records, headers, nil
}

// WriteCSV writes records back to CSV format
func WriteCSV(w io.Writer, records []Record, headers []string) error {
	writer := csv.NewWriter(w)
	defer writer.Flush()

	// Write headers
	if err := writer.Write(headers); err != nil {
		return err
	}

	// Write rows
	for _, record := range records {
		row := make([]string, len(headers))
		for i, header := range headers {
			row[i] = record[header]
		}
		if err := writer.Write(row); err != nil {
			return err
		}
	}

	return nil
}
