package controllers

import (
	"automotive/formatters"
	"automotive/inputters"
	"automotive/models"
	"automotive/utils"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var validate *validator.Validate

// Initialize the validator instance
func init() {
	validate = validator.New()
}

func GetAllUser(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)

	var user models.User

	result, err := user.GetAll(db)
	if err != nil {
		response := utils.APIResponse("Error to get Users", http.StatusInternalServerError, "error", err)
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	response := utils.APIResponse("List of Users", http.StatusOK, "success", formatters.FormatUsers(result))
	c.JSON(http.StatusOK, response)
}

func CreateUser(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	var input inputters.UserInput

	// Bind the JSON request body to the UserInput struct
	if err := c.ShouldBindJSON(&input); err != nil {
		response := utils.APIResponse("Failed to Process the Input", http.StatusBadRequest, "error", err)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	// Validate the UserInput struct
	if err := validate.Struct(input); err != nil {
		fmt.Println(err)
		response := utils.APIResponse("Validation Error", http.StatusBadRequest, "error", err.Error())
		c.JSON(http.StatusBadRequest, response)
		return
	}

	createdBy := c.MustGet("username").(string)
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)

	user := models.User{
		Name:        input.Name,
		Identity:    input.Identity,
		ImageUrl:    input.ImageUrl,
		Age:         input.Age,
		TelephoneNo: input.TelephoneNo,
		RoleId:      input.RoleId,
		Username:    input.Username,
		Password:    string(hashedPassword),
		CreatedBy:   createdBy,
		CreatedAt:   time.Now(),
	}

	insertedUser, err := user.Insert(db)
	if err != nil {
		response := utils.APIResponse("Internal Server Error", http.StatusInternalServerError, "error", err)
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	response := utils.APIResponse("Success Create Data", http.StatusOK, "success", insertedUser)
	c.JSON(http.StatusOK, response)
}

func UpdateUser(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	userID := c.Param("id")
	var input inputters.UserUpdateInput

	if err := c.ShouldBindJSON(&input); err != nil {
		response := utils.APIResponse("Failed to Process the Input", http.StatusBadRequest, "error", err)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	// Validate the UserInput struct
	if err := validate.Struct(input); err != nil {
		fmt.Println(err)
		response := utils.APIResponse("Validation Error", http.StatusBadRequest, "error", err.Error())
		c.JSON(http.StatusBadRequest, response)
		return
	}

	var existingUser models.User

	if err := db.First(&existingUser, userID).Error; err != nil {
		response := utils.APIResponse("User not Found", http.StatusNotFound, "error", err)
		c.JSON(http.StatusNotFound, response)
		return
	}

	if err := existingUser.Update(input, db); err != nil {
		response := utils.APIResponse("Internal Server Error", http.StatusInternalServerError, "error", err)
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	response := utils.APIResponse("Success Update Data", http.StatusOK, "success", input)
	c.JSON(http.StatusOK, response)
}

func DeleteUser(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	userID := c.Param("id")

	var existingUser models.User
	if err := db.First(&existingUser, &userID).Error; err != nil {
		response := utils.APIResponse("User not Found", http.StatusNotFound, "error", err)
		c.JSON(http.StatusNotFound, response)
		return
	}

	if err := existingUser.Delete(db); err != nil {
		response := utils.APIResponse("Internal Server Error", http.StatusInternalServerError, "error", err)
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	response := utils.APIResponse("Success Delete Data", http.StatusOK, "success", existingUser)
	c.JSON(http.StatusOK, response)
}
