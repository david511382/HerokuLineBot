package common

import (
	"strings"

	"github.com/gin-gonic/gin"
)

func GetLogName(c *gin.Context) string {
	return strings.Replace(c.Request.URL.Path, "/", "-", -1)
}
