package domain

// NotFoundError se lanza cuando un recurso no se encuentra
type NotFoundError struct {
	Message string
}

func (e NotFoundError) Error() string {
	return e.Message
}

// ConflictError se lanza cuando hay un conflicto en la creación o actualización
type ConflictError struct {
	Message string
}

func (e ConflictError) Error() string {
	return e.Message
}

// BadRequestError se lanza cuando la entrada no es válida
type BadRequestError struct {
	Message string
}

func (e BadRequestError) Error() string {
	return e.Message
}

// InternalServerError se lanza cuando hay un error interno del servidor
type InternalServerError struct {
	Message string
}

func (e InternalServerError) Error() string {
	return e.Message
}

// UnauthorizedError se lanza cuando no hay autorización
type UnauthorizedError struct {
	Message string
}

func (e UnauthorizedError) Error() string {
	return e.Message
}
