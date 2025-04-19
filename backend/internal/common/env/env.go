package env

import (
	"log"
	"os"
	"strconv"
)

func MustGet(name string) string {
	val, exists := os.LookupEnv(name)
	if !exists {
		log.Panicf("Environment variable %s must be defined\n", name)
	}
	return val
}

func GetBool(name string, defaultValue bool) bool {
	val, exists := os.LookupEnv(name)
	if !exists {
		return defaultValue
	}
	boolVal, err := strconv.ParseBool(val)
	if err != nil {
		log.Panicf("Environment variable %s must be a boolean\n", name)
	}
	return boolVal
}
