package storage

import (
	"os"
	"strconv"
)

type DbServerConfig struct {
	Server   string
	Port     int
	User     string
	Password string
	Database string
}

type DbAddrConfig struct {
	Host string
	Port int
}

type ExchangeGener struct {
	Key string
}

func NewDbSreverConfig() *DbServerConfig {
	return &DbServerConfig{
		Server:   getEnv("SERVER", ""),
		Port:     getEnvAsInt("SERVER_PORT", 1),
		User:     getEnv("SERVER_USER", ""),
		Password: getEnv("SERVER_PASSWORD", ""),
		Database: getEnv("SERVER_DATABASE", ""),
	}
}

func NewAddrServerConfig() *DbAddrConfig {
	return &DbAddrConfig{
		Host: getEnv("ADDR_HOST", ""),
		Port: getEnvAsInt("ADDR_PORT", 1),
	}
}

func NewExch() *ExchangeGener {
	return &ExchangeGener{
		Key: getEnv("API_KEY", ""),
	}
}

func getEnv(key string, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultVal
}

func getEnvAsInt(name string, defaultVal int) int {
	valueStr := getEnv(name, "")
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}
	return defaultVal
}

// func getEnvAsBool(name string, defaultVal bool) bool {
// 	valStr := getEnv(name, "")
// 	if val, err := strconv.ParseBool(valStr); err == nil {
// 		return val
// 	}
// 	return defaultVal
// }

// func getEnvAsSlice(name string, defaultVal []string, sep string) []string {
// 	valStr := getEnv(name, "")

// 	if valStr == "" {
// 		return defaultVal
// 	}
// 	val := strings.Split(valStr, sep)

// 	return val
// }
