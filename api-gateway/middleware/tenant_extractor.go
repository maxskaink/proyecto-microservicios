package main

import (
	"context"
	"net/http"
	"regexp"
)

// TenantExtractorMiddleware extrae el tenant de la ruta y lo pasa como header
func TenantExtractorMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Regex para extraer el tenant de la ruta: /tenantX/...
		tenantRegex := regexp.MustCompile(`^/([a-zA-Z0-9_-]+)/api/(.*)$`)
		matches := tenantRegex.FindStringSubmatch(r.URL.Path)

		if len(matches) == 3 {
			tenantID := matches[1]
			apiPath := matches[2]

			// Reescribir la ruta sin el tenant
			r.URL.Path = "/api/" + apiPath

			// Agregar el tenant al header
			r.Header.Set("X-Tenant-ID", tenantID)

			// Loguear la extracción del tenant
			// log.Printf("Tenant extraído: %s, Nueva ruta: %s", tenantID, r.URL.Path)
		}

		// Si no coincide, significa que no hay tenant en la ruta
		// continuamos sin agregar el header (para rutas públicas)

		next.ServeHTTP(w, r)
	})
}

// Plugin es la estructura que Traefik espera para los plugins
type Config struct{}

// CreateConfig crea la configuración del plugin
func CreateConfig() *Config {
	return &Config{}
}

// New crea una nueva instancia del middleware
func New(ctx context.Context, next http.Handler, config *Config, name string) (http.Handler, error) {
	return TenantExtractorMiddleware(next), nil
}

// SetupPlugin se encarga de la configuración del plugin en Traefik
func init() {}

// TenantExtractorPlugin es el plugin en sí
type TenantExtractorPlugin struct {
	next http.Handler
	name string
}

// ServeHTTP implementa http.Handler
func (p *TenantExtractorPlugin) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Extraer tenant de la ruta
	tenantRegex := regexp.MustCompile(`^/([a-zA-Z0-9_-]+)/api/(.*)$`)
	matches := tenantRegex.FindStringSubmatch(r.URL.Path)

	if len(matches) == 3 {
		tenantID := matches[1]
		apiPath := matches[2]

		// Reescribir la ruta
		r.URL.Path = "/api/" + apiPath

		// Agregar header
		r.Header.Set("X-Tenant-ID", tenantID)
	}

	p.next.ServeHTTP(w, r)
}
