package controllers

import (
	"automotive/inputters"
	"automotive/models"
	"automotive/utils"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GetAllCars(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)

	var cars []models.Car
	if err := db.Preload("Brand", func(db *gorm.DB) *gorm.DB {
		return db.Select("id, name")
	}).Preload("Type", func(db *gorm.DB) *gorm.DB {
		return db.Select("id, name")
	}).Find(&cars).Error; err != nil {
		response := utils.APIResponse("Internal Server Error", http.StatusInternalServerError, "error", err)
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	var carArray []map[string]interface{}

	for _, car := range cars {

		numberStr := strconv.Itoa(int(car.Price))

		// Split the number into parts (whole number and decimal)
		parts := strings.Split(numberStr, ".")

		// Format the whole number part
		formattedWhole := formatWholeNumber(parts[0])

		// If there's a decimal part, append it
		formattedNumber := formattedWhole
		if len(parts) > 1 {
			formattedNumber += "," + parts[1]
		}

		// Add currency symbol
		formattedNumber = "USD $" + formattedNumber

		carMap := map[string]interface{}{
			"id":         car.ID,
			"name":       car.Name,
			"brand_name": car.Brand.Name,
			"year_made":  car.YearMade,
			"image_url":  car.ImageUrl,
			"price":      formattedNumber,
			"stock":      car.Stock,
			"type_name":  car.Type.Name,
		}

		carArray = append(carArray, carMap)
	}

	response := utils.APIResponse("List of Cars", http.StatusOK, "success", carArray)
	c.JSON(http.StatusOK, response)
}

func UpdateCar(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	carID := c.Param("id")
	var input inputters.CarInput

	if err := c.ShouldBindJSON(&input); err != nil {
		response := utils.APIResponse("Failed to Process the Input", http.StatusBadRequest, "error", err.Error())
		c.JSON(http.StatusBadRequest, response)
		return
	}

	var existingBrand models.Car
	var checkBrand models.Brand

	if err := db.First(&existingBrand, carID).Error; err != nil {
		response := utils.APIResponse("Car not Found", http.StatusNotFound, "error", err)
		c.JSON(http.StatusNotFound, response)
		return
	}

	if err := db.First(&checkBrand, input.BrandID).Error; err != nil {
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

func CreateCar(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	var input inputters.CarInput

	if err := c.ShouldBindJSON(&input); err != nil {
		response := utils.APIResponse("Failed to Process the Input", http.StatusBadRequest, "error", err.Error())
		c.JSON(http.StatusBadRequest, response)
		return
	}

	var brand models.Brand

	if err := db.First(&brand, input.BrandID).Error; err != nil {
		response := utils.APIResponse("Brand not Found", http.StatusNotFound, "error", err)
		c.JSON(http.StatusNotFound, response)
		return
	}

	createdBy := c.MustGet("username").(string)
	car := models.Car{
		Name:        input.Name,
		BrandID:     input.BrandID,
		YearMade:    input.YearMade,
		Price:       input.Price,
		Stock:       input.Stock,
		TypeID:      input.TypeID,
		ImageUrl:    input.ImageUrl,
		Description: input.Description,
		CreatedAt:   time.Now(),
		CreatedBy:   createdBy,
	}

	insertedCar, err := car.Insert(db)
	if err != nil {
		response := utils.APIResponse("Internal Server Error", http.StatusInternalServerError, "error", err)
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	response := utils.APIResponse("Success Create Data", http.StatusOK, "success", insertedCar)
	c.JSON(http.StatusOK, response)
}

func formatWholeNumber(whole string) string {
	// Add dots every three digits
	var formatted string
	for i := len(whole); i > 0; i -= 3 {
		if i-3 > 0 {
			formatted = "." + whole[i-3:i] + formatted
		} else {
			formatted = whole[:i] + formatted
		}
	}
	return formatted
}

func DeleteCar(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	carID := c.Param("id")

	var existingCar models.Car
	if err := db.First(&existingCar, &carID).Error; err != nil {
		response := utils.APIResponse("Car not Found", http.StatusNotFound, "error", err)
		c.JSON(http.StatusNotFound, response)
		return
	}

	if err := existingCar.Delete(db); err != nil {
		response := utils.APIResponse("Internal Server Error", http.StatusInternalServerError, "error", err)
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	response := utils.APIResponse("Success Delete Data", http.StatusOK, "success", existingCar)
	c.JSON(http.StatusOK, response)
}
