package club

import (
	"fmt"
	"heroku-line-bot/logic/club/domain"
	commonLogic "heroku-line-bot/logic/common"
	commonLogicDomain "heroku-line-bot/logic/common/domain"
	"heroku-line-bot/service/linebot"
	linebotDomain "heroku-line-bot/service/linebot/domain"
	linebotModel "heroku-line-bot/service/linebot/domain/model"
	"heroku-line-bot/storage/database"
	dbReqs "heroku-line-bot/storage/database/domain/model/reqs"
	"heroku-line-bot/util"
	"math"
	"strconv"
)

type getActivities struct {
	context    domain.ICmdHandlerContext `json:"-"`
	activities []*newActivity
}

func (b *getActivities) Init(context domain.ICmdHandlerContext) error {
	*b = getActivities{
		context:    context,
		activities: make([]*newActivity, 0),
	}
	arg := dbReqs.Activity{
		IsComplete: util.GetBoolP(false),
	}
	if dbDatas, err := database.Club.Activity.DatePlaceCourtsSubsidyDescriptionPeopleLimit(arg); err != nil {
		return err
	} else {
		for _, v := range dbDatas {
			activity := &newActivity{
				context:     context,
				Date:        v.Date,
				Place:       v.Place,
				Description: v.Description,
				PeopleLimit: v.PeopleLimit,
				ClubSubsidy: v.ClubSubsidy,
				IsComplete:  false,
			}
			if err := activity.parseCourts(v.CourtsAndTime); err != nil {
				return err
			}
			b.activities = append(b.activities, activity)
		}
	}

	return nil
}

func (b *getActivities) GetSingleParam(attr string) string {
	switch attr {
	default:
		return ""
	}
}

func (b *getActivities) LoadSingleParam(attr, text string) (resultValue interface{}, resultErr error) {
	switch attr {
	default:
	}

	return
}

func (b *getActivities) Do(text string) (resultErr error) {
	size := linebotDomain.MD_FLEX_MESSAGE_SIZE
	valueSize := linebotDomain.SM_FLEX_MESSAGE_SIZE

	carouselContents := []*linebotModel.FlexMessagBubbleComponent{}
	for _, activity := range b.activities {
		contents := []interface{}{}
		activityContents := activity.getLineComponents(domain.NewActivityLineTemplate{})
		contents = append(contents, activityContents...)

		minTime, maxTime := activity.getCourtTimeRange()
		valueText := fmt.Sprintf("%s~%s", minTime.Format(commonLogicDomain.TIME_HOUR_MIN_FORMAT), maxTime.Format(commonLogicDomain.TIME_HOUR_MIN_FORMAT))
		valueSize = linebotDomain.MD_FLEX_MESSAGE_SIZE
		timeContent := linebot.GetKeyValueEditComponent(
			"時間",
			valueText,
			nil, nil, &valueSize,
		)
		contents = util.InsertAtIndex(contents, 1, timeContent)

		estimateTitleComponent := linebot.GetFlexMessageTextComponent(
			0,
			linebot.GetFlexMessageTextComponentSpan(
				"預估費用",
				linebotDomain.XL_FLEX_MESSAGE_SIZE,
				linebotDomain.BOLD_FLEX_MESSAGE_WEIGHT,
			),
		)
		contents = append(contents, estimateTitleComponent)

		size = linebotDomain.MD_FLEX_MESSAGE_SIZE
		valueSize = linebotDomain.SM_FLEX_MESSAGE_SIZE
		courtHours := activity.getCourtHours()
		totalBallConsume := domain.ESTIMATE_BALL_CONSUME_PER_HOUR * courtHours
		ballComponent := linebot.GetKeyValueEditComponent(
			"預估羽球消耗",
			strconv.FormatFloat(totalBallConsume, 'f', -1, 64)+" 顆",
			nil,
			&size,
			&valueSize,
		)
		contents = append(contents, ballComponent)

		estimateBallFee := totalBallConsume * domain.PRICE_PER_BALL
		ballFeeComponent := linebot.GetKeyValueEditComponent(
			"預估羽球費用",
			strconv.FormatFloat(estimateBallFee, 'f', -1, 64),
			nil,
			&size,
			&valueSize,
		)
		contents = append(contents, ballFeeComponent)

		courtFee := activity.getCourtFee()
		courtFeeComponent := linebot.GetKeyValueEditComponent(
			"場地費用",
			strconv.FormatFloat(courtFee, 'f', -1, 64),
			nil,
			&size,
			&valueSize,
		)
		contents = append(contents, courtFeeComponent)

		estimateActivityFee := commonLogic.FloatPlus(estimateBallFee, courtFee)
		activityFeeComponent := linebot.GetKeyValueEditComponent(
			"活動費用",
			strconv.FormatFloat(estimateActivityFee, 'f', -1, 64),
			nil,
			&size,
			&valueSize,
		)
		contents = append(contents, activityFeeComponent)

		if activity.PeopleLimit != nil {
			people := int(*activity.PeopleLimit)
			shareMoney := commonLogic.FloatMinus(estimateActivityFee, float64(activity.ClubSubsidy))

			p := people * domain.MONEY_UNIT
			clubMemberPay := math.Ceil(shareMoney/float64(p)) * domain.MONEY_UNIT
			clubMemberFeeComponent := linebot.GetKeyValueEditComponent(
				"預估社員費用",
				strconv.FormatFloat(clubMemberPay, 'f', -1, 64),
				nil,
				&size,
				&valueSize,
			)
			contents = append(contents, clubMemberFeeComponent)

			guestPay := math.Ceil(estimateActivityFee/float64(p)) * domain.MONEY_UNIT
			guestFeeComponent := linebot.GetKeyValueEditComponent(
				"預估自費費用",
				strconv.FormatFloat(guestPay, 'f', -1, 64),
				nil,
				&size,
				&valueSize,
			)
			contents = append(contents, guestFeeComponent)
		}

		carouselContents = append(
			carouselContents,
			linebot.GetFlexMessageBubbleContent(
				linebot.GetFlexMessageBoxComponent(
					linebotDomain.VERTICAL_MESSAGE_LAYOUT,
					nil,
					linebot.GetFlexMessageBoxComponent(
						linebotDomain.VERTICAL_MESSAGE_LAYOUT,
						nil,
						contents...,
					),
				),
				nil,
			),
		)
	}

	var replyMessge interface{}
	if len(carouselContents) == 0 {
		replyMessge = linebot.GetTextMessage("沒有活動")
	} else {
		replyMessge = linebot.GetFlexMessage(
			"查看活動",
			linebot.GetFlexMessageCarouselContent(carouselContents...),
		)
	}
	replyMessges := []interface{}{
		replyMessge,
	}
	if err := b.context.Reply(replyMessges); err != nil {
		return err
	}

	return nil
}
