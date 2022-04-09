package middleware

import (
	"bytes"
	"encoding/json"
	"heroku-line-bot/src/logger"
	errUtil "heroku-line-bot/src/pkg/util/error"
	"heroku-line-bot/src/server/common"
	"heroku-line-bot/src/server/domain"
	"io"
	"io/ioutil"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

func Logger() gin.HandlerFunc {
	notlogged := []string{
		"/",
	}

	var skip map[string]struct{}
	if length := len(notlogged); length > 0 {
		skip = make(map[string]struct{}, length)

		for _, path := range notlogged {
			skip[path] = struct{}{}
		}
	}

	return func(c *gin.Context) {
		// Start timer
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery
		resultErrInfo := errUtil.New(
			"API Log",
			zerolog.InfoLevel,
		)
		if bs, err := ioutil.ReadAll(c.Request.Body); err == nil {
			resultErrInfo.Attr("Body", string(bs))
			c.Request.Body = io.NopCloser(bytes.NewReader(bs))
		}

		// Process request
		c.Next()

		// Stop timer
		nowTime := time.Now()

		if _, ok := skip[path]; ok {
			return
		}

		// Log only when path is not being skipped

		resultErrInfo.Attr("ClientIP", c.ClientIP())
		resultErrInfo.Attr("Method", c.Request.Method)
		if raw != "" {
			path = path + "?" + raw
		}
		resultErrInfo.Attr("Path", path)
		resultErrInfo.Attr("Proto", c.Request.Proto)
		resultErrInfo.Attr("Status", c.Writer.Status())
		resultErrInfo.Attr("Duration", nowTime.Sub(start).String())
		resultErrInfo.Attr("UserAgent", c.Request.UserAgent())
		if responseValue, isExist := c.Get(domain.KEY_RESPONSE_CONTEXT); isExist {
			resultErrInfo.Attr("Response", responseValue)
		}
		if claims := common.GetClaims(c); claims != nil {
			bs, err := json.Marshal(claims)
			if err == nil {
				resultErrInfo.Attr("Claims", string(bs))
			}
		}

		logger.LogError(common.GetLogName(c), resultErrInfo)
	}
}
