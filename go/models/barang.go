package models

// Barang model
type Barang struct {
	ID    uint    `gorm:"primary_key" json:"id"`
	Nama  string  `json:"nama"`
	Harga float64 `json:"harga"`
	Stok  int     `json:"stok"`
	// tambahkan field lain sesuai kebutuhan
}
