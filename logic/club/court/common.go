package court

import (
	"heroku-line-bot/logic/club/court/domain"
	dbLogicDomain "heroku-line-bot/logic/database/domain"
	dbDomain "heroku-line-bot/storage/database/domain"
)

func ReasonMessage(cancelReason dbDomain.ReasonType) string {
	switch cancelReason {
	case dbLogicDomain.CANCEL_REASON_TYPE:
		return "取消"
	default:
		return ""
	}
}

func GetStatus(isPay, isRefund bool) domain.RentalCourtsStatus {
	if isRefund {
		if isPay {
			// 訂場已取消已退款
			return domain.RENTAL_COURTS_STATUS_CANCEL
		}
		// 訂場已取消還沒退款
		return domain.RENTAL_COURTS_STATUS_NOT_REFUND
	}

	if !isPay {
		// 訂場沒付款
		return domain.RENTAL_COURTS_STATUS_NOT_PAY
	}
	// 訂場已付款
	return domain.RENTAL_COURTS_STATUS_OK
}
