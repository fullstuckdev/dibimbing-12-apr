package controllers

import (
	"webroutes/models"
	"webroutes/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AuthController struct {
	DB *gorm.DB
}

func NewAuthController(db *gorm.DB) *AuthController {
	return &AuthController{DB: db}
}

func (ac *AuthController) Register(c *gin.Context) {
	var user models.User

	// VALIDASI MASUKAN dari USER
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(400, gin.H{"Error": err.Error()})
		return
	}

	// HASH PASSWORD
	if err := user.HashPassword(user.Password); err != nil {
		c.JSON(500, gin.H{"Error": err.Error()})
		return
	}

	// INSERT KE DATABASE
	result := ac.DB.Create(&user)
	if result.Error != nil {
		c.JSON(400, gin.H{"Error": "Error creating User"})
		return
	}

	// GENERATE TOKEN JWT
	token, err := utils.GenerateToken(user.ID) 
	if err != nil {
		c.JSON(500, gin.H{"Error": "Error generating token"})
		return
	}


	// KELUARKAN OUTPUT
	c.JSON(201, gin.H{
		"message": "User registered successfully!",
		"token": token,
	})
}

func (ac *AuthController) Login(c *gin.Context) {
    var loginReq models.LoginRequest

	// harus bentuknya JSON
    if err := c.ShouldBindJSON(&loginReq); err != nil {
        c.JSON(400, gin.H{"error": err.Error()})
        return
    }

    var user models.User

	// ngecheck datanya ada atau ga?
    if err := ac.DB.Where("email = ?", loginReq.Email).First(&user).Error; err != nil {
        c.JSON(401, gin.H{"error": "Invalid email or password"})
        return
    }

	// check password ke si JWT
    if err := user.CheckPassword(loginReq.Password); err != nil {
        c.JSON(401, gin.H{"error": "Invalid email or password"})
        return
    }

	// Generate Token
    token, err := utils.GenerateToken(user.ID)
    if err != nil {
        c.JSON(500, gin.H{"error": "Error generating token"})
        return
    }

    c.JSON(200, gin.H{
        "message": "Login successful",
        "token":   token,
    })
}