package config

// Config representa la configuración del servicio.
// Se cargará de variables de entorno en el futuro.
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
		DatabaseURL: "postgres://user:password@localhost:5432/users?sslmode=disable",
	}
}
