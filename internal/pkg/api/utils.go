package api

import "github.com/gin-gonic/gin"

func notFound(c *gin.Context) {
	respMessage(c, 404, "not found")
}

func respMessage(c *gin.Context, code uint, message string) {
	c.JSON(int(code), gin.H{"message": message})
}
