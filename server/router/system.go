package router

import (
	// TODO verify
	//clubLogicDomain"heroku-line-bot/logic/club/domain"
	indexApi "heroku-line-bot/server/api"
	badmintonApi "heroku-line-bot/server/api/badminton"
	clubApi "heroku-line-bot/server/api/club"
	clubLineBotApi "heroku-line-bot/server/api/clublinebot"
	"heroku-line-bot/server/common"
	"heroku-line-bot/server/middleware"
	"heroku-line-bot/server/validation"
	docsView "heroku-line-bot/server/view/docs"
	"io"
	"os"

	"github.com/gin-gonic/gin"
)

func SystemRouter() *gin.Engine {
	// 取消打印文字顏色
	gin.DisableConsoleColor()
	// 使用打印文字顏色
	gin.ForceConsoleColor()

	// 設定輸出的物件(本地文字檔)
	f, _ := os.Create("gin.log")
	// 指定輸出的目標(本地文字檔、Console)
	gin.DefaultWriter = io.MultiWriter(f, os.Stdout)

	// 設定gin
	router := gin.New()

	router.Use(gin.Logger())
	router.Use(middleware.Cors)

	// 客製參數驗證
	validation.RegisterValidation()

	// docs
	doc := router.Group("/docs")
	doc.GET("/*any", docsView.Swag)

	// api
	api := router.Group("/api")
	lineAuth := api.Group("/")
	lineAuth.Use(middleware.GetTokenAuthorize(common.NewLineTokenVerifier(), true))
	lineAuth.GET("/user-info", indexApi.GetUserInfo)

	// api/badminton
	apiBadminton := api.Group("/badminton")
	apiBadminton.Use(middleware.VerifyAuthorize(map[int16]bool{
		// TODO verify
		//int16(clubLogicDomain.ADMIN_CLUB_ROLE): true,
		//int16(clubLogicDomain.CADRE_CLUB_ROLE): true,
	}))
	apiBadminton.POST("/rental-courts", badmintonApi.AddRentalCourt)

	// api/club
	linebotClub := api.Group("/club")
	linebotClub.GET("/rental-courts", clubApi.GetRentalCourts)

	clubLineBotEvent := router.Group("/")
	clubLineBotEvent.POST("/club-line-bot", clubLineBotApi.Index)

	return router
}
