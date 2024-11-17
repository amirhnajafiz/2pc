package config

func Default() Config {
	return Config{
		DatastorePath: "",
		ShardsPath:    "",
	}
}
