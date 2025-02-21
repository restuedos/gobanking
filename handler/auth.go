package handler

import (
	"gobanking/config"
	"gobanking/model"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthHandler struct {
	db       *gorm.DB
	cfg      *config.Config
}

func NewAuthHandler(db *gorm.DB, cfg *config.Config) *AuthHandler {
	return &AuthHandler{
		db:  db,
		cfg: cfg,
	}
}

// @Summary Register new user
// @Description Register a new user with email and password
// @Tags auth
// @Accept json
// @Produce json
// @Param user body model.RegisterRequest true "User registration details"
// @Success 201 {object} model.TokenResponse
// @Failure 400 {object} model.ErrorResponse
// @Router /register [post]
func (h *AuthHandler) Register(c echo.Context) error {
	var req model.RegisterRequest
	if err := c.Bind(&req); err != nil {
		h.cfg.Logger.Warn("format request salah", "error", err)
		return c.JSON(http.StatusBadRequest, model.ErrorResponse{Remark: "Format request salah"})
	}

	// Check if email already exists
	var existingUser model.User
	result := h.db.Where("email = ?", req.Email).First(&existingUser)
	if result.RowsAffected > 0 {
		h.cfg.Logger.Info("gagal register: email sudah terdaftar", "email", req.Email)
		return c.JSON(http.StatusBadRequest, model.ErrorResponse{Remark: "Email sudah terdaftar"})
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		h.cfg.Logger.Error("gagal hash password", "error", err)
		return c.JSON(http.StatusInternalServerError, model.ErrorResponse{Remark: "Internal server error"})
	}

	user := model.User{
		Email:    req.Email,
		Password: string(hashedPassword),
	}

	if err := h.db.Create(&user).Error; err != nil {
		h.cfg.Logger.Error("gagal mendaftarkan user", "error", err)
		return c.JSON(http.StatusInternalServerError, model.ErrorResponse{Remark: "Internal server error"})
	}

	// Generate JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"email":   user.Email,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, err := token.SignedString([]byte(h.cfg.JWT.Secret))
	if err != nil {
		h.cfg.Logger.Error("gagal generate token", "error", err)
		return c.JSON(http.StatusInternalServerError, model.ErrorResponse{Remark: "Internal server error"})
	}

	h.cfg.Logger.Info("user berhasil didaftarkan", "email", user.Email)
	return c.JSON(http.StatusCreated, model.TokenResponse{Token: tokenString})
}

// @Summary Login user
// @Description Login with email and password
// @Tags auth
// @Accept json
// @Produce json
// @Param user body model.LoginRequest true "User login details"
// @Success 200 {object} model.TokenResponse
// @Failure 400 {object} model.ErrorResponse
// @Router /login [post]
func (h *AuthHandler) Login(c echo.Context) error {
	var req model.LoginRequest
	if err := c.Bind(&req); err != nil {
		h.cfg.Logger.Warn("format request salah", "error", err)
		return c.JSON(http.StatusBadRequest, model.ErrorResponse{Remark: "Format request salah"})
	}

	var user model.User
	if err := h.db.Where("email = ?", req.Email).First(&user).Error; err != nil {
		h.cfg.Logger.Info("gagal login: email tidak ditemukan", "email", req.Email)
		return c.JSON(http.StatusBadRequest, model.ErrorResponse{Remark: "Email atau password salah"})
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		h.cfg.Logger.Info("gagal login: password salah", "email", req.Email)
		return c.JSON(http.StatusBadRequest, model.ErrorResponse{Remark: "Email atau password salah"})
	}

	// Generate JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"email":   user.Email,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, err := token.SignedString([]byte(h.cfg.JWT.Secret))
	if err != nil {
		h.cfg.Logger.Error("gagal generate token", "error", err)
		return c.JSON(http.StatusInternalServerError, model.ErrorResponse{Remark: "Internal server error"})
	}

	h.cfg.Logger.Info("user berhasil login", "email", user.Email)
	return c.JSON(http.StatusOK, model.TokenResponse{Token: tokenString})
}