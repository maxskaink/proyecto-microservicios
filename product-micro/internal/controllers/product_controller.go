package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/maxskaink/proyecto-microservicios/product-micro/internal/domain"
	"github.com/maxskaink/proyecto-microservicios/product-micro/internal/dto"
	"github.com/maxskaink/proyecto-microservicios/product-micro/internal/services"
)

// ProductController maneja endpoints de productos.
type ProductController struct {
	productService services.IProductService
}

// NewProductController crea una nueva instancia de ProductController.
func NewProductController(productService services.IProductService) *ProductController {
	return &ProductController{
		productService: productService,
	}
}

// RegisterRoutes registra rutas HTTP relacionadas a productos.
func (pc *ProductController) RegisterRoutes(rg *gin.RouterGroup, auth gin.HandlerFunc) {
	products := rg.Group("/products")
	{
		products.GET("", auth, pc.ListProducts)
		products.GET("/:id", auth, pc.GetProductByID)
		products.POST("", auth, pc.CreateProduct)
		products.PUT("/:id", auth, pc.UpdateProduct)
	}
}

// CreateProduct maneja POST /products
// @Summary Crear un nuevo producto
// @Description Crea un nuevo producto. Solo productores y administradores pueden crear productos.
// @Tags products
// @Accept json
// @Produce json
// @Param product body dto.ProductDTORequest true "Datos del producto"
// @Success 201 {object} dto.ProductDTOResponse
// @Failure 400 {object} dto.ErrorDTO
// @Failure 401 {object} dto.ErrorDTO
// @Failure 404 {object} dto.ErrorDTO
// @Router /products [post]
// @Security Bearer
func (pc *ProductController) CreateProduct(c *gin.Context) {
	var productReq dto.ProductDTORequest

	// Validar y parsear el request
	if err := c.ShouldBindJSON(&productReq); err != nil {
		handleUserError(c, domain.BadRequestError{Message: "Datos inválidos: " + err.Error()})
		return
	}

	// Obtener el UID del usuario del contexto (asignado por middleware)
	userUID := c.GetString("uid")
	if userUID == "" {
		handleUserError(c, domain.InternalServerError{Message: "UID de usuario no encontrado"})
		return
	}

	// Llamar al servicio
	result, err := pc.productService.CreateProduct(productReq, userUID)
	if err != nil {
		handleUserError(c, err)
		return
	}

	c.JSON(http.StatusCreated, result)
}

// GetProductByID maneja GET /products/:id
// @Summary Obtener producto por ID
// @Description Obtiene los detalles de un producto específico.
// @Tags products
// @Produce json
// @Param id path string true "ID del producto"
// @Success 200 {object} dto.ProductDTOResponse
// @Failure 404 {object} dto.ErrorDTO
// @Router /products/{id} [get]
// @Security Bearer
func (pc *ProductController) GetProductByID(c *gin.Context) {
	productID := c.Param("id")

	if productID == "" {
		handleUserError(c, domain.BadRequestError{Message: "ID de producto requerido"})
		return
	}

	// Llamar al servicio
	result, err := pc.productService.GetByIdProduct(productID)
	if err != nil {
		handleUserError(c, err)
		return
	}

	c.JSON(http.StatusOK, result)
}

// ListProducts maneja GET /products
// @Summary Listar productos con paginación
// @Description Obtiene una lista de productos con soporte para paginación.
// @Tags products
// @Produce json
// @Param page query int false "Número de página (default: 1)"
// @Param pageSize query int false "Tamaño de la página (default: 10, min: 2, max: 100)"
// @Success 200 {array} dto.ProductDTOResponse
// @Failure 400 {object} dto.ErrorDTO
// @Router /products [get]
// @Security Bearer
func (pc *ProductController) ListProducts(c *gin.Context) {
	// Obtener parámetros de query con valores por defecto
	pageStr := c.DefaultQuery("page", "1")
	pageSizeStr := c.DefaultQuery("page_size", "10")

	// Convertir a enteros
	page, err := strconv.Atoi(pageStr)
	if err != nil {
		handleUserError(c, domain.BadRequestError{Message: "Parámetro 'page' inválido"})
		return
	}

	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil {
		handleUserError(c, domain.BadRequestError{Message: "Parámetro 'pageSize' inválido"})
		return
	}

	// Llamar al servicio
	results, err := pc.productService.ListProduct(page, pageSize)
	if err != nil {
		handleUserError(c, err)
		return
	}

	c.JSON(http.StatusOK, results)
}

// UpdateProduct maneja PUT /products/:id
// @Summary Actualizar un producto
// @Description Actualiza los datos de un producto. Solo el productor propietario o administradores pueden actualizar.
// @Tags products
// @Accept json
// @Produce json
// @Param id path string true "ID del producto"
// @Param product body dto.ProductDTORequest true "Datos actualizados del producto"
// @Success 200 {object} dto.ProductDTOResponse
// @Failure 400 {object} dto.ErrorDTO
// @Failure 401 {object} dto.ErrorDTO
// @Failure 404 {object} dto.ErrorDTO
// @Router /products/{id} [put]
// @Security Bearer
func (pc *ProductController) UpdateProduct(c *gin.Context) {
	productID := c.Param("id")

	if productID == "" {
		handleUserError(c, domain.BadRequestError{Message: "ID de producto requerido"})
		return
	}

	var productReq dto.ProductDTORequest

	// Validar y parsear el request
	if err := c.ShouldBindJSON(&productReq); err != nil {
		handleUserError(c, domain.BadRequestError{Message: "Datos inválidos: " + err.Error()})
		return
	}

	// Obtener el UID del usuario del contexto (asignado por middleware)
	userUID := c.GetString("uid")
	if userUID == "" {
		handleUserError(c, domain.InternalServerError{Message: "UID de usuario no encontrado"})
		return
	}

	// Llamar al servicio
	result, err := pc.productService.UpdateProduct(productID, productReq, userUID)
	if err != nil {
		handleUserError(c, err)
		return
	}

	c.JSON(http.StatusOK, result)
}
