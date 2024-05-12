package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/hifisultonauliya/goGinCore/src/helper"
	"github.com/hifisultonauliya/goGinCore/src/models"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var (
	// test user using my email hifisultonauliya@gmail.com
	oauthConfig = oauth2.Config{
		ClientID:     "test",
		ClientSecret: "test2",
		RedirectURL:  "http://localhost:8080/auth/google/callback",
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email"},
		Endpoint:     google.Endpoint,
	}
)

type AuthController struct{}

func NewAuthController() *AuthController {
	return &AuthController{}
}

func (ac *AuthController) Login(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	db := c.MustGet("db").(*gorm.DB)
	if !user.Validate(db) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	userid := strconv.FormatUint(uint64(user.ID), 10)
	token, err := helper.GenerateToken(userid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}

func (ac *AuthController) GoogleLogin(c *gin.Context) {
	url := oauthConfig.AuthCodeURL("state")
	c.Redirect(http.StatusFound, url)
}

func (ac *AuthController) GoogleCallback(c *gin.Context) {
	code := c.Query("code")
	token, err := oauthConfig.Exchange(c, code)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to exchange token"})
		return
	}

	// Use the token to get user info
	userInfo, err := getUserInfo(token)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user info"})
		return
	}

	// Here you can authenticate the user and create a session
	c.JSON(http.StatusOK, userInfo)

	// get email and generate token.... was missing here...
	email := userInfo["email"].(string)
	var result models.User
	db := c.MustGet("db").(*gorm.DB)
	if err := db.Where("email = ?", email).First(&result).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user data"})
		return
	}

	userid := strconv.FormatUint(uint64(result.ID), 10)
	token, err = helper.GenerateToken(userid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}

func getUserInfo(token *oauth2.Token) (map[string]interface{}, error) {
	client := oauthConfig.Client(oauth2.NoContext, token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var userInfo map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		return nil, err
	}

	return userInfo, nil
}
