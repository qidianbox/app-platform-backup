package file

import (
	"net/http"
	"github.com/gin-gonic/gin"
)

func Upload(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "File uploaded", "data": gin.H{"url": ""}})
}

func List(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"code": 0, "data": []interface{}{}})
}

func Delete(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "File deleted"})
}

func Stats(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"code": 0, "data": gin.H{"total": 0, "size": 0}})
}
