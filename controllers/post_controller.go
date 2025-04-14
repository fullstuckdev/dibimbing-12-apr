package controllers

import (
	"net/http"
	"webroutes/models"
	"webroutes/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)


type PostController struct {
	DB *gorm.DB
}


func NewPostController(db *gorm.DB) *PostController {
	return &PostController{DB: db}
}

func (pc *PostController) CreateTag(c *gin.Context) {
	// request yang masuk dari users
	var tag models.Tag

	// validasi request tag
	if err := utils.Validate(c, &tag); err != nil {}

	// buat bikin tag baru
	if err := pc.DB.Create(&tag).Error; err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
	}

	// response sukses
	c.JSON(http.StatusCreated, gin.H{"data": models.TagResponse{
		ID: tag.ID,
		Name: tag.Name,
	}})
}

func (pc *PostController) CreatePost(c *gin.Context) {
	// request yang masuk dari users
	var req models.CreatePostRequest

	// validasi apakah ada isinya atau kosong
	if err := c.ShouldBindBodyWithJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
	}

	// dia bakal membaca dari si userId yang ada di token JWT
	userId, exists := c.Get("userId") // userId dari token si JWT

	// kalau ga ada userId, forbidden
	if !exists {
		c.JSON(401, gin.H{"error": "Unauthorized"})
		return
	}

	// models menerima 3 data, Title, Content, dan UserId
	post := models.Post {
		Title: req.Title,  // Judul
		Content: req.Content, // Content
		UserId: userId.(uint), // UserId
	}

	// Transaksi di mulai
	tx := pc.DB.Begin()

	// Fungsi buat create post (membuat post)
	if err := tx.Create(&post).Error; err != nil {
		tx.Rollback() // ketika gagal, semuanya bakal di kembalikan ke kondisi awal.
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	// pengecekan tag.. apakah lebih dari 1
	if len(req.TagIds) > 0 {
		var tags []models.Tag // request tags, dari si models tags

		// dia bakal mencari data tags berdasarkan tagIds
		// dia searching lewat database
		if err := tx.Find(&tags, req.TagIds).Error; err != nil {
			tx.Rollback() // ketika ini di jalankan
			c.JSON(400, gin.H{"error": "invalid tag IDs"})
			return
		}

		// kalau panjangnya kurang sama dengan panjang yang data aktual.
		// dia gagal
		if len(tags) != len(req.TagIds) {
			tx.Rollback() // ketika ini di jalankan
			c.JSON(400, gin.H{"error": "Beberapa tag tidak ditemukan..."})
			return
		}

		// untuk menghubungkan table post dengan table tags. 
		// biar dia bisa masuk ke dalam table postTag
		if err := tx.Model(&post).Association("Tags").Append(&tags); err != nil {
			tx.Rollback()
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}
	}

	// ketika sudah sukses semua, commit agar data tersimpan permanen di db
	tx.Commit()

	if err := pc.DB.Preload("User").Preload("Tags").First(&post, post.ID).Error; err != nil {
		c.JSON(400, gin.H{"error": "Error loading post data"})
		return
	} 

	c.JSON(201, gin.H{"data": post})

}


func (pc *PostController) UpdatePost(c *gin.Context) {
	var req models.UpdatePostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	userId, exists := c.Get("userId")
	if !exists {
		c.JSON(401, gin.H{"error": "Unauthorized"})
		return
	}

	var post models.Post
	// Check if post exists and belongs to user
	if err := pc.DB.Where("id = ? AND user_id = ?", c.Param("id"), userId).First(&post).Error; err != nil {
		c.JSON(404, gin.H{"error": "Post not found or unauthorized"})
		return
	}

	tx := pc.DB.Begin()

	// Update basic post info
	if req.Title != "" {
		post.Title = req.Title
	}
	if req.Content != "" {
		post.Content = req.Content
	}

	if err := tx.Save(&post).Error; err != nil {
		tx.Rollback()
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	// Update tags if provided
	if len(req.TagIds) > 0 {
		var tags []models.Tag
		if err := tx.Find(&tags, req.TagIds).Error; err != nil {
			tx.Rollback()
			c.JSON(400, gin.H{"error": "Invalid tag IDs"})
			return
		}

		if len(tags) != len(req.TagIds) {
			tx.Rollback()
			c.JSON(400, gin.H{"error": "Some tags were not found"})
			return
		}

		// Replace existing tags
		if err := tx.Model(&post).Association("Tags").Replace(&tags); err != nil {
			tx.Rollback()
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}
	}

	tx.Commit()

	// Reload post with associations
	if err := pc.DB.Preload("User").Preload("Tags").First(&post, post.ID).Error; err != nil {
		c.JSON(400, gin.H{"error": "Error loading updated post"})
		return
	}

	c.JSON(200, gin.H{"data": "update success"})
}


func (pc *PostController) DeletePost(c *gin.Context) {
	userId, exists := c.Get("userId")

	if !exists {
		c.JSON(401, gin.H{"error": "Uanthorized"})
		return
	}

	var post models.Post

	if err := pc.DB.Where("id = ? AND user_id = ?", c.Param("id"), userId).First(&post).Error; err != nil {
		c.JSON(404, gin.H{"error": "Post bot found or Uanthorized"})
		return
	}

	tx := pc.DB.Begin()

	if err := tx.Model(&post).Association("Tags").Clear(); err != nil {
		tx.Rollback()
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	// Hard delete
	if err := tx.Unscoped().Delete(&post).Error; err != nil {
		tx.Rollback()
		c.JSON(400, gin.H{"Error": err.Error()})
		return
	}

	if err := tx.Delete(&post).Error; err != nil {
		tx.Rollback()
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	tx.Commit()

	c.JSON(200, gin.H{"data": "Postingan berhasil dihapus"})
}