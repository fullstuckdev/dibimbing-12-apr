package models

import (
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// One to One
type User struct {
	gorm.Model
	Name string `json:"name"`
	Email string `json:"email"`
	Password string `json:"password"`
	Profile UserProfile `json:"profile" gorm:"foreignKey:UserID"`
}

type UserProfile struct {
	gorm.Model
	UserID uint `gorm:"unique"` // Foreign Key
	Address string `json:"address"`
	Bio string `json:"bio"`
}

// DTO
type LoginRequest struct {
	Email string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// HASH PASSWORD 
func (u *User) HashPassword(password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	u.Password = string(hashedPassword)

	return nil
}

// PENGECEKAN PASSWORD
func (u *User) CheckPassword(password string) error {
	return bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
}