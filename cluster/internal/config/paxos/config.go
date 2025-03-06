package paxos

// Config holds values of paxos consensus protocol parameters.
type Config struct {
	CSMReplicas        int `koanf:"state_machine_replicas"`
	CSMBufferSize      int `koanf:"state_machine_queue_size"`
	Majority           int `koanf:"majority"`
	LeaderTimeout      int `koanf:"leader_timeout"`
	LeaderPingInterval int `koanf:"leader_ping_interval"`
	ConsensusTimeout   int `koanf:"consensus_timeout"`
}
