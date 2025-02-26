package controllers

import (
	"golangapi/models"
	"golangapi/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AuthController struct {
	DB *gorm.DB
}

func NewAuthController(db *gorm.DB) *AuthController {
	return &AuthController{DB: db}
}

func (ac *AuthController) Register(c *gin.Context) {
	var user models.User

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}	

	// gimana cara hash?
	if err := user.HashPassword(user.Password); err != nil {
		c.JSON(500, gin.H{"Error": "Gagal encrypt password!"})
		return
	}

	result := ac.DB.Create(&user)

	if result.Error != nil {
		c.JSON(400, gin.H{"Error": "Gagal membuat users!"})
		return
	}

	token, err := utils.GenerateToken(user.ID)

	if err != nil {
		c.JSON(400, gin.H{"Error": "Gagal generate token!"})
		return
	}

	c.JSON(201, gin.H{
		"Message": "User berhasil dibuat",
		"Token": token,
	})
}

func (ac *AuthController) Login(c *gin.Context) {
	var loginReq models.LoginRequest

	if err := c.ShouldBindJSON(&loginReq); err != nil {
		c.JSON(400, gin.H{"error:": err.Error()})
	}

	var user models.User

	if err := ac.DB.Where("email = ?", loginReq.Email).First(&user).Error; err != nil {
		c.JSON(401, gin.H{"Error": "Email salah!"})
		return
	}

	// check password function
	if err := user.CheckPassword(loginReq.Password); err != nil {
		c.JSON(401, gin.H{"Error": "Password salah!"})
		return
	}

	token, err := utils.GenerateToken(user.ID)
	if err != nil {
		c.JSON(500, gin.H{"Error": "Gagal generate token!"})
		return
	}

	c.JSON(200, gin.H{
		"Message": "Login succesfully!",
		"Token": token,
	})
}