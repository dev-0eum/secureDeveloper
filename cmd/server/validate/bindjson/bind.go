package bindjson

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func JSONValidate(c *gin.Context, request any) {
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid update request"})
		return
	}
}
