package config

import "os"

func GetDBSN() string {
	return os.Getenv("DB_DSN")
}

func JwtSecret() string {
	return os.Getenv("JWT_SECRET")
}