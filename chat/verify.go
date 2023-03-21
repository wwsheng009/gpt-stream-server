package chat

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func Verify(c *gin.Context) {
	var request struct {
		Token string `json:"token" binding:"required"`
	}

	if err := c.BindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "Fail", "message": "Bad request", "data": nil})
		return
	}

	if request.Token != os.Getenv("AUTH_SECRET_KEY") {
		// c.JSON(http.StatusUnauthorized, gin.H{"status": "Fail", "message": "密钥无效 | Secret key is invalid", "data": nil})
		// c.Errors = append(c.Errors)
		// c.AbortWithError(http.StatusUnauthorized, fmt.Errorf("密钥无效 | Secret key is invalid"))
		c.AbortWithStatusJSON(http.StatusOK, gin.H{"status": "Fail", "message": "密钥无效 | Secret key is invalid", "data": nil})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "Success", "message": "Verify successfully", "data": nil})
}
