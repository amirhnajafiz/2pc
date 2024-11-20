package utils

import (
	"bufio"
	"os"
	"strings"
)

// IPTableParseFile returns a map of nodes and their IPs.
func IPTableParseFile(path string) (map[string]string, error) {
	// open the iptable file
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// create a new map
	hashMap := make(map[string]string)

	// create a file scanner
	fileScanner := bufio.NewScanner(file)
	fileScanner.Split(bufio.ScanLines)

	// read file line-by-line
	for fileScanner.Scan() {
		parts := strings.Split(fileScanner.Text(), "-")
		hashMap[parts[0]] = parts[1]
	}

	return hashMap, nil
}
