package env

import (
	"log"
	"os"
)

func Must(key string) string {
	v, ok := os.LookupEnv(key)
	if !ok {
		log.Fatalf("environment variable required for %s\n", key)
	}
	return v
}
