package validators

import (
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/maxskaink/proyecto-microservicios/product-micro/internal/domain"
)

// RegisterValidators registra validadores personalizados
func RegisterValidators() {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("category", validateCategory)
		v.RegisterValidation("unit", validateUnit)
	}
}

// validateCategory valida que sea una categoría válida usando la función del dominio
func validateCategory(fl validator.FieldLevel) bool {
	category := fl.Field().String()
	return domain.IsValidProductCategory(category)
}

// validateUnit valida que sea una unidad válida usando la función del dominio
func validateUnit(fl validator.FieldLevel) bool {
	unit := fl.Field().String()
	return domain.IsValidProductUnit(unit)
}
