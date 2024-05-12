package controllers

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/hifisultonauliya/goGinCore/src/helper"
	"github.com/hifisultonauliya/goGinCore/src/models"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

type TaskController struct{}

func NewTaskController() *TaskController {
	return &TaskController{}
}

func (ic *TaskController) GetTasks(c *gin.Context) {
	var tasks []models.Task
	db := c.MustGet("db").(*gorm.DB)
	db.Find(&tasks)
	c.JSON(http.StatusOK, tasks)
}

func (ic *TaskController) GetTask(c *gin.Context) {
	id := c.Param("id")
	var task models.Task
	db := c.MustGet("db").(*gorm.DB)
	if err := db.Where(idQuery, id).First(&task).Error; err != nil {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}
	c.JSON(http.StatusOK, task)
}

func (ic *TaskController) CreateTask(c *gin.Context) {
	var task models.Task
	if err := c.ShouldBindJSON(&task); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// get userid from auth jwt
	tokenString := c.GetHeader("Authorization")
	log.Println(">>> tokenString", tokenString)
	UserID, err := helper.GetUserIDFromJWT(tokenString)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	log.Println(">>> UserID", UserID)

	u64, err := strconv.ParseUint(UserID, 10, 32)
	log.Println(">>> u64", u64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	task.UserID = uint(u64)
	log.Println(">>>", task.UserID, task)

	db := c.MustGet("db").(*gorm.DB)
	db.Create(&task)
	c.JSON(http.StatusCreated, task)
}

func (ic *TaskController) UpdateTask(c *gin.Context) {
	id := c.Param("id")
	var task models.Task
	db := c.MustGet("db").(*gorm.DB)
	if err := db.Where(idQuery, id).First(&task).Error; err != nil {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}
	if err := c.ShouldBindJSON(&task); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// get userid from auth jwt
	tokenString := c.GetHeader("Authorization")
	UserID, err := helper.GetUserIDFromJWT(tokenString)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	u64, err := strconv.ParseUint(UserID, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	task.UserID = uint(u64)
	log.Println(">>>", task.UserID, task)

	db.Save(&task)
	c.JSON(http.StatusOK, task)
}

func (ic *TaskController) DeleteTask(c *gin.Context) {
	id := c.Param("id")
	var task models.Task
	db := c.MustGet("db").(*gorm.DB)
	if err := db.Where(idQuery, id).First(&task).Error; err != nil {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}
	db.Delete(&task)
	c.JSON(http.StatusOK, gin.H{"message": "Task deleted successfully"})
}
