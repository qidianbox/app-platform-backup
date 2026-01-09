package user

import (
	"net/http"
	"github.com/gin-gonic/gin"
)

func List(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"code": 0, "data": []interface{}{}, "message": "User list"})
}

func Detail(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"code": 0, "data": gin.H{}, "message": "User detail"})
}

func UpdateStatus(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "User status updated"})
}

func Stats(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"code": 0, "data": gin.H{"total": 0, "active": 0}})
}
