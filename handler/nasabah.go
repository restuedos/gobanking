package handler

import (
	"fmt"
	"gobanking/config"
	"gobanking/model"
	"math/rand"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type NasabahHandler struct {
	db       *gorm.DB
	cfg      *config.Config
	validate *validator.Validate
}

func NewNasabahHandler(db *gorm.DB, cfg *config.Config) *NasabahHandler {
	return &NasabahHandler{
		db:       db,
		cfg:      cfg,
		validate: validator.New(),
	}
}

func generateNoRekening() string {
	return fmt.Sprintf("%010d", rand.Intn(9000000000)+1000000000)
}

func (h *NasabahHandler) Daftar(c echo.Context) error {
	var req model.DaftarRequest
	if err := c.Bind(&req); err != nil {
		h.cfg.Logger.Warn("format request salah", "error", err)
		return c.JSON(http.StatusBadRequest, model.ErrorResponse{Remark: "Format request salah"})
	}

	if err := h.validate.Struct(req); err != nil {
		h.cfg.Logger.Warn("gagal validasi", "error", err)
		return c.JSON(http.StatusBadRequest, model.ErrorResponse{Remark: "Semua field harus diisi"})
	}

	var existingNasabah model.Nasabah
	result := h.db.Where("nik = ? OR no_hp = ?", req.NIK, req.NoHP).First(&existingNasabah)
	if result.RowsAffected > 0 {
		h.cfg.Logger.Info("gagal daftar: duplikat NIK atau No Handphone",
			"nik", req.NIK,
			"no_hp", req.NoHP,
		)
		return c.JSON(http.StatusBadRequest, model.ErrorResponse{Remark: "NIK atau No Handphone sudah terdaftar"})
	}

	nasabah := model.Nasabah{
		Nama:       req.Nama,
		NIK:        req.NIK,
		NoHP:       req.NoHP,
		NoRekening: generateNoRekening(),
		Saldo:      0,
	}

	if err := h.db.Create(&nasabah).Error; err != nil {
		h.cfg.Logger.Error("gagal mendaftarkan nasabah", "error", err)
		return c.JSON(http.StatusInternalServerError, model.ErrorResponse{Remark: "Internal server error"})
	}

	h.cfg.Logger.Info("nasabah berhasil didaftarkan",
		"name", nasabah.Nama,
		"no_rekening", nasabah.NoRekening,
	)

	return c.JSON(http.StatusOK, model.RekeningResponse{NoRekening: nasabah.NoRekening})
}

func (h *NasabahHandler) Tabung(c echo.Context) error {
	var req model.TransaksiRequest
	if err := c.Bind(&req); err != nil {
		h.cfg.Logger.Warn("format request salah", "error", err)
		return c.JSON(http.StatusBadRequest, model.ErrorResponse{Remark: "Format request salah"})
	}

	if err := h.validate.Struct(req); err != nil {
		h.cfg.Logger.Warn("gagal validasi", "error", err)
		return c.JSON(http.StatusBadRequest, model.ErrorResponse{Remark: "Semua field harus diisi"})
	}

	var nasabah model.Nasabah
	if err := h.db.Where("no_rekening = ?", req.NoRekening).First(&nasabah).Error; err != nil {
		h.cfg.Logger.Info("gagal tabungan: rekening tidak ditemukan",
			"no_rekening", req.NoRekening,
		)
		return c.JSON(http.StatusBadRequest, model.ErrorResponse{Remark: "No Rekening tidak ditemukan"})
	}

	nasabah.Saldo += req.Nominal
	if err := h.db.Save(&nasabah).Error; err != nil {
		h.cfg.Logger.Error("gagal memperbarui saldo", "error", err)
		return c.JSON(http.StatusInternalServerError, model.ErrorResponse{Remark: "Internal server error"})
	}

	h.cfg.Logger.Info("tabungan berhasil",
		"no_rekening", nasabah.NoRekening,
		"nominal", req.Nominal,
		"saldo_baru", nasabah.Saldo,
	)

	return c.JSON(http.StatusOK, model.SaldoResponse{Saldo: nasabah.Saldo})
}

func (h *NasabahHandler) Tarik(c echo.Context) error {
	var req model.TransaksiRequest
	if err := c.Bind(&req); err != nil {
		h.cfg.Logger.Warn("format request salah", "error", err)
		return c.JSON(http.StatusBadRequest, model.ErrorResponse{Remark: "Format request salah"})
	}

	if err := h.validate.Struct(req); err != nil {
		h.cfg.Logger.Warn("gagal validasi", "error", err)
		return c.JSON(http.StatusBadRequest, model.ErrorResponse{Remark: "Semua field harus diisi"})
	}

	var nasabah model.Nasabah
	if err := h.db.Where("no_rekening = ?", req.NoRekening).First(&nasabah).Error; err != nil {
		h.cfg.Logger.Info("gagal penarikan: rekening tidak ditemukan",
			"no_rekening", req.NoRekening,
		)
		return c.JSON(http.StatusBadRequest, model.ErrorResponse{Remark: "No Rekening tidak ditemukan"})
	}

	if nasabah.Saldo < req.Nominal {
		h.cfg.Logger.Info("gagal penarikan: saldo tidak mencukupi",
			"no_rekening", nasabah.NoRekening,
			"saldo", nasabah.Saldo,
			"nominal_ditarik", req.Nominal,
		)
		return c.JSON(http.StatusBadRequest, model.ErrorResponse{Remark: "Saldo tidak mencukupi"})
	}

	nasabah.Saldo -= req.Nominal
	if err := h.db.Save(&nasabah).Error; err != nil {
		h.cfg.Logger.Error("gagal memperbarui saldo", "error", err)
		return c.JSON(http.StatusInternalServerError, model.ErrorResponse{Remark: "Internal server error"})
	}

	h.cfg.Logger.Info("penarikan berhasil",
		"no_rekening", nasabah.NoRekening,
		"nominal", req.Nominal,
		"saldo_baru", nasabah.Saldo,
	)

	return c.JSON(http.StatusOK, model.SaldoResponse{Saldo: nasabah.Saldo})
}

func (h *NasabahHandler) Saldo(c echo.Context) error {
	noRekening := c.Param("no_rekening")

	var nasabah model.Nasabah
	if err := h.db.Where("no_rekening = ?", noRekening).First(&nasabah).Error; err != nil {
		h.cfg.Logger.Info("gagal pengecekan saldo: rekening tidak ditemukan",
			"no_rekening", noRekening,
		)
		return c.JSON(http.StatusBadRequest, model.ErrorResponse{Remark: "No Rekening tidak ditemukan"})
	}

	h.cfg.Logger.Info("pengecekan saldo berhasil",
		"no_rekening", nasabah.NoRekening,
		"saldo", nasabah.Saldo,
	)

	return c.JSON(http.StatusOK, model.SaldoResponse{Saldo: nasabah.Saldo})
}
