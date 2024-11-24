package handlers

import (
	"fmt"
	"log"
	"strconv"

	"github.com/F24-CSE535/2pc/cli/internal/database"
	"github.com/F24-CSE535/2pc/cli/pkg/models"
)

// RebalanceHandler runs an aggregated query to perform the shards rebalance.
type RebalanceHandler struct{}

func (r *RebalanceHandler) GetName() string {
	return "rebalance"
}

func (r *RebalanceHandler) Execute(argc int, args []string) error {
	// four arguments are needed
	if argc != 4 {
		return fmt.Errorf("mismatch input arguments: count %d expected 4", argc)
	}

	// open database connection
	db, err := database.NewDatabase(args[2], args[3])
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
		fmt.Printf("accounts: %s and %s\n", result.Account1, result.Account2)
		fmt.Printf("transaction Count: %d\n", result.TransactionCount)
		fmt.Printf("total Amount: %.2f\n", result.TotalAmount)
		fmt.Printf("type: %v\n", result.Types)
		fmt.Printf("clusters: %v\n", result.Participants)

		// change shards by sending events
		if args[1] == "apply" && result.Types[0] == "cross-shard" {
			fmt.Println("making changes to shards ...")
			if err := db.InsertEvent(&models.Event{
				Cluster:   result.Participants[0],
				Operation: fmt.Sprintf("rb:%s:%s", result.Account1, result.Account2),
			}); err != nil {
				log.Printf("failed to make changes: %v\n", err)
			} else {
				fmt.Println("changes applied.")
			}
		}

		fmt.Println("---")
	}

	return nil
}
