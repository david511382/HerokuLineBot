package common

import (
	errLogic "heroku-line-bot/logic/error"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, data)
}

func FailRequest(c *gin.Context, errInfo errLogic.IError) {
	Abort(c, http.StatusBadRequest, errInfo)
}

func FailAuth(c *gin.Context, errInfo errLogic.IError) {
	Abort(c, http.StatusUnauthorized, errInfo)
}

func FailForbidden(c *gin.Context, errInfo errLogic.IError) {
	Abort(c, http.StatusForbidden, errInfo)
}

func FailInternal(c *gin.Context, errInfo errLogic.IError) {
	Abort(c, http.StatusInternalServerError, errInfo)
}

func Abort(c *gin.Context, code int, errInfo errLogic.IError) {
	if errInfo == nil {
		c.AbortWithStatus(code)
	} else {
		c.AbortWithError(code, errInfo)
	}
}
