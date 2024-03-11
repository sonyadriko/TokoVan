package main

import (
	"log"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var db *gorm.DB

type Barang struct {
	gorm.Model
	Nama  string `json:"nama"`
	Harga string `json:"harga"`
	Stok  string `json:"stok"`
}

func init() {
	// Mengganti informasi koneksi sesuai dengan konfigurasi MySQL Anda
	dsn := "root:@tcp(localhost:3306)/tokovan?parseTime=true"
	var err error
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}
}

// TambahBarang menambahkan data barang ke database
func getBarang(c *gin.Context) {
	var barang []Barang
	if err := db.Find(&barang).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, barang)
}

// AmbilSemuaBarang mengambil semua data barang dari database
func getBarangByID(c *gin.Context) {
	id := c.Param("id")
	var barang Barang
	if err := db.First(&barang, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Barang not found"})
		return
	}
	c.JSON(http.StatusOK, barang)
}

// AmbilBarangByID mengambil data barang berdasarkan ID dari database
func tambahBarang(c *gin.Context) {
	var input Barang
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := db.Create(&input).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, input)
}

// PerbaruiHargaBarang memperbarui harga barang berdasarkan ID di database
func perbaruiBarang(c *gin.Context) {
	id := c.Param("id")
	var barang Barang
	if err := db.First(&barang, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Barang not found"})
		return
	}

	var input Barang
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	barang.Nama = input.Nama
	barang.Harga = input.Harga
	barang.Stok = input.Stok

	if err := db.Save(&barang).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, barang)
}

// HapusBarang menghapus data barang berdasarkan ID dari database
func hapusBarang(c *gin.Context) {
	id := c.Param("id")
	var barang Barang
	if err := db.First(&barang, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Barang not found"})
		return
	}

	if err := db.Delete(&barang).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Barang deleted successfully"})
}

func softDeleteBarang(c *gin.Context) {
	id := c.Param("id")
	var barang Barang
	if err := db.First(&barang, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Barang not found"})
		return
	}

	// Melakukan soft delete dengan mengatur DeletedAt menjadi waktu saat ini
	if err := db.Delete(&barang).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Barang soft deleted successfully"})
}

func main() {
	r := gin.Default()

	// Middleware CORS
	corsMiddleware := cors.New(cors.Config{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders: []string{"Content-Type"},
	})

	// Menggunakan middleware CORS pada router GIN
	r.Use(corsMiddleware)

	// Endpoints untuk CRUD
	r.GET("/barang", getBarang)
	r.GET("/barang/:id", getBarangByID)
	r.POST("/barang", tambahBarang)
	r.PUT("/barang/:id", perbaruiBarang)
	r.DELETE("/barang/:id", softDeleteBarang)

	// Menjalankan server
	r.Run(":8080")
}
