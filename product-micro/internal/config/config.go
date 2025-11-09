package config

import (
	"os"
	"strconv"
)

// Config representa la configuración del servicio.
// Se cargará desde variables de entorno (.env en desarrollo, o envs del contenedor en prod).
type Config struct {
	Env          string
	Port         string
	DatabaseURL  string
	FirebaseProj string

	// Cloudflare R2 / S3-compatible storage
	R2Endpoint        string
	R2Region          string
	R2AccessKeyID     string
	R2SecretAccessKey string
	R2Bucket          string
	R2PublicBaseURL   string
	PresignExpiresSec int
}

// Load retorna una configuración, leyendo de variables de entorno con defaults razonables.
func Load() *Config {
	return &Config{
		Env:               getEnv("ENV", "dev"),
		Port:              getEnv("PORT", "8080"),
		DatabaseURL:       getEnv("DATABASE_URL", "postgres://user:password@localhost:5432/users?sslmode=disable"),
		FirebaseProj:      getEnv("FIREBASE_PROJECT", ""),
		R2Endpoint:        getEnv("R2Endpoint", ""),
		R2Region:          getEnv("R2_REGION", "auto"),
		R2AccessKeyID:     getEnv("R2_ACCESS_KEY_ID", ""),
		R2SecretAccessKey: getEnv("R2_SECRET_ACCESS_KEY", ""),
		R2Bucket:          getEnv("R2_BUCKET_NAME", ""),
		R2PublicBaseURL:   getEnv("R2_PUBLIC_BASE_URL", "plaza"),
		PresignExpiresSec: getEnvAsInt("PRESIGN_EXPIRES_SECONDS", 300),
	}
}

func getEnv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

func getEnvAsInt(key string, def int) int {
	if v := os.Getenv(key); v != "" {
		if i, err := strconv.Atoi(v); err == nil {
			return i
		}
	}
	return def
}
