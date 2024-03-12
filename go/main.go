package main

import (
	"errors"
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
	Nama  string  `json:"nama"`
	Harga float64 `json:"harga"`
	Stok  int     `json:"stok"`
	// tambahkan field lain sesuai kebutuhan
}

// type Transaksi struct {
// 	ID      uint      `gorm:"primary_key" json:"id"`
// 	Tanggal time.Time `json:"tanggal"`
// 	// tambahkan field lain sesuai kebutuhan
// }

type Transaksi struct {
	ID         uint    `gorm:"primary_key" json:"id"`
	Tanggal    string  `json:"tanggal"`
	Barang     string  `json:"barang"`
	Jumlah     int     `json:"jumlah"`
	HargaTotal float64 `json:"harga_total"`
}

type DetailTransaksi struct {
	ID          uint    `gorm:"primary_key" json:"id"`
	IDTransaksi uint    `json:"id_transaksi"`
	IDBarang    uint    `json:"id_barang"`
	Jumlah      int     `json:"jumlah"`
	HargaSatuan float64 `json:"harga_satuan"`
	TotalHarga  float64 `json:"total_harga"`
	// tambahkan field lain sesuai kebutuhan
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
	if err := db.Where("deleted_at IS NULL").Find(&barang).Error; err != nil {
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

func getTransaksi(c *gin.Context) {
	var transaksi []Transaksi
	if err := db.Find(&transaksi).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, transaksi)
}

// getTransaksiByID mengambil data transaksi berdasarkan ID dari database
func getTransaksiByID(c *gin.Context) {
	id := c.Param("id")
	var transaksi Transaksi
	if err := db.First(&transaksi, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Transaksi not found"})
		return
	}
	c.JSON(http.StatusOK, transaksi)
}

// func tambahTransaksi(c *gin.Context) {
// 	var input Transaksi
// 	if err := c.ShouldBindJSON(&input); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		return
// 	}

// 	if err := db.Create(&input).Error; err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 		return
// 	}

// 	c.JSON(http.StatusOK, input)
// }

// PerbaruiTransaksi memperbarui data transaksi berdasarkan ID di database
func perbaruiTransaksi(c *gin.Context) {
	id := c.Param("id")
	var transaksi Transaksi
	if err := db.First(&transaksi, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Transaksi not found"})
		return
	}

	var input Transaksi
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	transaksi.Tanggal = input.Tanggal

	if err := db.Save(&transaksi).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, transaksi)
}

// HapusTransaksi menghapus data transaksi berdasarkan ID dari database
func hapusTransaksi(c *gin.Context) {
	id := c.Param("id")
	var transaksi Transaksi
	if err := db.First(&transaksi, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Transaksi not found"})
		return
	}

	if err := db.Delete(&transaksi).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Transaksi deleted successfully"})
}

func tambahDetailTransaksi(c *gin.Context) {
	var input DetailTransaksi
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

// HapusDetailTransaksi menghapus data detail transaksi berdasarkan ID dari database
func hapusDetailTransaksi(c *gin.Context) {
	id := c.Param("id")
	var detail DetailTransaksi
	if err := db.First(&detail, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Detail transaksi not found"})
		return
	}

	if err := db.Delete(&detail).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Detail transaksi deleted successfully"})
}

func getDetailTransaksiByID(c *gin.Context) {
	id := c.Param("id")
	var detail DetailTransaksi
	if err := db.First(&detail, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Detail transaksi not found"})
		return
	}
	c.JSON(http.StatusOK, detail)
}

func tambahTransaksi(c *gin.Context) {
	var input Transaksi
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Membuat transaksi baru di database
	if err := db.Create(&input).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Mendapatkan detail transaksi dari request
	var detailInput []DetailTransaksi
	if err := c.ShouldBindJSON(&detailInput); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Memproses setiap detail transaksi
	for _, detail := range detailInput {
		// Mengurangkan stok barang berdasarkan jumlah transaksi
		if err := kurangiStokBarang(detail.IDBarang, detail.Jumlah); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// Mengaitkan detail transaksi dengan transaksi utama
		detail.IDTransaksi = input.ID
		if err := db.Create(&detail).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "Transaksi berhasil ditambahkan"})
}

// fungsi ini digunakan untuk mengurangkan stok barang berdasarkan ID barang dan jumlah yang dibeli
func kurangiStokBarang(idBarang uint, jumlah int) error {
	var barang Barang
	if err := db.First(&barang, idBarang).Error; err != nil {
		return err
	}

	// Memastikan stok cukup untuk transaksi
	if barang.Stok < jumlah {
		return errors.New("stok barang tidak cukup")
	}

	// Mengurangkan stok barang
	barang.Stok -= jumlah

	// Menyimpan perubahan ke database
	if err := db.Save(&barang).Error; err != nil {
		return err
	}

	return nil
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

	// r.GET("/transaksi", getTransaksi)
	// r.GET("/transaksi/:id", getTransaksiByID)
	// r.POST("/transaksi", tambahTransaksi)
	// r.PUT("/transaksi/:id", perbaruiTransaksi)
	// r.DELETE("/transaksi/:id", hapusTransaksi)

	// Endpoint untuk menambah transaksi
	r.POST("/transaksi", tambahTransaksi)

	// Endpoint untuk mendapatkan semua transaksi
	r.GET("/transaksi", getTransaksi)

	// Endpoint untuk mendapatkan transaksi berdasarkan ID
	r.GET("/transaksi/:id", getTransaksiByID)

	// Endpoint untuk mengupdate transaksi berdasarkan ID
	r.PUT("/transaksi/:id", perbaruiTransaksi)

	// Endpoint untuk menghapus transaksi berdasarkan ID
	r.DELETE("/transaksi/:id", hapusTransaksi)

	r.GET("/detail-transaksi/:id", getDetailTransaksiByID)
	r.POST("/detail-transaksi", tambahDetailTransaksi)
	r.DELETE("/detail-transaksi/:id", hapusDetailTransaksi)
	// Menjalankan server
	r.Run(":8080")
}
