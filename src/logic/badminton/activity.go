package badminton

import (
	"fmt"
	"heroku-line-bot/src/logic/badminton/domain"
	"heroku-line-bot/src/pkg/util"
	errUtil "heroku-line-bot/src/pkg/util/error"
	"strings"
	"time"
)

func ParseActivityDbCourts(
	courtsStr string,
) (
	respCourts []*domain.ActivityCourt,
	resultErrInfo errUtil.IError,
) {
	courtsStrs := strings.Split(courtsStr, ",")
	for _, courtsStr := range courtsStrs {
		court := &domain.ActivityCourt{}
		timeStr := ""
		if _, err := fmt.Sscanf(
			courtsStr,
			"%d-%f-%s",
			&court.Count,
			&court.PricePerHour,
			&timeStr); err != nil {
			resultErrInfo = errUtil.NewError(err)
			return
		}
		times := strings.Split(timeStr, "~")
		if len(times) != 2 {
			errInfo := errUtil.New("時間格式錯誤")
			resultErrInfo = errUtil.Append(resultErrInfo, errInfo)
			return
		}
		fromTimeStr := times[0]
		toTimeStr := times[1]
		if t, err := time.Parse(util.TIME_HOUR_MIN_FORMAT, fromTimeStr); err != nil {
			resultErrInfo = errUtil.NewError(err)
			return
		} else {
			court.FromTime = t
		}
		if t, err := time.Parse(util.TIME_HOUR_MIN_FORMAT, toTimeStr); err != nil {
			resultErrInfo = errUtil.NewError(err)
			return
		} else {
			court.ToTime = t
		}

		respCourts = append(respCourts, court)
	}
	return
}

func FormatActivityDbCourts(
	courts []*domain.ActivityCourt,
) string {
	courtStrs := []string{}
	for _, court := range courts {
		courtStr := fmt.Sprintf(
			"%d-%.1f-%s~%s",
			court.Count,
			court.PricePerHour,
			court.FromTime.Format(util.TIME_HOUR_MIN_FORMAT),
			court.ToTime.Format(util.TIME_HOUR_MIN_FORMAT),
		)
		courtStrs = append(courtStrs, courtStr)
	}
	return strings.Join(courtStrs, ",")
}
