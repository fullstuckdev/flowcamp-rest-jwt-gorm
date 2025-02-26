package controllers

import (
	"golangapi/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type RelationController struct {
	DB *gorm.DB
}

func NewRelationController(db *gorm.DB) *RelationController {
	return &RelationController{DB: db}
}

// One to One
func (rc *RelationController) CreateProfile(c *gin.Context) {
	var profile models.Profile
	
	if err := c.ShouldBindJSON(&profile); err != nil {
		c.JSON(400, gin.H{"error:": err.Error()})
	}

	if err := rc.DB.Create(&profile).Error; err != nil {
		c.JSON(500, gin.H{"error:": err.Error()})
		return
	}

	c.JSON(201, gin.H{"data": profile})
}

// One to Many
func (rc *RelationController) CreatePost(c *gin.Context) {
	var post models.Post

	if err := c.ShouldBindJSON(&post); err != nil {
		c.JSON(400, gin.H{"error:": err.Error()})
	}

	if err := rc.DB.Create(&post).Error; err != nil {
		c.JSON(500, gin.H{"error:": err.Error()})
		return
	}

	c.JSON(201, gin.H{"data": post})
}

// Create
func (rc *RelationController) CreateGroups(c *gin.Context) {
	var group models.Group
	
	if err := c.ShouldBindJSON(&group); err != nil {
		c.JSON(400, gin.H{"error:": err.Error()})
	}


	if err := rc.DB.Create(&group).Error; err != nil {
		c.JSON(500, gin.H{"error:": err.Error()})
		return
	}

	c.JSON(201, gin.H{"data": group})
}

// many to many
func (rc *RelationController) AddUserToGroup(c *gin.Context) {
	var request struct {
		UserID uint `json:"user_id"`
		GroupID uint `json:"group_id"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(400, gin.H{"error:": err.Error()})
	}

	var user models.User
	var group models.Group

	if err := rc.DB.First(&user, request.UserID).Error; err != nil {
		c.JSON(404, gin.H{"error": "user not found"})
	}

	if err := rc.DB.First(&group, request.GroupID).Error; err != nil {
		c.JSON(404, gin.H{"error": "group not found"})
	}

	if err := rc.DB.Model(&group).Association("Users").Append(&user); err != nil {
		c.JSON(400, gin.H{"error:": err.Error()})
		return
	}

	c.JSON(201, gin.H{"message": "berhasil menambahkan data!"})


}