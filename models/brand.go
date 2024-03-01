package models

import (
	"automotive/inputters"
	"time"

	"gorm.io/gorm"
)

type (
	Brand struct {
		ID        uint      `gorm:"primary_key" json:"id"`
		Name      string    `json:"name"`
		Country   string    `json:"country"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		CreatedBy string    `json:"created_by"`
		UpdatedBy string    `json:"updated_by"`
		Cars      []Car     `json:"-"`
	}
)

func (b *Brand) Insert(db *gorm.DB) (*Brand, error) {
	if err := db.Create(&b).Error; err != nil {
		return b, err
	}
	return b, nil
}

func (b *Brand) Update(input inputters.BrandInput, updatedBy string, db *gorm.DB) error {
	updatedData := map[string]interface{}{
		"Name":      input.Name,
		"Country":   input.Country,
		"UpdatedAt": time.Now(),
		"UpdatedBy": updatedBy,
	}

	err := db.Model(b).Updates(updatedData).Error

	if err != nil {
		return err
	}

	return nil
}
