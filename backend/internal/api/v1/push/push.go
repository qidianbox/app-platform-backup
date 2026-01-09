package push

import (
	"net/http"
	"github.com/gin-gonic/gin"
)

func Send(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "Push sent"})
}

func Tasks(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"code": 0, "data": []interface{}{}})
}

func Stats(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"code": 0, "data": gin.H{"total": 0, "success": 0}})
}

func Templates(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"code": 0, "data": []interface{}{}})
}
