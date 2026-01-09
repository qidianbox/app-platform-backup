package event

import (
	"net/http"
	"github.com/gin-gonic/gin"
)

func Report(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "Event reported"})
}

func List(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"code": 0, "data": []interface{}{}})
}

func Stats(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"code": 0, "data": gin.H{}})
}

func Funnel(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"code": 0, "data": gin.H{}})
}

func Definitions(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"code": 0, "data": []interface{}{}})
}
