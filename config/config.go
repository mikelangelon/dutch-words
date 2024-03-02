package config

import "os"

type Config struct {
	MongoURL string
}

func Parse() (c *Config, err error) {
	return &Config{
		MongoURL: getEnvWithFallback("MONGO_URL", "mongodb://localhost:27017/"),
	}, nil
}

func getEnvWithFallback(name string, fallback string) string {
	v := os.Getenv(name)
	if len(v) == 0 {
		return fallback
	}
	return v
}
