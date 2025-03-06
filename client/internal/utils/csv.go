package utils

import (
	"encoding/csv"
	"fmt"
	"os"
	"strings"

	"github.com/F24-CSE535/2pc/client/pkg/models"
)

// CSVParseTestcaseFile accepts a testcase path and returns all testcases.
func CSVParseTestcaseFile(path string) (map[string]*models.Testcase, error) {
	list := make(map[string]*models.Testcase)

	// open CSV file
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open the CSV file %s: %v", path, err)
	}
	defer file.Close()

	// create CSV reader
	reader := csv.NewReader(file)
	reader.FieldsPerRecord = -1 // allow variable number of fields per row

	// read all data
	data, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("failed to read CSV file %s: %v", path, err)
	}

	// loop variables
	var (
		index    = ""
		servers  []string
		contacts map[string]string
		sets     []*models.Testset
	)

	// read row by row
	for i, row := range data {
		if row[0] == "" { // old row
			tmp := strings.Split(strings.Replace(strings.Replace(row[1], "(", "", -1), ")", "", -1), ", ")
			sets = append(sets, &models.Testset{
				Sender:   strings.TrimSpace(tmp[0]),
				Receiver: strings.TrimSpace(tmp[1]),
				Amount:   strings.TrimSpace(tmp[2]),
			})
		} else {
			// save the current values
			if index != "" {
				list[index] = &models.Testcase{
					LiveServers:    servers,
					ContactServers: contacts,
					Sets:           sets,
				}
			}

			// reset values
			index = row[0]
			servers = make([]string, 0)
			contacts = make(map[string]string)
			sets = make([]*models.Testset, 0)

			// set servers and contact servers
			tmpS := append(servers, strings.Split(strings.Replace(strings.Replace(row[2], "[", "", -1), "]", "", -1), ", ")...)
			for _, item := range tmpS {
				if item != "" {
					servers = append(servers, item)
				}
			}
			tmpB := strings.Split(strings.Replace(strings.Replace(row[3], "[", "", -1), "]", "", -1), ", ")
			for index, value := range tmpB {
				contacts[fmt.Sprintf("C%d", index+1)] = value
			}

			// process the first row transactions
			tmp := strings.Split(strings.Replace(strings.Replace(row[1], "(", "", -1), ")", "", -1), ", ")
			sets = append(sets, &models.Testset{
				Sender:   strings.TrimSpace(tmp[0]),
				Receiver: strings.TrimSpace(tmp[1]),
				Amount:   strings.TrimSpace(tmp[2]),
			})
		}

		// save the last set
		if i == len(data)-1 {
			list[index] = &models.Testcase{
				LiveServers:    servers,
				ContactServers: contacts,
				Sets:           sets,
			}
		}
	}

	return list, nil
}
