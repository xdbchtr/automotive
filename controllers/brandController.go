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

func GetAllBrand(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)

	name := c.Query("name")
	country := c.Query("country")
	carname := c.Query("carname")

	if name != "" {
		db = db.Where("name ILIKE ?", "%"+name+"%")
	}

	if country != "" {
		db = db.Where("country ILIKE ?", "%"+country+"%")
	}

	var brands []models.Brand
	db = db.Preload("Cars", func(db *gorm.DB) *gorm.DB {
		if carname != "" {
			db = db.Where("cars.name ILIKE ?", "%"+carname+"%")
		}
		return db.Select("id", "name", "year_made", "brand_id")
	}).Select("id", "name", "country").Find(&brands)

	if err := db.Error; err != nil {
		response := utils.APIResponse("Internal Server Error", http.StatusInternalServerError, "error", err)
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	var arrays []map[string]interface{}

	for _, brand := range brands {
		var carArray []map[string]interface{}
		for _, car := range brand.Cars {
			carMap := map[string]interface{}{
				"id":        car.ID,
				"name":      car.Name,
				"year_made": car.YearMade,
			}
			carArray = append(carArray, carMap)
		}

		if len(carArray) == 0 {
			carArray = []map[string]interface{}{}
		}

		brandMap := map[string]interface{}{
			"id":      brand.ID,
			"name":    brand.Name,
			"country": brand.Country,
			"cars":    carArray,
		}
		arrays = append(arrays, brandMap)
	}

	response := utils.APIResponse("List of Brands", http.StatusOK, "success", arrays)
	c.JSON(http.StatusOK, response)
}

func CreateBrand(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	var input inputters.BrandInput

	if err := c.ShouldBindJSON(&input); err != nil {
		response := utils.APIResponse("Failed to Process the Input", http.StatusBadRequest, "error", err)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	createdBy := c.MustGet("username").(string)
	brand := models.Brand{
		Name:      input.Name,
		Country:   input.Country,
		CreatedAt: time.Now(),
		CreatedBy: createdBy,
	}

	insertedBrand, err := brand.Insert(db)
	if err != nil {
		response := utils.APIResponse("Internal Server Error", http.StatusInternalServerError, "error", err)
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	response := utils.APIResponse("Success Create Data", http.StatusOK, "success", insertedBrand)
	c.JSON(http.StatusOK, response)
}

func UpdateBrand(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	brandID := c.Param("id")
	var input inputters.BrandInput

	if err := c.ShouldBindJSON(&input); err != nil {
		response := utils.APIResponse("Failed to Process the Input", http.StatusBadRequest, "error", err)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	var existingBrand models.Brand

	if err := db.First(&existingBrand, brandID).Error; err != nil {
		response := utils.APIResponse("Brand not Found", http.StatusNotFound, "error", err)
		c.JSON(http.StatusNotFound, response)
		return
	}

	updatedBy := c.MustGet("username").(string)

	if err := existingBrand.Update(input, updatedBy, db); err != nil {
		response := utils.APIResponse("Internal Server Error", http.StatusInternalServerError, "error", err)
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	response := utils.APIResponse("Success Update Data", http.StatusOK, "success", input)
	c.JSON(http.StatusOK, response)
}
