# Observabilidad con Prometheus y Grafana

Este directorio contiene la configuración para el sistema de observabilidad del proyecto de microservicios.

## 🏗️ Arquitectura

```
Grafana (Visualización)
    ↓
Prometheus (Recolección)
    ↓
    ├── users-micro:8080/metrics
    ├── product-micro:8080/metrics
    ├── tenant-micro:8080/metrics
    ├── shipping-micro:8080/metrics
    ├── Consul:8500/v1/agent/metrics
    └── RabbitMQ:15692/metrics
```

## 📊 Servicios

### Prometheus
- **URL**: http://localhost:9090
- **Puerto**: 9090
- **Configuración**: `prometheus/prometheus.yml`
- **Función**: Recolecta métricas de todos los microservicios cada 15 segundos

### Grafana
- **URL**: http://localhost:3000
- **Credenciales**: `admin` / `admin`
- **Puerto**: 3000
- **Datasources**: Auto-provisionado con Prometheus

## 📈 Métricas Disponibles

### HTTP Metrics
- `http_requests_total` - Total de requests HTTP
  - Labels: `method`, `endpoint`, `status`
- `http_request_duration_seconds` - Duración de requests (histograma)
  - Labels: `method`, `endpoint`

### Database Metrics
- `db_query_duration_seconds` - Duración de queries
  - Labels: `operation`, `table`
- `db_connections_active` - Conexiones activas a la BD

## 🔍 Consultas Útiles en Prometheus

### Tasa de Requests por Segundo
```promql
rate(http_requests_total[5m])
```

### Latencia P95 por Servicio
```promql
histogram_quantile(0.95, rate(http_request_duration_seconds_bucket[5m]))
```

### Tasa de Errores (5xx)
```promql
rate(http_requests_total{status=~"5.."}[5m])
```

### Requests por Método y Status
```promql
sum by (method, status) (rate(http_requests_total[5m]))
```

### Duración Promedio de Queries DB
```promql
rate(db_query_duration_seconds_sum[5m]) / rate(db_query_duration_seconds_count[5m])
```

## 📊 Dashboards en Grafana

### Importar Dashboards Públicos

1. Accede a Grafana: http://localhost:3000
2. Ve a **Dashboards** → **New** → **Import**
3. Usa los siguientes IDs:
   - **10826** - Go Application Metrics
   - **10991** - RabbitMQ Overview
   - **12049** - Consul Cluster Monitoring
   - **9628** - PostgreSQL Database

### Crear Dashboard Personalizado

1. Ve a **Dashboards** → **New Dashboard**
2. Agrega paneles con las consultas de arriba
3. Configura alertas si es necesario

## 🔧 Endpoints de Métricas

Cada microservicio expone sus métricas en `/metrics`:

- Users: http://localhost:8080/metrics
- Products: http://localhost:8081/metrics
- Tenants: http://localhost:8082/metrics
- Shipping: http://localhost:8083/metrics

## 🚀 Uso

### Iniciar el Sistema
```bash
./start.sh
```

### Verificar Targets en Prometheus
1. Accede a http://localhost:9090/targets
2. Verifica que todos los servicios estén "UP"

### Explorar Métricas
1. Ve a http://localhost:9090/graph
2. Prueba consultas PromQL
3. Visualiza gráficos

### Crear Visualizaciones en Grafana
1. Accede a http://localhost:3000
2. Crea un nuevo dashboard
3. Agrega paneles con consultas Prometheus

## 📝 Configuración

### Agregar Nuevos Targets

Edita `prometheus/prometheus.yml`:

```yaml
scrape_configs:
  - job_name: 'nuevo-servicio'
    static_configs:
      - targets: ['nuevo-servicio:puerto']
    metrics_path: '/metrics'
```

### Cambiar Intervalo de Scrape

En `prometheus/prometheus.yml`:

```yaml
global:
  scrape_interval: 15s  # Cambiar aquí
```

## 🔔 Alertas (Opcional)

Para configurar alertas, crea `prometheus/alerts.yml`:

```yaml
groups:
  - name: microservices
    rules:
      - alert: HighErrorRate
        expr: rate(http_requests_total{status=~"5.."}[5m]) > 0.05
        for: 5m
        labels:
          severity: warning
        annotations:
          summary: "Alto rate de errores"
```

## 📚 Referencias

- [Prometheus Documentation](https://prometheus.io/docs/)
- [Grafana Documentation](https://grafana.com/docs/)
- [PromQL Basics](https://prometheus.io/docs/prometheus/latest/querying/basics/)
- [Prometheus Client Go](https://github.com/prometheus/client_golang)
