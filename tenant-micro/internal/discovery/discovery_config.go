package discovery

import (
	"fmt"
	"log"
	"net"
	"os"

	"github.com/hashicorp/consul/api"
)

// ServiceRegistration maneja el registro del servicio en Consul
type ServiceRegistration struct {
	client       *api.Client
	registration *api.AgentServiceRegistration
}

// NewServiceRegistration crea una nueva instancia de ServiceRegistration
func NewServiceRegistration() (*ServiceRegistration, error) {
	config := api.DefaultConfig()

	// Obtener dirección de Consul desde variables de entorno
	consulAddr := os.Getenv("CONSUL_HTTP_ADDR")
	if consulAddr != "" {
		config.Address = consulAddr
	}

	client, err := api.NewClient(config)
	if err != nil {
		return nil, err
	}

	return &ServiceRegistration{
		client: client,
	}, nil
}

// Register registra el servicio en Consul
func (s *ServiceRegistration) Register(id string, name string, port int, tags []string) error {
	// Obtener IP del contenedor
	ip, err := getOutboundIP()
	if err != nil {
		return err
	}

	log.Printf("Registrando servicio %s con IP %s y puerto %d", name, ip.String(), port)

	// Configurar el registro
	s.registration = &api.AgentServiceRegistration{
		ID:      id,
		Name:    name,
		Port:    port,
		Address: ip.String(),
		Tags:    tags,
		Check: &api.AgentServiceCheck{
			HTTP:                           fmt.Sprintf("http://%s:%d/health", ip.String(), port),
			Interval:                       "10s",
			Timeout:                        "3s",
			DeregisterCriticalServiceAfter: "30s",
		},
	}

	// Registrar el servicio
	return s.client.Agent().ServiceRegister(s.registration)
}

// Deregister elimina el registro del servicio de Consul
func (s *ServiceRegistration) Deregister() error {
	if s.registration == nil {
		return nil
	}
	log.Printf("Desregistrando servicio %s", s.registration.ID)
	return s.client.Agent().ServiceDeregister(s.registration.ID)
}

// getOutboundIP obtiene la IP preferida para conexiones externas
func getOutboundIP() (net.IP, error) {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)
	return localAddr.IP, nil
}
