package controllers

import (
	"automotive/inputters"
	"automotive/models"
	"automotive/utils"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GetAllTypes(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)

	var types []models.Type
	if err := db.Select("id", "name").Find(&types).Error; err != nil {
		response := utils.APIResponse("Internal Server Error", http.StatusInternalServerError, "error", err)
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	var typesArray []map[string]interface{}

	for _, typeObject := range types {
		obj := map[string]interface{}{
			"id":   typeObject.ID,
			"name": typeObject.Name,
		}

		typesArray = append(typesArray, obj)
	}

	response := utils.APIResponse("List of Types", http.StatusOK, "success", typesArray)
	c.JSON(http.StatusOK, response)
}

func CreateType(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	var input inputters.TypeInput

	if err := c.ShouldBindJSON(&input); err != nil {
		response := utils.APIResponse("Failed to Process the Input", http.StatusBadRequest, "error", err.Error())
		c.JSON(http.StatusBadRequest, response)
		return
	}

	createdBy := c.MustGet("username").(string)
	typeData := models.Type{
		Name:      input.Name,
		CreatedAt: time.Now(),
		CreatedBy: createdBy,
	}

	insertedType, err := typeData.Insert(db)
	if err != nil {
		response := utils.APIResponse("Internal Server Error", http.StatusInternalServerError, "error", err)
		c.JSON(http.StatusInternalServerError, response)
		return
	}
	response := utils.APIResponse("Success Create Data", http.StatusOK, "success", insertedType)
	c.JSON(http.StatusOK, response)
}

func UpdateType(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	typeID := c.Param("id")
	var input inputters.TypeInput

	if err := c.ShouldBindJSON(&input); err != nil {
		response := utils.APIResponse("Failed to Process the Input", http.StatusBadRequest, "error", err)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	var existingType models.Type

	if err := db.First(&existingType, typeID).Error; err != nil {
		response := utils.APIResponse("Type not Found", http.StatusNotFound, "error", err)
		c.JSON(http.StatusNotFound, response)
		return
	}

	updatedBy := c.MustGet("username").(string)

	if err := existingType.Update(input, updatedBy, db); err != nil {
		response := utils.APIResponse("Internal Server Error", http.StatusInternalServerError, "error", err)
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	response := utils.APIResponse("Success Update Data", http.StatusOK, "success", input)
	c.JSON(http.StatusOK, response)
}

func DeleteType(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	typeID := c.Param("id")

	var existingType models.Type

	if err := db.First(&existingType, typeID).Error; err != nil {
		response := utils.APIResponse("Type not Found", http.StatusNotFound, "error", err)
		c.JSON(http.StatusNotFound, response)
		return
	}

	if err := existingType.Delete(db); err != nil {
		response := utils.APIResponse("Internal Server Error", http.StatusInternalServerError, "error", err)
		c.JSON(http.StatusInternalServerError, response)
		return
	}
	response := utils.APIResponse("Success Delete Data", http.StatusOK, "success", existingType)
	c.JSON(http.StatusOK, response)
}
