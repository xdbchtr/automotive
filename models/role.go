package models

import (
	"automotive/inputters"
	"time"

	"gorm.io/gorm"
)

type (
	// Role
	Role struct {
		ID          uint      `gorm:"primary_key" json:"id"`
		Name        string    `json:"name"`
		Description *string   `json:"description"`
		CreatedAt   time.Time `json:"created_at"`
		UpdatedAt   time.Time `json:"updated_at"`
		CreatedBy   string    `json:"created_by"`
		UpdatedBy   string    `json:"updated_by"`
		Users       []User    `json:"-"`
	}
)

func (r *Role) GetAll(db *gorm.DB) ([]Role, error) {
	var roles []Role
	if err := db.Select("id", "name", "description").Find(&roles).Error; err != nil {
		return nil, err
	}
	return roles, nil
}

func (r *Role) Insert(db *gorm.DB) (*Role, error) {
	if err := db.Create(&r).Error; err != nil {
		return r, err
	}
	return r, nil
}

func (r *Role) Update(input inputters.RoleInput, db *gorm.DB) error {
	updateData := map[string]interface{}{
		"Name":        input.Name,
		"Description": input.Description,
		"UpdatedBy":   "test",
		"UpdatedAt":   time.Now(),
	}
	err := db.Model(r).Updates(updateData).Error

	if err != nil {
		return err
	}
	return nil
}

func (r *Role) Delete(db *gorm.DB) error {
	if err := db.Delete(&r).Error; err != nil {
		return err
	}
	return nil
}
