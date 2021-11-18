package config

import "os"

type env struct {
	Port                      string
	MONGODB_CONNECTION_STRING string
	MESSAGE_BROKER_CS         string
}

var ENV env

func init() {
	ENV = env{
		Port:                      getFromEnvOrDefault("PORT", "3001"),
		MONGODB_CONNECTION_STRING: getFromEnvOrDefault("MONGODB_CONNECTION_STRING", "mongodb://localhost:27018/"),
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
