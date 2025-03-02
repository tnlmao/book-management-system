package config

import (
	"log"
	"os"
)

type Config struct {
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	RedisURL   string
	KafkaURL   string
}

var CFG Config

func LoadConfig() {
	CFG = Config{
		DBHost:     os.Getenv("DB_HOST"),
		DBPort:     os.Getenv("DB_PORT"),
		DBUser:     os.Getenv("DB_USER"),
		DBPassword: os.Getenv("DB_PASSWORD"),
		DBName:     os.Getenv("DB_NAME"),
		RedisURL:   os.Getenv("REDIS_URL"),
		KafkaURL:   os.Getenv("KAFKA_URL"),
	}

	// Validate required fields
	if CFG.DBHost == "" || CFG.DBUser == "" || CFG.DBPassword == "" || CFG.DBName == "" {
		log.Fatal("Missing required database configuration")
	}
	if CFG.RedisURL == "" {
		log.Fatal("Missing Redis configuration")
	}
	if CFG.KafkaURL == "" {
		log.Fatal("Missing Kafka configuration")
	}

}
