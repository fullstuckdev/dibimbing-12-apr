package controllers

import (
	"net/http"
	"webroutes/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type UserController struct {
	DB *gorm.DB
}

// Bikin variabel buat menampung
var usersInMemory = []models.User{}

func NewUserController(db *gorm.DB) *UserController {
	return &UserController{DB: db}
}

// GET user tanpa DB
func GetUserWithoutDB(c *gin.Context) {
	// c.json = formatting JSON
	// gin.H = penampung data
	c.JSON(http.StatusOK, gin.H{"data": usersInMemory})
}

// CREATE user tanpa DB
func CreateUserWithoutDB(c *gin.Context) {
	// Data User
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}
	// misal panjang array = 0.. berarti user.ID = 1
	// misal panjang array = 1.. berarti user.ID = 2
	user.ID = uint(len(usersInMemory) + 1)

	// Buat insert data ke dalam array
	usersInMemory = append(usersInMemory, user)

	// Sukses insert data
	c.JSON(http.StatusCreated, gin.H{"data": user})
}

// GET USERS dengan DATABASE
func (uc *UserController) GetUsers(c *gin.Context) {
	var users []models.User
	uc.DB.Find(&users)

	c.JSON(http.StatusOK, gin.H{"data": users})
}

// CREATE User dengan DATABASE
func(uc *UserController) CreateUser(c *gin.Context) {
	var users models.User

	if err := c.ShouldBindJSON(&users); err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}

	if err := users.HashPassword(users.Password); err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}

	result := uc.DB.Create(&users)

	if result.Error != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": users})
}