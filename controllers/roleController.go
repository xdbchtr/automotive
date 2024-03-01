package controllers

import (
	"automotive/formatters"
	"automotive/inputters"
	"automotive/models"
	"automotive/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// GetAllRoles godoc
// @Summary Get all Roles.
// @Description Get a list of Roles.
// @Tags Role
// @Produce json
// @Success 200 {object} interface{}
// @Router /roles [get]
func GetAllRoles(c *gin.Context) {
	// get db from gin context
	db := c.MustGet("db").(*gorm.DB)

	var role models.Role

	result, err := role.GetAll(db)

	if err != nil {
		response := utils.APIResponse("Error to get Roles", http.StatusInternalServerError, "error", err)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	response := utils.APIResponse("List of Roles", http.StatusOK, "success", formatters.FormatRoles(result))
	c.JSON(http.StatusOK, response)
}

// GetAllRolesWithUsers godoc
// @Summary Get all Roles with Users.
// @Description Get a list of Roles and the Users.
// @Tags Role
// @Produce json
// @Success 200 {object} interface{}
// @Router /roles-with-users [get]
func GetAllRolesWithUsers(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)

	var roles []models.Role

	// Retrieve all roles
	if err := db.Find(&roles).Error; err != nil {
		response := utils.APIResponse("Error to get Roles with Users", http.StatusInternalServerError, "error", err)
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	var formatResponses []formatters.RoleWithUsersFormatter

	for _, role := range roles {
		var users []models.User
		if err := db.Model(&role).Association("Users").Find(&users); err != nil {
			response := utils.APIResponse("Error to retrieve Users for Role", http.StatusInternalServerError, "error", err)
			c.JSON(http.StatusInternalServerError, response)
			return
		}

		// Create RoleWithUsersResponse object
		roleResponse := formatters.RoleWithUsersFormatter{
			ID:          role.ID,
			Name:        role.Name,
			Description: role.Description,
			Users:       formatters.FormatUsers(users),
		}

		formatResponses = append(formatResponses, roleResponse)
	}

	response := utils.APIResponse("List of Roles with Users", http.StatusOK, "success", formatResponses)
	c.JSON(http.StatusOK, response)
}

// CreateRole godoc
// @Summary Create New Role.
// @Description Creating a new Role.
// @Tags Role
// @Param Body body inputters.RoleInput true "the body to create a new Role"
// @Produce json
// @Success 200 {object} interface{}
// @Router /role [post]
func CreateRole(c *gin.Context) {
	// Input Validation
	var input inputters.RoleInput
	if err := c.ShouldBindJSON(&input); err != nil {
		response := utils.APIResponse("Failed to Process the Input", http.StatusBadRequest, "error", err)
		c.JSON(http.StatusBadRequest, response)
		return
	}
	createdBy := c.MustGet("username").(string)
	role := models.Role{
		Name:        input.Name,
		Description: input.Description,
		CreatedBy:   createdBy,
	}

	db := c.MustGet("db").(*gorm.DB)

	insertedRole, err := role.Insert(db)

	if err != nil {
		response := utils.APIResponse("Failed to Process the Input", http.StatusInternalServerError, "error", err)
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	response := utils.APIResponse("Success Create Data", http.StatusOK, "success", insertedRole)
	c.JSON(http.StatusOK, response)
}

// UpdateRole godoc
// @Summary Update Role.
// @Description Update Role by id.
// @Tags Role
// @Produce json
// @Param id path string true "Role id"
// @Param Body body RoleInput true "the body to update role"
// @Success 200 {object} interface{}
// @Router /role/{id} [put]
func UpdateRole(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	roleID := c.Param("id")
	var input inputters.RoleInput

	if err := c.ShouldBindJSON(&input); err != nil {
		response := utils.APIResponse("Failed to Process the Input", http.StatusBadRequest, "error", err)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	var existingRole models.Role

	if err := db.First(&existingRole, roleID).Error; err != nil {
		response := utils.APIResponse("Role not Found", http.StatusNotFound, "error", err)
		c.JSON(http.StatusNotFound, response)
		return
	}

	if err := existingRole.Update(input, db); err != nil {
		response := utils.APIResponse("Internal Server Error", http.StatusInternalServerError, "error", err)
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	response := utils.APIResponse("Success Update Data", http.StatusOK, "success", input)
	c.JSON(http.StatusOK, response)
}

func DeleteRole(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	roleID := c.Param("id")

	var existingRole models.Role
	if err := db.First(&existingRole, &roleID).Error; err != nil {
		response := utils.APIResponse("Role not Found", http.StatusNotFound, "error", err)
		c.JSON(http.StatusNotFound, response)
		return
	}

	if err := existingRole.Delete(db); err != nil {
		response := utils.APIResponse("Internal Server Error", http.StatusInternalServerError, "error", err)
		c.JSON(http.StatusInternalServerError, response)
		return
	}
	response := utils.APIResponse("Success Delete Data", http.StatusOK, "success", existingRole)
	c.JSON(http.StatusOK, response)
}
