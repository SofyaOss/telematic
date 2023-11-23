package config

import (
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

type AppConf struct {
	Redis struct {
		Port string
		Host string
	}
	Kafka struct {
		Port string
		Host string
	}
	GRPC struct {
		Port string
		Host string
	}
}

func ReadConf(cfg *AppConf) *AppConf {
	f, err := os.Open("/app/internal/config/config.yaml")
	if err != nil {
		log.Fatalf("Could not open config file: %s", err)
	}
	defer f.Close()

	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(&cfg)
	if err != nil {
		log.Fatalf("Could not decode config file: %s", err)
	}
	return cfg
}
