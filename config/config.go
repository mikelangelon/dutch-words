package config

import "os"

type Config struct {
	MongoURL string
}

func Parse() (c *Config, err error) {
	return &Config{
		MongoURL: getEnvWithFallback("MONGO_URL", "mongodb+srv://vercel-admin-user:iMlPRmLSCOYqy5Rq@cluster0.0rmqvms.mongodb.net/myFirstDatabase?retryWrites=true&w=majority"),
	}, nil
}

func getEnvWithFallback(name string, fallback string) string {
	v := os.Getenv(name)
	if len(v) == 0 {
		return fallback
	}
	return v
}
