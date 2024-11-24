package handlers

import (
	"fmt"
	"os"
	"strconv"

	"github.com/F24-CSE535/2pc/cli/internal/database"
)

// RebalanceHandler runs an aggregated query to perform the shards rebalance.
type RebalanceHandler struct{}

func (r *RebalanceHandler) GetName() string {
	return "rebalance"
}

func (r *RebalanceHandler) Execute(argc int, args []string) error {
	// four arguments are needed
	if argc != 3 {
		return fmt.Errorf("mismatch input arguments: count %d expected 3", argc)
	}

	// open database connection
	db, err := database.NewDatabase(args[1], args[2])
	if err != nil {
		return fmt.Errorf("open database failed: %v", err)
	}

	// get the count
	count, err := strconv.Atoi(args[0])
	if err != nil {
		return fmt.Errorf("input argument is not a number: %v", err)
	}

	// get aggregations
	aggregations, err := db.GetAggregation(count)
	if err != nil {
		return fmt.Errorf("database failed: %v", err)
	}

	// open the rebalance schema file
	file, err := os.Create("new-schema.txt")
	if err != nil {
		return fmt.Errorf("failed to open new-schema file: %v", err)
	}
	defer file.Close()

	// loop and store them inside file
	for _, result := range aggregations {
		if result.Types[0] == "cross-shard" {
			record := fmt.Sprintf(
				"{clusters: %v, accounts: [%s, %s], transactions: %d, total: %2f}\n",
				result.Participants,
				result.Account1,
				result.Account2,
				result.TransactionCount,
				result.TotalAmount,
			)
			_, err = file.WriteString(record)
			if err != nil {
				return fmt.Errorf("failed to store record: %v", err)
			}
		}
	}

	return nil
}
