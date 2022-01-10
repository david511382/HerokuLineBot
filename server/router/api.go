package router

import (
	clubLogicDomain "heroku-line-bot/logic/club/domain"
	indexApi "heroku-line-bot/server/api"
	badmintonApi "heroku-line-bot/server/api/badminton"
	"heroku-line-bot/server/common"
	"heroku-line-bot/server/middleware"
	"heroku-line-bot/server/validation"

	"github.com/gin-gonic/gin"
)

func SetupApiRouter(router *gin.Engine) *gin.Engine {
	// 客製參數驗證
	validation.RegisterValidation()

	lineTokenVerifier := common.NewLineTokenVerifier()

	// api
	api := router.Group("/api")
	api.Use(middleware.AuthorizeToken(lineTokenVerifier, false))
	// api auth
	apiAuth := api.Group("/")
	apiAuth.Use(middleware.AuthorizeToken(lineTokenVerifier, true))
	apiAuth.GET("/user-info", indexApi.GetUserInfo)

	// api/badminton
	apiBadminton := api.Group("/badminton")
	apiBadminton.GET("/rental-courts", badmintonApi.GetRentalCourts)
	// api/badminton auth
	apiBadminton.Use(middleware.AuthorizeToken(lineTokenVerifier, true))
	apiBadminton.Use(middleware.VerifyAuthorize(map[int16]bool{
		int16(clubLogicDomain.ADMIN_CLUB_ROLE): true,
		int16(clubLogicDomain.CADRE_CLUB_ROLE): true,
	}))
	apiBadminton.POST("/rental-courts", badmintonApi.AddRentalCourt)

	// TODO wait for delete
	// api/club
	apiClub := api.Group("/club")
	apiClub.GET("/rental-courts", badmintonApi.GetRentalCourts)

	return router
}
