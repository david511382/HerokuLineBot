package router

import (
	"heroku-line-bot/bootstrap"
	badmintonApi "heroku-line-bot/src/server/api/badminton"
	"heroku-line-bot/src/server/domain"
	badmintonWs "heroku-line-bot/src/server/ws/badminton"

	"github.com/gin-gonic/gin"
)

func SetupCustomerRouter(cfg *bootstrap.Config, tokenVerifier domain.ITokenVerifier, api *gin.RouterGroup) {
	// api/badminton
	apiBadminton := api.Group("/badminton")
	apiBadminton.GET("/activitys", badmintonApi.GetActivitys)

	// TODO wait for delete
	// api/club
	apiClub := api.Group("/club")
	apiClub.GET("/rental-courts", badmintonApi.GetRentalCourts)
}

func SetupWsCustomerRouter(webSocket *gin.RouterGroup) {
	// ws/badminton
	wsBadminton := webSocket.Group("/badminton")
	wsBadminton.GET("/activitys", badmintonWs.GetActivitys)
}
