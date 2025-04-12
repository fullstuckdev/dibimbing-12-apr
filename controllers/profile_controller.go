package controllers

import (
	"net/http"
	"webroutes/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)


type ProfileController struct {
	DB *gorm.DB
}


func NewProfileController(db *gorm.DB) *ProfileController {
	return &ProfileController{DB: db}
}

func (pc *ProfileController) CreateProfile(c *gin.Context) {

	// Check userID yang ada di token JWT
	userId, exists := c.Get("userId")

	if !exists {
		c.JSON(http.StatusForbidden, gin.H{"Error": "unauthorized"})
	}
	
	var profile models.UserProfile

	if err := c.ShouldBindJSON(&profile); err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	profile.UserID = userId.(uint)

	// Cek apakah user sudah memiliki profile
	var existingProfile models.UserProfile
	if err := pc.DB.Where("user_id = ?", userId).First(&existingProfile).Error; err == nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "Profile already exists"})
		return
	}

	// Create data
	if err := pc.DB.Create(&profile).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": profile})
}