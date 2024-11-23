package models

// Testset is a transaction holder.
type Testset struct {
	Sender   string
	Receiver string
	Amount   string
}

// Testcase contains a scenario parameters to test.
type Testcase struct {
	ContactServers map[string]string
	LiveServers    []string
	Sets           []*Testset
}
