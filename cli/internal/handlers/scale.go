package handlers

import (
	"fmt"

	"github.com/F24-CSE535/2pc/cli/internal/database"
	"github.com/F24-CSE535/2pc/cli/pkg/models"
)

// ScaleHandler gets scale up or down commands and generates events.
type ScaleHandler struct{}

func (c *ScaleHandler) GetName() string {
	return "scale"
}

func (c *ScaleHandler) Execute(argc int, args []string) error {
	// four arguments are needed
	if argc != 4 {
		return fmt.Errorf("mismatch input arguments: count %d expected 4", argc)
	}

	// open database connection
	db, err := database.NewDatabase(args[2], args[3])
	if err != nil {
		return fmt.Errorf("open database failed: %v", err)
	}

	// create a new event
	if err := db.InsertEvent(&models.Event{
		Operation: fmt.Sprintf("scale-%s", args[0]),
		Cluster:   args[1],
	}); err != nil {
		return fmt.Errorf("failed to insert event: %v", err)
	}

	return nil
}
