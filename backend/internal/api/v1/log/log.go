package log

import (
	"net/http"
	"github.com/gin-gonic/gin"
)

func System(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"code": 0, "data": []interface{}{}})
}

func Operation(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"code": 0, "data": []interface{}{}})
}

func Stats(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"code": 0, "data": gin.H{}})
}

func Clean(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "Logs cleaned"})
}
