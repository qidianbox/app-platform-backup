package config

import (
	"net/http"
	"github.com/gin-gonic/gin"
)

func List(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"code": 0, "data": []interface{}{}})
}

func Create(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "Config created"})
}

func Update(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "Config updated"})
}

func Publish(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "Config published"})
}

func History(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"code": 0, "data": []interface{}{}})
}
