package models

// DetailTransaksi model
type DetailTransaksi struct {
	ID          uint    `gorm:"primary_key" json:"id"`
	IDTransaksi uint    `json:"id_transaksi"`
	IDBarang    uint    `json:"id_barang"`
	Jumlah      int     `json:"jumlah"`
	HargaSatuan float64 `json:"harga_satuan"`
	TotalHarga  float64 `json:"total_harga"`
	// tambahkan field lain sesuai kebutuhan
}
