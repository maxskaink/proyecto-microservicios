# Sistema de Manejo de Eventos con Dispatcher

## Descripción

Este es un sistema limpio y escalable para manejar eventos de RabbitMQ usando el patrón **Event Dispatcher**. 

## Arquitectura

```
Events (user.created, product.created, etc.)
    ↓
RabbitMQ Consumer
    ↓
EventProcessor (detecta y mapea el tipo)
    ↓
EventDispatcher (router central)
    ↓
Handlers específicos (UserEventHandler, ProductEventHandler, etc.)
```

## Estructura de carpetas

```
internal/messaging/
├── events/
│   ├── event.go          # Definición de tipos de eventos
│   └── user_handler.go   # Point de entrada (configuración)
├── handlers/
│   ├── dispatcher.go        # Router central de eventos
│   ├── event_processor.go   # Procesador de mensajes
│   ├── user_handler.go      # Lógica de eventos de usuario
│   └── product_handler.go   # Lógica de eventos de producto
└── rabbitmq/
    └── consumer.go       # Consumer de RabbitMQ (sin cambios)
```

## Cómo agregar un nuevo evento

### 1. Definir el tipo de evento en `events/event.go`

```go
const (
    OrderCreated EventType = "order.created"
)

type OrderCreatedEvent struct {
    Event
    OrderID string `json:"order_id"`
    UserID  string `json:"user_id"`
}
```

### 2. Crear el handler en una nueva carpeta o archivo

```go
// handlers/order_handler.go
type OrderEventHandler struct{}

func (h *OrderEventHandler) HandleOrderCreated(data []byte) error {
    var event events.OrderCreatedEvent
    if err := json.Unmarshal(data, &event); err != nil {
        logger.Error(fmt.Sprintf("Error: %v", err))
        return err
    }
    
    logger.Info(fmt.Sprintf("Orden creada: ID=%s", event.OrderID))
    // Tu lógica aquí
    return nil
}
```

### 3. Registrar el handler en `dispatcher.go`

```go
func NewEventDispatcher() *EventDispatcher {
    dispatcher := &EventDispatcher{
        handlers: make(map[events.EventType]func([]byte) error),
    }

    userHandler := &UserEventHandler{}
    productHandler := &ProductEventHandler{}
    orderHandler := &OrderEventHandler{}  // ← Agregar

    dispatcher.Register(events.UserCreated, userHandler.HandleUserCreated)
    dispatcher.Register(events.ProductCreated, productHandler.HandleProductCreated)
    dispatcher.Register(events.OrderCreated, orderHandler.HandleOrderCreated)  // ← Agregar

    return dispatcher
}
```

## Flujo de un evento

1. **RabbitMQ Consumer** recibe el mensaje
2. Llama a `EventProcessor.ProcessMessage()`
3. `EventProcessor` parsea el mensaje y extrae el tipo de evento
4. Envía el tipo al `EventDispatcher`
5. `EventDispatcher` busca el handler registrado
6. **El handler específico** ejecuta su lógica

## Ventajas

✅ **Escalable**: Agregar nuevos eventos es solo registrar un nuevo handler  
✅ **Mantenible**: Cada evento tiene su propio handler  
✅ **Testeable**: Handlers pueden probarse independientemente  
✅ **Limpio**: Separación clara de responsabilidades  
✅ **Flexible**: Fácil cambiar comportamientos sin tocar el consumer  

## Eventos actuales

| Evento | Handler | Acción |
|--------|---------|--------|
| `user.created` | `UserEventHandler` | Registra usuario creado |
| `user.deleted` | `UserEventHandler` | Limpia datos del usuario |
| `product.created` | `ProductEventHandler` | Registra producto creado |
| `product.updated` | `ProductEventHandler` | Actualiza producto |

## Cómo hacer que un evento se procese sin cambiar código

Si el evento tiene el campo `type` correctamente, el dispatcher lo procesará automáticamente sin necesidad de modificar el consumer.

```json
{
  "type": "user.created",
  "timestamp": "2025-10-30T10:00:00Z",
  "id": "user-123",
  "email": "user@example.com",
  "name": "John Doe",
  "rol": "producer"
}
```
