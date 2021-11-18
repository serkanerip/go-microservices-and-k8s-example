package config

import (
	"os"
)

type env struct {
	CommandService            string
	Port                      string
	MONGODB_CONNECTION_STRING string
	MESSAGE_BROKER_CS         string
}

var ENV env

func init() {
	ENV = env{
		CommandService:            getFromEnvOrDefault("COMMAND_SERVICE", "http://localhost:3001"),
		Port:                      getFromEnvOrDefault("PORT", "3000"),
		MONGODB_CONNECTION_STRING: getFromEnvOrDefault("MONGODB_CONNECTION_STRING", "mongodb://localhost:27017/"),
		MESSAGE_BROKER_CS:         getFromEnvOrDefault("MESSAGE_BROKER_CS", "amqp://guest:guest@localhost:31457/"),
	}
}

func getFromEnvOrDefault(key string, defaultValue string) (e string) {
	value := defaultValue
	if e, ok := os.LookupEnv(key); ok {
		value = e
	}
	return value
}
