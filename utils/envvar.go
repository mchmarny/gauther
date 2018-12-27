package utils

import (
	"log"
	"os"
)

// MustGetEnv gets sets value or sets it to default when not set
func MustGetEnv(key, defaultValue string) string {
	val := os.Getenv(key)
	if val == "" {
		val = defaultValue
	}
	if val == "" {
		log.Fatalf("Required env var (%q) not set", key)
	}
	return val
}
