package controllers

import (
	"net/http"
	"os"
	"path/filepath"
	"webroutes/dto"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type SysController struct {
	DB *gorm.DB
}

func NewSysController(db *gorm.DB) *SysController {
	return &SysController{DB: db}
}

// untuk membuat folder
func (sc *SysController) CreateDirectory(c *gin.Context) {

	// DTO dari directory request
	var req dto.CreateDirectoryRequest

	// Validasi folder tidak boleh kosong
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "Nama folder tidak boleh kosong!"})
		return
	}

	// Buat directory dengan hak akses izin, 0755 (full akses)
	err := os.Mkdir(req.DirectoryName, 0755)

	// Kalau status error, tampilkan error
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "Gagal membuat directory"})
		return
	}

	// Ketika sukses tampilkan sukses
	c.JSON(http.StatusCreated, gin.H{"message": "Folder berhasil dibuat!"})
}

// untuk membuat file
func (sc *SysController) CreateFile(c *gin.Context) {
	// DTO dari file Request
	var req dto.CreateFileRequest

	// Validasi tidak boleh kosong
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "Nama field wajib di isi!"})
		return
	}

	// Ini akan check folder yang di minta.
	// kalau foldernya ada, tidak perlu bikin.
	// kalau foldernya tidak ada, otomatis dibikin.
	if err := os.MkdirAll(req.DirectoryName, 0755); err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "Gagal membuat folder!"})
		return
	}

	// Menggabungkan directoryname + filename 
	filePath := filepath.Join(req.DirectoryName, req.Filename)

	// Fungsi untuk membuat sebuah file
	file, err := os.Create(filePath) 

	// Validasi jika gagal
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "Gagal membuat file!"})
		return
	}

	// Tutup file ketika tidak digunakan lagi (optimize memory)
	defer file.Close()

	// Menulis content / isi di dalam sebuah file
	_, err = file.WriteString(req.Content)

	// Validasi jika gagal
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "Gagal menulis file!"})
		return
	}

	// Response ketika sukses
	c.JSON(http.StatusCreated, gin.H{
		"message": "File berhasil dibuat dan ditulis!",
		"path": filePath,
	})
}

// untuk membaca file
func (sc *SysController) ReadFile(c *gin.Context) {
	// request user pada Read File Request
	var req dto.ReadFileRequest

	// validasi untuk nama folder dan nama file
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "Nama folder dan nama file wajib di isi!"})
		return
	}

	// Untuk menggabungkan sebuah file directory + file name
	filePath := filepath.Join(req.DirectoryName, req.Filename)

	// Untuk membaca sebuah file
	data, err := os.ReadFile(filePath)

	// Validasi ketika file gagal dibaca
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "Error membaca file!",
		})
	}

	// Untuk menampilkan response sukses ketika data berhasil dibaca
	c.JSON(http.StatusOK, gin.H{
		"data": string(data),
	})
}

// untuk rename sebuah file
func (sc *SysController) RenameFile(c *gin.Context) {
	// request file untuk DTO
	var req dto.RenameFileRequest

	// Validasi kalau folder, nama file lama, dan nama file baru tidak boleh kosong
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "Nama directory, file lama dan file baru tidak boleh kosong!"})
		return
	}

	// Untuk mengetahui keberadaan file yang lama ada dimana
	oldPath := filepath.Join(req.DirectoryName, req.OldFileName)
	
	// Untuk menentukan keberadaan file baru, yang akan di simpan
	newPath := filepath.Join(req.DirectoryName, req.NewFileName)

	// Menggunakan fungsi Rename untuk mengubah data pada folder lama, ke folder yang baru
	err := os.Rename(oldPath, newPath)

	// Kalau misal ketika file di rename, terjadi error
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "Gagal untuk rename sebuah file"})
		return
	}

	// Response ketika sukses
	c.JSON(http.StatusCreated, gin.H{
		"message": "data berhasil diubah",
		"path_lama": oldPath,
		"path_baru": newPath,
	})
}

// untuk mengupload sebuah file
func (sc *SysController) UploadFile(c *gin.Context) {
	// untuk menentukan folder upload file yang digunakan
	const uploadDir = "uploads"

	// Fungsi MkdirAll, akan check folder yang di minta.
	// kalau foldernya ada, tidak perlu bikin.
	// kalau foldernya tidak ada, otomatis dibikin.
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "Gagal membuat directory"})
		return
	}

	// Menggunakan context formFile untuk keperluan mengupload sebuah file
	file, err := c.FormFile("file")

	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "File tidak boleh kosong!",
		})
		return
	}

	// konversi name extension. misalkan taufik.txt => taufik
	filename := filepath.Base(file.Filename)

	// Untuk menggabungkan antara Upload Directory + Filename 
	PathFolder := filepath.Join(uploadDir, filename)

	// Untuk menentukan dimana sebuah file akan di upload
	if err := c.SaveUploadedFile(file, PathFolder); err != nil {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "Gagal mengupload file",
		})
	}

	// Untuk menentukan response File Upload
	response := dto.FileUploadResponse{
		Message: "File berhasil diupload",
		Filename: filename,
		Path: PathFolder,
	}

	// Untuk menampilkan sukses response
	c.JSON(http.StatusOK, gin.H{
		"data": response,
	})
}

// untuk mendownload sebuah file
func (sc *SysController) DownloadFile(c *gin.Context) {
	filename := c.Query("file_name")
	dirName := c.Query("directory_name")

	if filename == "" || dirName == "" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Filename dan Directory Name tidak boleh kosong!"})
		return
	}

	// untuk menggabungkan directoryName + FileName
	filePath := filepath.Join(dirName, filename)

	// cek apakah filenya ada atau tidak.
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		c.JSON(http.StatusForbidden, gin.H{"error": "File tidak ditemukan!"})
		return
	}

	// Untuk setting header ke content-type
	c.Header("Content-Type", "application/octet-stream")

	// Menampilkan filePath dari hasil merge
	c.File(filePath)
}