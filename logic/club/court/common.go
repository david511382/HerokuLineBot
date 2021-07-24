package court

import (
	"heroku-line-bot/logic/club/court/domain"
	dbLogicDomain "heroku-line-bot/logic/database/domain"
)

func ReasonMessage(cancelReason dbLogicDomain.ReasonType) string {
	switch cancelReason {
	case domain.CANCEL_REASON_TYPE:
		return "取消"
	default:
		return ""
	}
}

func GetStatus(isPay, isRefund, isCancel bool) domain.RentalCourtsStatus {
	if isCancel {
		if isRefund {
			return domain.RENTAL_COURTS_STATUS_CANCEL
		}
		return domain.RENTAL_COURTS_STATUS_NOT_REFUND
	}

	if !isPay {
		return domain.RENTAL_COURTS_STATUS_NOT_PAY
	}
	return domain.RENTAL_COURTS_STATUS_OK
}
