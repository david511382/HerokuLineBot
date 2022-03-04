package badminton

import (
	"encoding/json"
	"heroku-line-bot/bootstrap"
	clubLogicDomain "heroku-line-bot/src/logic/club/domain"
	"heroku-line-bot/src/pkg/util"
	"heroku-line-bot/src/repo"
	"heroku-line-bot/src/server/common"
	"heroku-line-bot/src/server/domain"
	"heroku-line-bot/src/server/middleware"
	"heroku-line-bot/src/server/validation"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
)

var (
	testServer *util.TestGinServer
	location   *time.Location
)

func TestMain(m *testing.M) {
	if err := bootstrap.SetEnvWorkDir(bootstrap.DEFAULT_WORK_DIR); err != nil {
		panic(err)
	}
	if err := bootstrap.SetEnvConfig("local"); err != nil {
		panic(err)
	}

	loc, err := time.LoadLocation("Asia/Taipei")
	if err != nil {
		panic(err)
	}
	location = loc

	if errInfo := repo.Init(); errInfo != nil {
		panic(errInfo.ErrorWithTrace())
	}
	defer repo.Dispose()

	r := setupRouter()
	testServer = util.NewTestServer(r)
	testServer.SetRequest(func(req *http.Request) error {
		claims := domain.JwtClaims{
			RoleID: int16(clubLogicDomain.ADMIN_CLUB_ROLE),
		}
		bs, err := json.Marshal(claims)
		if err != nil {
			return err
		}

		req.Header.Set(domain.HeaderAuthorization, string(bs))
		return nil
	})

	exitVal := m.Run()

	os.Exit(exitVal)
}

func setupRouter() *gin.Engine {
	// 使用打印文字顏色
	gin.ForceConsoleColor()
	gin.SetMode(gin.DebugMode)

	r := gin.New()
	// Recovery middleware recovers from any panics and writes a 500 if there was one.
	r.Use(gin.Recovery())
	r.Use(gin.Logger())

	// 客製參數驗證
	validation.RegisterValidation()

	jsonTokenVerifier := common.NewJsonTokenVerifier()

	// api
	api := r.Group("/api")
	api.Use(middleware.AuthorizeToken(jsonTokenVerifier, false))

	// api/badminton
	apiBadminton := api.Group("/badminton")
	apiBadminton.GET("/activitys", GetActivitys)
	// api/badminton auth
	apiBadminton.Use(middleware.AuthorizeToken(jsonTokenVerifier, true))
	apiBadminton.Use(middleware.VerifyAuthorize(map[int16]bool{
		int16(clubLogicDomain.ADMIN_CLUB_ROLE): true,
		int16(clubLogicDomain.CADRE_CLUB_ROLE): true,
	}))
	apiBadminton.GET("/rental-courts", GetRentalCourts)
	apiBadminton.POST("/rental-courts", AddRentalCourt)

	return r
}
