package utils

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
)

// CSVParseDatastoreFile to get client names and their balances.
func CSVParseDatastoreFile(path string) (map[string]int, error) {
	// open the CSV file
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// create a CSV reader
	reader := csv.NewReader(file)
	reader.FieldsPerRecord = -1 // allow variable number of fields per row

	// read all data
	data, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("failed to read CSV file %s: %v", path, err)
	}

	// create a map
	hashMap := make(map[string]int)

	// loop over rows
	for _, row := range data[1:] {
		tmp, err := strconv.Atoi(row[1])
		if err != nil {
			return nil, err
		}

		hashMap[row[0]] = tmp
	}

	return hashMap, nil
}

// CSVPaseShardsFile to get shards and clusters info.
func CSVParseShardsFile(path string) ([][]string, error) {
	// open the CSV file
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// create a CSV reader
	reader := csv.NewReader(file)
	reader.FieldsPerRecord = -1 // allow variable number of fields per row

	// read all data
	data, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("failed to read CSV file %s: %v", path, err)
	}

	return data, nil
}
