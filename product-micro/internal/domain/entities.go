package domain

import (
	"strings"
	"time"
)

// ==== User roles ====

const (
	UserRoleAdmin    string = "admin"
	UserRoleProducer string = "producer"
	UserRoleClient   string = "client"
)

// ==== Product Category ====
type ProductCategory string

const (
	ProductCategoryTuberculos ProductCategory = "tuberculo"
	ProductCategoryMedicional ProductCategory = "medicional"
	ProductCategoryFrutas     ProductCategory = "fruta"
	ProductCategoryVerdura    ProductCategory = "verdura"
	ProductCategoryHortaliza  ProductCategory = "hortaliza"
)

var productsCategories = map[ProductCategory]struct{}{
	ProductCategoryTuberculos: {},
	ProductCategoryMedicional: {},
	ProductCategoryFrutas:     {},
	ProductCategoryVerdura:    {},
	ProductCategoryHortaliza:  {},
}

// ==== Product Unit ====
type ProductUnit string

const (
	ProductUnitKg    ProductUnit = "kg"
	ProductUnitAtado ProductUnit = "atado"
	ProductUnitLibra ProductUnit = "libra"
)

var productsUnits = map[ProductUnit]struct{}{
	ProductUnitKg:    {},
	ProductUnitAtado: {},
	ProductUnitLibra: {},
}

// Product represents the entity of the product published by a producer
type Product struct {
	ID          string
	ProducerID  string
	Category    ProductCategory
	Price       int //in COP for the moment
	Description string
	Stock       int //cuantity of the product in the units
	Unit        ProductUnit
	PhotoUrl    string
	CreatedAt   time.Time
	DeletedAt   time.Time
}

func IsValidProductCategory(category string) bool {
	_, exist := productsCategories[ProductCategory(category)]
	return exist
}

func StringValidProductCategories() string {
	categories := make([]string, len(productsCategories))
	i := 0
	for cateogry := range productsCategories {
		categories[i] = string(cateogry)
		i++
	}
	return strings.Join(categories, ", ")
}

func IsValidProductUnit(unit string) bool {
	_, exist := productsUnits[ProductUnit(unit)]
	return exist
}

func StringValidProductUnits() string {
	units := make([]string, len(productsUnits))
	i := 0
	for unit := range productsUnits {
		units[i] = string(unit)
		i++
	}
	return strings.Join(units, ", ")
}

func GetCategoriesList() []string {
	categories := make([]string, 0, len(productsCategories))
	for category := range productsCategories {
		categories = append(categories, string(category))
	}
	return categories
}
