package models

import (
	"automotive/inputters"
	"time"

	"gorm.io/gorm"
)

type (
	Cart struct {
		ID             uint      `gorm:"primary_key" json:"id"`
		UserId         int       `json:"user_id"`
		CarId          int       `json:"car_id"`
		Quantity       int       `json:"quantity"`
		ShipmentMethod *string   `json:"shipment_method"`
		NoResi         *string   `json:"no_resi"`
		IsPaid         bool      `json:"is_paid"`
		SalesId        int       `json:"sales_id"`
		CreatedAt      time.Time `json:"created_at"`
		UpdatedAt      time.Time `json:"updated_at"`
		CreatedBy      string    `json:"created_by"`
		UpdatedBy      *string   `json:"updated_by"`
		Car            Car       `json:"-"`
		User           User      `json:"-"`
		Sales          User      `gorm:"foreignKey:SalesId" json:"-"`
	}
)

func (c *Cart) Create(db *gorm.DB) (*Cart, error) {
	if err := db.Create(&c).Error; err != nil {
		return c, err
	}
	return c, nil
}

func (c *Cart) Insert(db *gorm.DB) (*Cart, error) {
	if err := db.Create(&c).Error; err != nil {
		return c, err
	}
	return c, nil
}

func (c *Cart) Update(input inputters.CartInput, db *gorm.DB, updatedBy string) error {
	updateData := map[string]interface{}{
		"Quantity":  input.Quantity,
		"SalesId":   input.SaledID,
		"UpdatedAt": time.Now(),
		"UpdatedBy": updatedBy,
	}
	err := db.Model(c).Updates(updateData).Error

	if err != nil {
		return err
	}
	return nil
}
