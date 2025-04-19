package utils

import (
	"log"
	"os"
)

func MustGetEnv(name string) string {
	val, exists := os.LookupEnv(name)
	if !exists {
		log.Panicf("Environment variable %s is not defined\n", name)
	}
	return val
}
