package config

func Default() Config {
	return Config{
		WatchInterval: 1,
		Subnet:        6000,
		Replicas:      0,
		MongoDB:       "",
		Database:      "",
		LogLevel:      "debug",
	}
}
