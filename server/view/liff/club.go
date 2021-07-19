package liff

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Club(c *gin.Context) {
	c.HTML(http.StatusOK, "liffClub.html", nil)
}
