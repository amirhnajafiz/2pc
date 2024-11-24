package handlers

import (
	"fmt"
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
		return fmt.Errorf("mismatch input arguments: count%d expected 4", argc)
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

	// get suggestions
	agr, err := db.GetAggregation(count)
	if err != nil {
		return fmt.Errorf("database failed: %v", err)
	}

	// loop and print
	for _, result := range agr {
		fmt.Printf("Accounts: %s and %s\n", result.Account1, result.Account2)
		fmt.Printf("Transaction Count: %d\n", result.TransactionCount)
		fmt.Printf("Total Amount: %.2f\n", result.TotalAmount)
		fmt.Printf("Types: %v\n", result.Types)
		fmt.Printf("Participants: %v\n", result.Participants)
		fmt.Println("---")
	}

	return nil
}
