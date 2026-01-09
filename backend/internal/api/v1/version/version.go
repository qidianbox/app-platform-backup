package version

import (
	"net/http"
	"github.com/gin-gonic/gin"
)

func List(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"code": 0, "data": []interface{}{}})
}

func Create(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "Version created"})
}

func Publish(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "Version published"})
}

func Offline(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "Version offline"})
}

func CheckUpdate(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"code": 0, "data": gin.H{"has_update": false}})
}
