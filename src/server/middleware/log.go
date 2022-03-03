package middleware

import (
	"encoding/json"
	"heroku-line-bot/src/logger"
	"heroku-line-bot/src/server/common"
	"heroku-line-bot/src/server/domain"
	"heroku-line-bot/src/util"
	errUtil "heroku-line-bot/src/util/error"
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

		// Process request
		c.Next()

		// Stop timer
		nowTime := time.Now()

		if _, ok := skip[path]; ok {
			return
		}

		// Log only when path is not being skipped

		resultErrInfo := errUtil.New(
			"API Log",
			zerolog.InfoLevel,
		)

		if bs, err := ioutil.ReadAll(c.Request.Body); err == nil {
			resultErrInfo.Attr("Body", string(bs))
		}
		resultErrInfo.Attr("ClientIP", c.ClientIP())
		resultErrInfo.Attr("Time", nowTime.Format(util.DATE_TIME_FORMAT))
		resultErrInfo.Attr("Method", c.Request.Method)
		if raw != "" {
			path = path + "?" + raw
		}
		resultErrInfo.Attr("Path", path)
		resultErrInfo.Attr("Proto", c.Request.Proto)
		resultErrInfo.Attr("Status", c.Writer.Status())
		resultErrInfo.Attr("Duration", nowTime.Sub(start))
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

		logger.Log(common.GetLogName(c), resultErrInfo)
	}
}
