package config

import (
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
	"log"
	"os"
)

type NatsConfig struct {
	ClusterID string `yaml:"cluster_id"`
	ClientID  string `yaml:"client_id"`
}

type HTTPServer struct {
	Address string `yaml:"address"`
}

type DataBase struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Dbname   string `yaml:"dbname"`
	Sslmode  string `yaml:"sslmode" env-default:"disable"`
}

type Config struct {
	NatsConfig `yaml:"nats_config"`
	DataBase   `yaml:"data_base"`
	HTTPServer `yaml:"http_server"`
}

func MustLoad(app string) *Config {
	var configPath string

	if err := godotenv.Load("local.env"); err != nil {
		log.Print("No .env file found")
	}

	switch app {
	case "PRODUCER":
		configPath = os.Getenv("CONFIG_PATH_PRODUCER")
	case "SUBSCRIBER":
		configPath = os.Getenv("CONFIG_PATH_SUBSCRIBER")
	}
	if configPath == "" {
		log.Fatal("Config path producer is not set")
	}
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("config file does not exist: %s", configPath)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatalf("cannot read config: %s", err)
	}
	fmt.Println(cfg)
	return &cfg
}
