package config

// Config representa la configuración del servicio.
type Config struct {
	Env          string
	Port         string
	DatabaseURL  string
	FirebaseProj string
}

// Load retorna una configuración por defecto (placeholder).
func Load() *Config {
	return &Config{
		Env:         "dev",
		Port:        "8080",
		DatabaseURL: "postgres://shipping:shippingpass@localhost:5432/shipping_db?sslmode=disable",
	}
}
