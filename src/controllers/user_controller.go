package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hifisultonauliya/goGinCore/src/models"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"golang.org/x/crypto/bcrypt"
)

type UserController struct{}

const (
	idQuery = "id = ?"
)

func NewUserController() *UserController {
	return &UserController{}
}

func (ic *UserController) GetUsers(c *gin.Context) {
	var users []models.User
	db := c.MustGet("db").(*gorm.DB)
	db.Find(&users)
	c.JSON(http.StatusOK, users)
}

func (ic *UserController) GetUser(c *gin.Context) {
	id := c.Param("id")
	var user models.User
	db := c.MustGet("db").(*gorm.DB)
	if err := db.Where(idQuery, id).First(&user).Error; err != nil {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}
	c.JSON(http.StatusOK, user)
}

func (ic *UserController) CreateUser(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user.Password = string(hashedPassword)

	db := c.MustGet("db").(*gorm.DB)
	if err := db.Where("email = ?", user.Email).First(&user).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email already registered"})
		return
	}

	db.Create(&user)
	c.JSON(http.StatusCreated, user)
}

func (ic *UserController) UpdateUser(c *gin.Context) {
	id := c.Param("id")
	var user models.User
	db := c.MustGet("db").(*gorm.DB)

	if err := db.Where(idQuery, id).First(&user).Error; err != nil {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	db.Save(&user)
	c.JSON(http.StatusOK, user)
}

func (ic *UserController) DeleteUser(c *gin.Context) {
	id := c.Param("id")
	var user models.User

	db := c.MustGet("db").(*gorm.DB)
	if err := db.Where(idQuery, id).First(&user).Error; err != nil {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}
	db.Delete(&user)
	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}
