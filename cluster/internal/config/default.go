package config

func Default() Config {
	return Config{
		NodeName: "",
		GRPCPort: 0,
		MongoDB:  "",
		Database: "",
	}
}
