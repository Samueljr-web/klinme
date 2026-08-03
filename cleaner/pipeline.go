package cleaner

import "fmt"

type CleanResult struct {
	Records     []Record
	Headers     []string
	RowsIn      int
	RowsCleaned int
	RowsOut     int
}

func RunPipeline(records []Record, headers []string) (*CleanResult, error) {
	rowsIn := len(records)

	fmt.Println("Running rule-based cleaning...")
	rulesCleaned, cleanedCount := RuleBasedClean(records)

	fmt.Println(" Running AI cleaning...")
	aiCleaned, err := AIClean(rulesCleaned, headers)
	if err != nil {
		fmt.Println("AI cleaning failed, returning rule-based results:", err)
		return &CleanResult{
			Records:     rulesCleaned,
			Headers:     headers,
			RowsIn:      rowsIn,
			RowsCleaned: cleanedCount,
			RowsOut:     len(rulesCleaned),
		}, nil
	}

	return &CleanResult{
		Records:     aiCleaned,
		Headers:     headers,
		RowsIn:      rowsIn,
		RowsCleaned: cleanedCount,
		RowsOut:     len(aiCleaned),
	}, nil
}
