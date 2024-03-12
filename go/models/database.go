package models

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	// DB adalah instance untuk koneksi ke database
	DB *gorm.DB
)

// InisialisasiDatabase membuka koneksi ke database
func InisialisasiDatabase() {
	// Ganti dengan konfigurasi MySQL yang sesuai
	dsn := "root:@tcp(localhost:3306)/tokovan?parseTime=true"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("Failed to connect to database")
	}

	// Auto Migrate untuk membuat tabel sesuai dengan definisi model
	db.AutoMigrate(&Transaksi{}, &DetailTransaksi{}, &Barang{})

	// Assign instance DB
	DB = db
}
