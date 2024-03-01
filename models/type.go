package models

import (
	"automotive/inputters"
	"time"

	"gorm.io/gorm"
)

type (
	Type struct {
		ID        uint       `gorm:"primary_key" json:"id"`
		Name      string     `json:"name"`
		CreatedAt time.Time  `json:"created_at"`
		UpdatedAt *time.Time `json:"updated_at"`
		CreatedBy string     `json:"created_by"`
		UpdatedBy *string    `json:"updated_by"`
		Cars      []Car      `json:"-"`
	}
)

func (t *Type) Insert(db *gorm.DB) (*Type, error) {
	if err := db.Create(&t).Error; err != nil {
		return t, err
	}
	return t, nil
}

func (t *Type) Update(input inputters.TypeInput, updatedBy string, db *gorm.DB) error {
	updatedData := map[string]interface{}{
		"Name":      input.Name,
		"UpdatedBy": updatedBy,
		"UpdatedAt": time.Now(),
	}

	if err := db.Model(t).Updates(updatedData).Error; err != nil {
		return err
	}
	return nil
}

func (t *Type) Delete(db *gorm.DB) error {
	if err := db.Delete(&t).Error; err != nil {
		return err
	}
	return nil
}
