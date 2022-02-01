package docs

import (
	"net/http"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger" // gin-swagger middleware

	_ "heroku-line-bot/docs" // 務必引入 swag init 創建的目錄
)

func Swag(c *gin.Context) {
	// 套件設定的服務頁面為index.html
	// 跳轉到該路徑
	if c.Request.RequestURI == "/docs" || c.Request.RequestURI == "/docs/" {
		c.Redirect(http.StatusTemporaryRedirect, "/docs/index.html")
	} else {
		// 當有名為Rlease的環境變數設置的時候，頁面將顯示404
		ginSwagger.DisablingWrapHandler(swaggerFiles.Handler, "Release")(c)
	}
}
