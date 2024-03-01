package controllers

import (
	"automotive/inputters"
	"automotive/models"
	"automotive/utils"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GetUserCart(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	userID := c.MustGet("user_id").(float64)

	var carts []models.Cart

	if err := db.Preload("Car.Brand").Preload("Sales").Where("user_id = ?", userID).Select("id", "user_id", "car_id", "quantity", "is_paid", "sales_id").Find(&carts).Error; err != nil {
		response := utils.APIResponse("Cart not Found", http.StatusNotFound, "error", err)
		c.JSON(http.StatusNotFound, response)
		return
	}

	var newCartFormatters []map[string]interface{}

	for _, cart := range carts {
		cartFormatter := map[string]interface{}{
			"id":        cart.ID,
			"user_id":   cart.UserId,
			"car_id":    cart.CarId,
			"quantity":  cart.Quantity,
			"is_paid":   cart.IsPaid,
			"sales":     fmt.Sprintf("%v - %s", cart.Sales.ID, cart.Sales.Name),
			"car":       fmt.Sprintf("%s %s", cart.Car.Brand.Name, cart.Car.Name),
			"image_url": cart.Car.ImageUrl,
		}

		newCartFormatters = append(newCartFormatters, cartFormatter)
	}

	response := utils.APIResponse("Success Get Data", http.StatusOK, "success", newCartFormatters)
	c.JSON(http.StatusOK, response)
}

func UpsertCart(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)

	var input inputters.CartInput

	if err := c.ShouldBindJSON(&input); err != nil {
		response := utils.APIResponse("Failed to Process the Input", http.StatusBadRequest, "error", err.Error())
		c.JSON(http.StatusBadRequest, response)
		return
	}

	userID := c.MustGet("user_id").(float64)
	username := c.MustGet("username").(string)

	var cart models.Cart
	var car models.Car

	if err := db.First(&car, input.CarID).Error; err != nil {
		response := utils.APIResponse("Car not Found", http.StatusNotFound, "error", err)
		c.JSON(http.StatusNotFound, response)
		return
	}

	if car.Stock < input.Quantity {
		response := utils.APIResponse("Failed to Process the Input", http.StatusBadRequest, "error", "Stock Tidak Mencukupi")
		c.JSON(http.StatusBadRequest, response)
		return
	}

	if input.Quantity == 0 {
		db.Where("user_id = ? AND car_id = ?", userID, input.CarID).Delete(&cart)
	} else {
		err := db.Where("user_id = ? AND car_id = ?", userID, input.CarID).First(&cart).Error

		if err != nil {
			cart.CarId = input.CarID
			cart.Quantity = input.Quantity
			cart.UserId = int(userID)
			cart.SalesId = input.SaledID
			cart.CreatedBy = username
			_, err := cart.Create(db)
			if err != nil {
				response := utils.APIResponse("Internal Server Error", http.StatusInternalServerError, "error", err.Error())
				c.JSON(http.StatusInternalServerError, response)
				return
			}
		} else {
			cart.Quantity = input.Quantity
			cart.SalesId = input.SaledID
			if err := cart.Update(input, db, username); err != nil {
				response := utils.APIResponse("Internal Server Error", http.StatusInternalServerError, "error", err.Error())
				c.JSON(http.StatusInternalServerError, response)
				return
			}
		}
	}
	response := utils.APIResponse("Success Create Data", http.StatusOK, "success", cart)
	c.JSON(http.StatusOK, response)
}
