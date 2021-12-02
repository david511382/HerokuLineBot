package validation

import (
	"heroku-line-bot/server/domain/reqs"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10" //需使用Gin使用的版本
)

func RegisterValidation() {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		// 結構驗證，執行綁定reqs.time 的驗證
		v.RegisterStructValidation(timeValidation,
			reqs.FromTo{},
			reqs.MustFromTo{},
			reqs.MustFromBefore{},
			reqs.FromBefore{},
			reqs.FromToDate{},
			reqs.MustFromToDate{},
			reqs.MustFromBeforeDate{},
			reqs.FromBeforeDate{},
			reqs.AddRentalCourt{},
		)
	}
}
