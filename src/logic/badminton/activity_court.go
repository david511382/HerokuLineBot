package badminton

import (
	"fmt"
	"heroku-line-bot/src/pkg/util"
	errUtil "heroku-line-bot/src/pkg/util/error"
	"strings"
	"time"
)

type ActivityCourt struct {
	FromTime     time.Time `json:"from_time"`
	ToTime       time.Time `json:"to_time"`
	Count        uint8     `json:"count"`
	PricePerHour float64   `json:"price_per_hour"`
}

func (b *ActivityCourt) Cost() util.Float {
	return b.TotalHours().MulFloat(b.PricePerHour)
}

func (b *ActivityCourt) Hours() util.Float {
	return util.NewFloat(b.ToTime.Sub(b.FromTime).Hours())
}

func (b *ActivityCourt) TotalHours() util.Float {
	return b.Hours().MulFloat(float64(b.Count))
}

func (b *ActivityCourt) Time() string {
	return fmt.Sprintf(
		"%s~%s",
		b.FromTime.Format(util.TIME_HOUR_MIN_FORMAT),
		b.ToTime.Format(util.TIME_HOUR_MIN_FORMAT),
	)
}

type DbActivityCourtsStr string

func (courtsStr DbActivityCourtsStr) String() string {
	return string(courtsStr)
}

func (courtsStr DbActivityCourtsStr) ParseCourts() (
	respCourts ActivityCourts,
	resultErrInfo errUtil.IError,
) {
	courtsStrs := strings.Split(courtsStr.String(), ",")
	for _, courtsStr := range courtsStrs {
		court := &ActivityCourt{}
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

type ActivityCourts []*ActivityCourt

func (courts ActivityCourts) FormatDbCourts() DbActivityCourtsStr {
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
	return DbActivityCourtsStr(strings.Join(courtStrs, ","))
}
