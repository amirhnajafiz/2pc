package config

import "github.com/F24-CSE535/2pc/cluster/internal/config/paxos"

func Default() Config {
	return Config{
		WatchInterval: 0, // in seconds (if set to zero, it will be disabled)
		Subnet:        6001,
		Replicas:      1,
		LogLevel:      "debug",
		MongoDB:       "mongodb://localhost:27017",
		Database:      "global",
		PaxosConfig: paxos.Config{
			CSMReplicas:        1,
			Majority:           1,
			LeaderTimeout:      10,  // in seconds
			LeaderPingInterval: 5,   // in seconds (must be less than timeout)
			ConsensusTimeout:   100, // in milliseconds
		},
	}
}
