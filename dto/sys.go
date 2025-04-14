// DTO
package dto

// Request dari users untuk membuat sebuah directory (folder)
type CreateDirectoryRequest struct {
	DirectoryName string `json:"directory_name" binding:"required"` // Required, Wajib di isi..
}

// Request dari users untuk membuat sebuah file
type CreateFileRequest struct {
	DirectoryName string `json:"directory_name" binding:"required"` // Required, Wajib di isi..
	Filename string `json:"file_name" binding:"required"` // Wajib, untuk mengupload file
	Content string `json:"content" binding:"required"` // Wajib, untuk isi contentnya
}

// Request dari users untuk membaca sebuah file
type ReadFileRequest struct {
	DirectoryName string `json:"directory_name" binding:"required"` // Required, Wajib di isi..
	Filename string `json:"file_name" binding:"required"` // Wajib, untuk membaca file
}

// Request dari users untuk mengubah nama dari suatu file
type RenameFileRequest struct {
	DirectoryName string `json:"directory_name" binding:"required"` // Required, Wajib di isi..
	OldFileName string `json:"old_file_name" binding:"required"` // Nama lama dari file kita
	NewFileName string `json:"new_file_name" binding:"required"` // Nama baru dari file kita
}

// Request dari users untuk menampilkan response ketika sudah selesai Upload
type FileUploadResponse struct {
	Message string `json:"message"` // Pesan dari responsenya
	Filename string `json:"filename"` // File yang di upload
	Path string `json:"path"` // Tempat dimana kita mengupload sebuah file
}

// Request dari users untuk mendownload sebuah file
type DownloadFileRequest struct {
	DirectoryName string `json:"directory_name" binding:"required"` // untuk mengetahui path tempat file berada
	FileName string `json:"file_name" binding:"required"` // file yang ingin kita download
}