package router

import (
	"heroku-line-bot/bootstrap"
	clubLogicDomain "heroku-line-bot/src/logic/club/domain"
	badmintonApi "heroku-line-bot/src/server/api/badminton"
	"heroku-line-bot/src/server/domain"
	"heroku-line-bot/src/server/middleware"

	"github.com/gin-gonic/gin"
)

func SetupManagerRouter(cfg *bootstrap.Config, tokenVerifier domain.ITokenVerifier, api *gin.RouterGroup) {
	// api/badminton auth
	apiBadminton := api.Group("/badminton")
	apiBadminton.Use(middleware.AuthorizeToken(tokenVerifier, true))
	apiBadminton.Use(middleware.VerifyAuthorize(map[clubLogicDomain.ClubRole]bool{
		clubLogicDomain.ADMIN_CLUB_ROLE: true,
		clubLogicDomain.CADRE_CLUB_ROLE: true,
	}))
	apiBadminton.GET("/rental-courts", badmintonApi.GetRentalCourts)
	apiBadminton.POST("/rental-courts", badmintonApi.AddRentalCourt)
}

func SetupWsManagerRouter(webSocket *gin.RouterGroup) {
}
