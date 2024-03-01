package models

import (
	"automotive/inputters"
	"time"

	"gorm.io/gorm"
)

type (
	Car struct {
		ID          uint      `gorm:"primary_key" json:"id"`
		Name        string    `json:"name"`
		BrandID     int       `json:"brand_id"`
		YearMade    int       `json:"year_made"`
		Price       float64   `json:"price"`
		Stock       int       `json:"stock"`
		TypeID      *int      `json:"type_id"`
		ImageUrl    *string   `json:"image_url"`
		Description *string   `json:"description"`
		CreatedAt   time.Time `json:"created_at"`
		UpdatedAt   time.Time `json:"updated_at"`
		CreatedBy   string    `json:"created_by"`
		UpdatedBy   string    `json:"updated_by"`
		Type        Type      `json:"-"`
		Brand       Brand     `json:"-"`
		Carts       []Cart    `json:"-"`
	}
)

func (c *Car) Insert(db *gorm.DB) (*Car, error) {
	if err := db.Create(&c).Error; err != nil {
		return c, err
	}
	return c, nil
}

func (c *Car) Update(input inputters.CarInput, updatedBy string, db *gorm.DB) error {
	updatedData := map[string]interface{}{
		"Name":        input.Name,
		"BrandID":     input.BrandID,
		"YearMade":    input.YearMade,
		"Price":       input.Price,
		"Stock":       input.Stock,
		"TypeID":      input.TypeID,
		"ImageUrl":    input.ImageUrl,
		"Description": input.Description,
	}

	err := db.Model(c).Updates(updatedData).Error

	if err != nil {
		return err
	}
	return nil
}

func (c *Car) Delete(db *gorm.DB) error {
	if err := db.Delete(&c).Error; err != nil {
		return err
	}
	return nil
}
