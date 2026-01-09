package message

import (
	"net/http"
	"github.com/gin-gonic/gin"
)

func Send(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "Message sent"})
}

func List(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"code": 0, "data": []interface{}{}})
}

func Templates(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"code": 0, "data": []interface{}{}})
}

func UnreadCount(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"code": 0, "data": gin.H{"count": 0}})
}
