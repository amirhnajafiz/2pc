package handlers

// Handler interface is used in main.go to map user commands to handlers.
type Handler interface {
	Execute(argc int, args []string) error
	GetName() string
}
