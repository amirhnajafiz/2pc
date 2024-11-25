package config

import (
	"encoding/json"
	"log"
	"strings"

	"github.com/F24-CSE535/2pc/cluster/internal/config/paxos"

	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/providers/structs"
	"github.com/knadh/koanf/v2"
	"github.com/tidwall/pretty"
)

// Prefix indicates environment variables prefix.
const Prefix = "2pc_"

// Config struct is a module that stores system configs.
type Config struct {
	Subnet                int          `koanf:"subnet"`
	Replicas              int          `koanf:"replicas"`
	ReplicasStartingIndex int          `koanf:"replicas_starting_index"`
	ClusterName           string       `koanf:"cluster_name"`
	LogLevel              string       `koanf:"log_level"`
	MongoDB               string       `koanf:"mongodb"`
	Database              string       `koanf:"database"`
	PaxosConfig           paxos.Config `koanf:"paxos"`
}

// New reads configuration with koanf, by loading a yaml config path into the Config struct.
func New(path string) Config {
	var instance Config

	k := koanf.New(".")

	// load default configuration from file
	if err := k.Load(structs.Provider(Default(), "koanf"), nil); err != nil {
		log.Fatalf("error loading default: %s", err)
	}

	// load configuration from file
	if err := k.Load(file.Provider(path), yaml.Parser()); err != nil {
		log.Printf("error loading config.yml: %s", err)
	}

	// load environment variables
	if err := k.Load(env.Provider(Prefix, ".", func(s string) string {
		return strings.ReplaceAll(strings.ToLower(
			strings.TrimPrefix(s, Prefix)), "__", ".")
	}), nil); err != nil {
		log.Printf("error loading environment variables: %s", err)
	}

	if err := k.Unmarshal("", &instance); err != nil {
		log.Fatalf("error unmarshalling config: %s", err)
	}

	indent, err := json.MarshalIndent(instance, "", "\t")
	if err != nil {
		log.Fatalf("error marshaling config to json: %s", err)
	}

	indent = pretty.Color(indent, nil)
	tmpl := `
	================ Loaded Configuration ================
	%s
	=============================================
	`
	log.Printf(tmpl, string(indent))

	return instance
}
