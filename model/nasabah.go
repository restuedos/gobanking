package model

import "gorm.io/gorm"

type Nasabah struct {
	gorm.Model
	Nama       string  `gorm:"not null" json:"nama"`
	NIK        string  `gorm:"unique;not null" json:"nik"`
	NoHP       string  `gorm:"unique;not null" json:"no_hp"`
	NoRekening string  `gorm:"unique;not null" json:"-"`
	Saldo      float64 `gorm:"default:0;check:saldo >= 0" json:"-"`
}

type DaftarRequest struct {
	Nama string `json:"nama" validate:"required"`
	NIK  string `json:"nik" validate:"required"`
	NoHP string `json:"no_hp" validate:"required"`
}

type TransaksiRequest struct {
	NoRekening string  `json:"no_rekening" validate:"required"`
	Nominal    float64 `json:"nominal" validate:"required,gt=0"`
}

type ErrorResponse struct {
	Remark string `json:"remark"`
}

type SaldoResponse struct {
	Saldo float64 `json:"saldo"`
}

type RekeningResponse struct {
	NoRekening string `json:"no_rekening"`
}
