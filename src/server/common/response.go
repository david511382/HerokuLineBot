package common

import (
	"encoding/json"
	"heroku-line-bot/src/pkg/errorcode"
	errUtil "heroku-line-bot/src/pkg/util/error"
	"heroku-line-bot/src/server/domain"
	"heroku-line-bot/src/server/domain/resp"

	"github.com/gin-gonic/gin"
)

func Success(c *gin.Context, data any) {
	response(
		c,
		errorcode.ERROR_MSG_SUCCESS.New(nil),
		&resp.Base{
			Data: data,
		},
	)
}

func FailRequest(c *gin.Context, errInfo errUtil.IError) {
	FailErrorcode(c, errorcode.ERROR_MSG_REQUEST.New(errInfo))
}

func FailAuth(c *gin.Context, errInfo errUtil.IError) {
	FailErrorcode(c, errorcode.ERROR_MSG_AUTH.New(errInfo))
}

func FailForbidden(c *gin.Context, errInfo errUtil.IError) {
	FailErrorcode(c, errorcode.ERROR_MSG_FORBIDDEN.New(errInfo))
}

func FailInternal(c *gin.Context, errInfo errUtil.IError) {
	FailErrorcode(c, errorcode.ERROR_MSG_ERROR.New(errInfo))
}

// errInfo 不得為 nil
func Fail(c *gin.Context, errInfo errUtil.IError) {
	errCode := errorcode.GetErrorcode(errInfo)
	if errCode != nil {
		FailErrorcode(c, errCode)
	} else {
		FailInternal(c, errInfo)
	}
}

// errCode 不得為 nil
func FailErrorcode(c *gin.Context, errCode errorcode.IErrorcode) {
	response(c, errCode, nil)
}

// errCode 不得為 nil
func response(c *gin.Context, errCode errorcode.IErrorcode, result *resp.Base) {
	if result == nil {
		result = &resp.Base{}
	}
	result.Message = errCode.Error()
	if bs, err := json.Marshal(result); err == nil {
		c.Set(domain.KEY_RESPONSE_CONTEXT, string(bs))
	}

	code, errInfo := errCode.Log()
	if errInfo != nil {
		c.Set(domain.KEY_RESPONSE_ERROR, errInfo)
	}

	c.AbortWithStatusJSON(code, result)
}
