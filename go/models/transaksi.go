package models

import "time"

// Transaksi model
type Transaksi struct {
	ID      uint      `gorm:"primary_key" json:"id"`
	Tanggal time.Time `json:"tanggal"`
	// tambahkan field lain sesuai kebutuhan
}
