package config

import "os"

var (
	JWTSecret =    GetEnv("JWT_SECRET","R1BYcTVXVGNDU2JmWHVnZ1lnN0FKeGR3cU1RUU45QXV4SDJONFZ3ckhwS1N0ZjNCYVkzZ0F4RVBSS1UzRENwRw==")
	JWTExpirationMs = GetEnv("JWT_EXPIRATION_MS","86400000")
)

func GetEnv(key, defaultValue string) string{
	value := os.Getenv(key)
	if len(value) == 0 {
		return defaultValue
	}
	return value
}