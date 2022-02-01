package router

import (
	"heroku-line-bot/bootstrap"
	clubLogicDomain "heroku-line-bot/src/logic/club/domain"
	indexApi "heroku-line-bot/src/server/api"
	badmintonApi "heroku-line-bot/src/server/api/badminton"
	"heroku-line-bot/src/server/common"
	"heroku-line-bot/src/server/domain"
	"heroku-line-bot/src/server/middleware"
	"heroku-line-bot/src/server/validation"

	"github.com/gin-gonic/gin"
)

func SetupApiRouter(cfg *bootstrap.Config, router *gin.Engine) *gin.Engine {
	isDebug := cfg.Var.UseDebug

	// 客製參數驗證
	validation.RegisterValidation()

	lineTokenVerifier := common.NewLineTokenVerifier()
	var defaultTokenVerifier domain.ITokenVerifier = lineTokenVerifier
	if isDebug {
		defaultTokenVerifier = common.NewDebugTokenVerifier()
	}

	// api
	api := router.Group("/api")
	api.Use(middleware.AuthorizeToken(defaultTokenVerifier, false))
	// api auth
	apiAuth := api.Group("/")
	apiAuth.Use(middleware.AuthorizeToken(lineTokenVerifier, true))
	apiAuth.GET("/user-info", indexApi.GetUserInfo)

	// api/badminton
	apiBadminton := api.Group("/badminton")
	apiBadminton.GET("/activitys", badmintonApi.GetActivitys)
	// api/badminton auth
	apiBadminton.Use(middleware.AuthorizeToken(lineTokenVerifier, true))
	apiBadminton.Use(middleware.VerifyAuthorize(map[int16]bool{
		int16(clubLogicDomain.ADMIN_CLUB_ROLE): true,
		int16(clubLogicDomain.CADRE_CLUB_ROLE): true,
	}))
	apiBadminton.GET("/rental-courts", badmintonApi.GetRentalCourts)
	apiBadminton.POST("/rental-courts", badmintonApi.AddRentalCourt)

	// TODO wait for delete
	// api/club
	apiClub := api.Group("/club")
	apiClub.GET("/rental-courts", badmintonApi.GetRentalCourts)

	return router
}
