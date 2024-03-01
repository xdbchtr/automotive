package models

import (
	"automotive/inputters"
	"automotive/utils/token"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type (
	User struct {
		ID          uint      `gorm:"primary_key" json:"id"`
		Name        string    `json:"name"`
		Identity    string    `json:"identity"`
		ImageUrl    *string   `json:"image_url"`
		Age         *int      `json:"age"`
		TelephoneNo string    `json:"telephone_no"`
		RoleId      uint      `json:"role_id"`
		Username    string    `json:"username"`
		Password    string    `json:"password"`
		CreatedAt   time.Time `json:"created_at"`
		UpdatedAt   time.Time `json:"updated_at"`
		CreatedBy   string    `json:"created_by"`
		UpdatedBy   string    `json:"updated_by"`
		Role        Role      `json:"-"`
		Carts       []Cart    `gorm:"foreignKey:UserId" json:"-"`
		Sales       []Cart    `gorm:"foreignKey:SalesId" json:"-"`
	}
)

func (u *User) GetAll(db *gorm.DB) ([]User, error) {
	var users []User

	if err := db.Joins("Role").Select("users.id", "users.name", "identity", "image_url", "age", "telephone_no", "username").Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

func (u *User) Insert(db *gorm.DB) (*User, error) {
	if err := db.Create(&u).Error; err != nil {
		return u, err
	}
	return u, nil
}

func (u *User) Update(input inputters.UserUpdateInput, db *gorm.DB) error {
	updateData := map[string]interface{}{
		"Name":        input.Name,
		"Identity":    input.Identity,
		"ImageUrl":    input.ImageUrl,
		"Age":         input.Age,
		"TelephoneNo": input.TelephoneNo,
		"RoleId":      input.RoleId,
		"UpdatedAt":   time.Now(),
		"UpdatedBy":   "test",
	}

	err := db.Model(u).Updates(updateData).Error

	if err != nil {
		return err
	}
	return nil
}

func (u *User) Delete(db *gorm.DB) error {
	if err := db.Delete(&u).Error; err != nil {
		return err
	}
	return nil
}

func VerifyPassword(password, hashedPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

func LoginCheck(username string, password string, db *gorm.DB) (string, error) {
	var err error
	u := User{}

	err = db.Model(User{}).Where("username = ?", username).Take(&u).Error

	if err != nil {
		return "", err
	}

	err = VerifyPassword(password, u.Password)

	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword {
		return "", err
	}

	token, err := token.GenerateToken(u.Username, u.ID)

	if err != nil {
		return "", err
	}

	return token, nil
}
