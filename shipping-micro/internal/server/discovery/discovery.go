package discovery

import (
	"fmt"
	"net"
	"os"

	"github.com/hashicorp/consul/api"
	"github.com/maxskaink/proyecto-microservicios/shipping-micro/pkg/logger"
)

type ServiceRegistration struct {
	client *api.Client
}

func NewServiceRegistration() (*ServiceRegistration, error) {
	consulAddr := os.Getenv("CONSUL_HTTP_ADDR")
	if consulAddr == "" {
		consulAddr = "localhost:8500"
	}

	config := api.DefaultConfig()
	config.Address = consulAddr

	client, err := api.NewClient(config)
	if err != nil {
		return nil, fmt.Errorf("error al crear cliente Consul: %w", err)
	}

	return &ServiceRegistration{client: client}, nil
}

func (sr *ServiceRegistration) Register(serviceName, serviceID string, port int, tags []string) error {
	ip, err := getOutboundIP()
	if err != nil {
		return fmt.Errorf("no se pudo obtener outbound IP: %w", err)
	}

	registration := &api.AgentServiceRegistration{
		ID:      serviceID,
		Name:    serviceName,
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

	err = sr.client.Agent().ServiceRegister(registration)
	if err != nil {
		return fmt.Errorf("error al registrar servicio: %w", err)
	}

	logger.Info(fmt.Sprintf("Servicio %s registrado en Consul con ID %s", serviceName, serviceID))
	return nil
}

func (sr *ServiceRegistration) Deregister(serviceID string) error {
	err := sr.client.Agent().ServiceDeregister(serviceID)
	if err != nil {
		return fmt.Errorf("error al desregistrar servicio: %w", err)
	}

	logger.Info(fmt.Sprintf("Servicio %s desregistrado de Consul", serviceID))
	return nil
}

// getOutboundIP obtiene la IP preferida para conexiones externas desde el contenedor
func getOutboundIP() (net.IP, error) {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	localAddr := conn.LocalAddr().(*net.UDPAddr)
	return localAddr.IP, nil
}

// getHostname se eliminó: ya no se usa para el registro en Consul (usamos getOutboundIP)
