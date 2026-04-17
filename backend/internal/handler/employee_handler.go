package handler

import (
	"net/http"
	"strconv"

	"github.com/gaston-garcia-cegid/gonsgarage/internal/core/ports"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type EmployeeHandler struct {
	employeeService ports.EmployeeService
}

func NewEmployeeHandler(employeeService ports.EmployeeService) *EmployeeHandler {
	return &EmployeeHandler{employeeService: employeeService}
}

// CreateEmployee registra un empleado (rutas protegidas: admin/manager).
// @Summary     Crear empleado
// @Tags        employees
// @Security    BearerAuth
// @Accept      json
// @Produce     json
// @Param       body body ports.CreateEmployeeRequest true "Datos del empleado"
// @Success     201 {object} map[string]interface{} "message y employee"
// @Failure     400 {object} SwaggerMessage
// @Failure     401 {object} SwaggerMessage
// @Failure     403 {object} SwaggerMessage
// @Router      /api/v1/employees [post]
func (h *EmployeeHandler) CreateEmployee(c *gin.Context) {
	var req ports.CreateEmployeeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	employee, err := h.employeeService.CreateEmployee(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":  "Employee created successfully",
		"employee": employee,
	})
}

// GetEmployee obtiene un empleado por ID interno.
// @Summary     Obtener empleado
// @Tags        employees
// @Security    BearerAuth
// @Produce     json
// @Param       id path string true "ID del empleado"
// @Success     200 {object} map[string]interface{} "employee"
// @Failure     400 {object} SwaggerMessage
// @Failure     401 {object} SwaggerMessage
// @Failure     403 {object} SwaggerMessage
// @Failure     404 {object} SwaggerMessage
// @Router      /api/v1/employees/{id} [get]
func (h *EmployeeHandler) GetEmployee(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid employee ID"})
		return
	}

	employee, err := h.employeeService.GetEmployee(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Employee not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"employee": employee})
}

// ListEmployees lista empleados con paginación.
// @Summary     Listar empleados
// @Tags        employees
// @Security    BearerAuth
// @Produce     json
// @Param       limit query int false "Límite (default 10)"
// @Param       offset query int false "Offset"
// @Param       department query string false "Filtro departamento"
// @Success     200 {object} map[string]interface{} "employees, total, limit, offset"
// @Failure     401 {object} SwaggerMessage
// @Failure     403 {object} SwaggerMessage
// @Router      /api/v1/employees [get]
func (h *EmployeeHandler) ListEmployees(c *gin.Context) {
	limitStr := c.DefaultQuery("limit", "10")
	offsetStr := c.DefaultQuery("offset", "0")
	department := c.Query("department")

	limit, _ := strconv.Atoi(limitStr)
	offset, _ := strconv.Atoi(offsetStr)

	filters := &ports.EmployeeFilters{
		Limit:  limit,
		Offset: offset,
	}

	if department != "" {
		filters.Department = &department
	}

	employees, total, err := h.employeeService.ListEmployees(c.Request.Context(), filters)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"employees": employees,
		"total":     total,
		"limit":     limit,
		"offset":    offset,
	})
}

// UpdateEmployee actualiza un empleado.
// @Summary     Actualizar empleado
// @Tags        employees
// @Security    BearerAuth
// @Accept      json
// @Produce     json
// @Param       id path string true "ID del empleado"
// @Param       body body ports.UpdateEmployeeRequest true "Campos"
// @Success     200 {object} map[string]interface{} "message y employee"
// @Failure     400 {object} SwaggerMessage
// @Failure     401 {object} SwaggerMessage
// @Failure     403 {object} SwaggerMessage
// @Router      /api/v1/employees/{id} [put]
func (h *EmployeeHandler) UpdateEmployee(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid employee ID"})
		return
	}

	var req ports.UpdateEmployeeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	employee, err := h.employeeService.UpdateEmployee(c.Request.Context(), id, req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":  "Employee updated successfully",
		"employee": employee,
	})
}

// DeleteEmployee elimina un empleado.
// @Summary     Eliminar empleado
// @Tags        employees
// @Security    BearerAuth
// @Produce     json
// @Param       id path string true "ID del empleado"
// @Success     200 {object} SwaggerMessage "message"
// @Failure     400 {object} SwaggerMessage
// @Failure     401 {object} SwaggerMessage
// @Failure     403 {object} SwaggerMessage
// @Router      /api/v1/employees/{id} [delete]
func (h *EmployeeHandler) DeleteEmployee(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid employee ID"})
		return
	}

	if err := h.employeeService.DeleteEmployee(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Employee deleted successfully"})
}
