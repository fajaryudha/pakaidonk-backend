package config

import (
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

type JWTConfig struct {
	Issuer     string `yaml:"issuer"`
	TTLMinutes int    `yaml:"ttl_minutes"`
}

type MessageBrokerConfig struct {
	Type     string `yaml:"type"`
	RabbitMQ struct {
		URL   string `yaml:"url"`
		Queue string `yaml:"queue_name"`
	} `yaml:"rabbitmq"`
}

type AppConfig struct {
	JWT           JWTConfig           `yaml:"jwt"`
	MessageBroker MessageBrokerConfig `yaml:"message_broker"`
}

var Config AppConfig

func LoadConfig(path string) {
	data, err := os.ReadFile(path)
	if err != nil {
		log.Fatalf("Failed to read config file: %v", err)
	}
	err = yaml.Unmarshal(data, &Config)
	if err != nil {
		log.Fatalf("Failed to parse config: %v", err)
	}
}
