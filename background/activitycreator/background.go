package activitycreator

import (
	"heroku-line-bot/bootstrap"
	clubLogic "heroku-line-bot/logic/club"
	clubCourtLogic "heroku-line-bot/logic/club/court"
	clubCourtLogicDomain "heroku-line-bot/logic/club/court/domain"
	"heroku-line-bot/logic/club/domain"
	clubLineBotLogic "heroku-line-bot/logic/clublinebot"
	commonLogic "heroku-line-bot/logic/common"
	commonLogicDomain "heroku-line-bot/logic/common/domain"
	errLogic "heroku-line-bot/logic/error"
	rdsBadmintonSettingLogic "heroku-line-bot/logic/redis/badmintonsetting"
	"heroku-line-bot/models/storage"
	"heroku-line-bot/service/linebot"
	"heroku-line-bot/storage/database"
	"heroku-line-bot/util"
	"sort"
	"strconv"
	"time"
)

type BackGround struct{}

func (b *BackGround) Init(cfg bootstrap.Backgrounds) (name string, backgroundCfg bootstrap.Background, resultErrInfo errLogic.IError) {
	return "ActivityCreator", cfg.ActivityCreator, nil
}

func (b *BackGround) Run(runTime time.Time) (resultErrInfo errLogic.IError) {
	defer func() {
		if resultErrInfo != nil {
			resultErrInfo = resultErrInfo.NewParent(runTime.String())
		}
	}()

	rdsSetting, errInfo := rdsBadmintonSettingLogic.Get()
	if errInfo != nil {
		resultErrInfo = errInfo
		if resultErrInfo.IsError() {
			return
		}
	}
	if rdsSetting == nil {
		resultErrInfo = errLogic.New("no redis setting", errLogic.WARN)
		rdsSetting = &storage.BadmintonActivity{
			Description: "7人出團",
			ClubSubsidy: 0,
		}
	}

	newActivityHandlers := make([]*clubLogic.NewActivity, 0)
	createActivityDate := commonLogicDomain.WEEK_TIME_TYPE.Next(
		commonLogicDomain.DATE_TIME_TYPE.Of(runTime),
		1,
	)
	weekday := int16(createActivityDate.Weekday())
	if courtDatas, errInfo := clubCourtLogic.GetRentalCourts(
		createActivityDate,
		createActivityDate,
		nil,
		&weekday,
	); errInfo != nil {
		resultErrInfo = errInfo
		return
	} else if len(courtDatas) == 0 {
		return
	} else {
		for place, dateMap := range courtDatas {
			for dateInt, courtData := range dateMap {
				date := commonLogic.IntTime(dateInt, commonLogicDomain.DATE_TIME_TYPE)
				courts := b.combineCourts(courtData.Courts)
				newActivityHandler := &clubLogic.NewActivity{
					Date:        date,
					Place:       place,
					Description: rdsSetting.Description,
					ClubSubsidy: rdsSetting.ClubSubsidy,
					IsComplete:  false,
					Courts:      courts,
				}

				totalCourtCount := 0
				for _, court := range newActivityHandler.Courts {
					totalCourtCount += int(court.Count)
				}
				peopleLimit := int16(totalCourtCount * domain.PEOPLE_PER_HOUR * 2)

				newActivityHandler.PeopleLimit = util.GetInt16P(peopleLimit)

				newActivityHandlers = append(newActivityHandlers, newActivityHandler)
			}
		}
	}

	if transaction := database.Club.Begin(); transaction.Error != nil {
		resultErrInfo = errLogic.NewError(transaction.Error)
		return
	} else {
		for _, newActivityHandler := range newActivityHandlers {
			if resultErrInfo = newActivityHandler.InsertActivity(transaction); resultErrInfo != nil {
				return
			}
		}

		database.CommitTransaction(transaction, resultErrInfo)
	}

	getActivityHandler := &clubLogic.GetActivities{}
	pushMessage, err := getActivityHandler.GetActivitiesMessage("開放活動報名", false, false)
	if err != nil {
		resultErrInfo = errLogic.NewError(err)
		return
	}
	pushMessages := []interface{}{
		linebot.GetTextMessage("活動開放報名，請私下與我報名"),
		pushMessage,
	}
	linebotContext := clubLineBotLogic.NewContext("", "", &clubLineBotLogic.Bot)
	if err := linebotContext.PushRoom(pushMessages); err != nil {
		resultErrInfo = errLogic.NewError(err)
		return
	}

	return
}

type node struct {
	fromSideMap map[int]*node
	toSideMap   map[int]*node
	fromSides   []int
	toSides     []int
	target      *clubCourtLogicDomain.ActivityCourt
}

func (m *node) isFromSorted() bool {
	return len(m.fromSideMap) == len(m.fromSides)
}

func (m *node) isToSorted() bool {
	return len(m.toSideMap) == len(m.toSides)
}

func (m *node) fromLen() int {
	return len(m.fromSideMap)
}

func (m *node) toLen() int {
	return len(m.toSideMap)
}

func (m *node) value() util.Float {
	return m.target.TotalHours()
}

func (m *node) totalValue() util.Float {
	return m.fromSideMax().Plus(m.value(), m.toSideMax())
}

func (m *node) takeFromSideMax() []*clubCourtLogicDomain.ActivityCourt {
	result := make([]*clubCourtLogicDomain.ActivityCourt, 0)

	for _, index := range m.fromSides {
		n := m.fromSideMap[index]
		if n == nil || n.target == nil {
			continue
		}

		target := *n.target
		result = append(result, &target)
		result = append(result, n.takeFromSideMax()...)
		n.target = nil
		break
	}

	return result
}

func (m *node) takeToSideMax() []*clubCourtLogicDomain.ActivityCourt {
	result := make([]*clubCourtLogicDomain.ActivityCourt, 0)
	for _, index := range m.toSides {
		n := m.toSideMap[index]
		if n == nil || n.target == nil {
			continue
		}

		target := *n.target
		result = append(result, &target)
		result = append(result, n.takeToSideMax()...)
		n.target = nil
		break
	}

	return result
}

func (m *node) takeMax() []*clubCourtLogicDomain.ActivityCourt {
	result := make([]*clubCourtLogicDomain.ActivityCourt, 0)
	if m.target == nil {
		return result
	}

	fromSideNodes := m.takeFromSideMax()
	for i := len(fromSideNodes) - 1; i >= 0; i-- {
		result = append(result, fromSideNodes[i])
	}
	target := *m.target
	result = append(result, &target)
	result = append(result, m.takeToSideMax()...)

	m.target = nil
	return result
}

func (m *node) fromSideMax() util.Float {
	if m.fromLen() == 0 {
		return util.ToFloat(0)
	}

	index := m.fromSides[0]
	next := m.fromSideMap[index]
	return next.fromSideMax().Plus(next.value())
}

func (m *node) toSideMax() util.Float {
	if m.toLen() == 0 {
		return util.ToFloat(0)
	}

	index := m.toSides[0]
	next := m.toSideMap[index]
	return next.toSideMax().Plus(next.value())
}

func (m *node) sort() {
	m.sortFromSide()
	m.sortToSide()
}

func (m *node) sortFromSide() {
	fromSideValueMap := make(map[int]float64)
	m.fromSides = make([]int, 0)
	for i, n := range m.fromSideMap {
		if !n.isFromSorted() {
			n.sortFromSide()
		}

		fromSideValueMap[i] = n.fromSideMax().Plus(n.value()).Value()

		m.fromSides = append(m.fromSides, i)
	}
	sort.SliceStable(m.fromSides, func(i, j int) bool {
		iIndex := m.fromSides[i]
		jIndex := m.fromSides[j]
		return iIndex < jIndex
	})
	sort.SliceStable(m.fromSides, func(i, j int) bool {
		iIndex := m.fromSides[i]
		jIndex := m.fromSides[j]
		return fromSideValueMap[iIndex] > fromSideValueMap[jIndex]
	})
}

func (m *node) sortToSide() {
	toSideValueMap := make(map[int]float64)
	m.toSides = make([]int, 0)
	for i, n := range m.toSideMap {
		if !n.isToSorted() {
			n.sortToSide()
		}

		toSideValueMap[i] = n.toSideMax().Plus(n.value()).Value()

		m.toSides = append(m.toSides, i)
	}
	sort.SliceStable(m.toSides, func(i, j int) bool {
		iIndex := m.toSides[i]
		jIndex := m.toSides[j]
		return iIndex < jIndex
	})
	sort.SliceStable(m.toSides, func(i, j int) bool {
		iIndex := m.toSides[i]
		jIndex := m.toSides[j]
		return toSideValueMap[iIndex] > toSideValueMap[jIndex]
	})
}

func (b *BackGround) combineCourts(courts []*clubCourtLogicDomain.ActivityCourt) []*clubCourtLogicDomain.ActivityCourt {
	newCourts := make([]*clubCourtLogicDomain.ActivityCourt, 0)

	priceCourtsMap := make(map[float64][]*clubCourtLogicDomain.ActivityCourt)
	for _, court := range courts {
		price := court.PricePerHour
		if priceCourtsMap[price] == nil {
			priceCourtsMap[price] = make([]*clubCourtLogicDomain.ActivityCourt, 0)
		}
		priceCourtsMap[price] = append(priceCourtsMap[price], court)
	}
	for _, courts := range priceCourtsMap {
		newCourts = append(newCourts, b.combineSamePriceCourts(courts)...)
	}

	return newCourts
}

func (b *BackGround) combineSamePriceCourts(courts []*clubCourtLogicDomain.ActivityCourt) []*clubCourtLogicDomain.ActivityCourt {
	newCourts := make([]*clubCourtLogicDomain.ActivityCourt, 0)

	for _, court := range courts {
		for c := 1; c < int(court.Count); c++ {
			copyCourt := *court
			courts = append(courts, &copyCourt)
		}
	}

	fromIntMap := make(map[int][]int)
	toIntMap := make(map[int][]int)
	idNodeMap := make(map[int]*node)
	for i, court := range courts {
		idNodeMap[i] = &node{
			fromSideMap: make(map[int]*node),
			toSideMap:   make(map[int]*node),
			target:      court,
		}

		fromClockInt := commonLogic.ClockInt(court.FromTime, commonLogicDomain.MINUTE_TIME_TYPE)
		if len(toIntMap[fromClockInt]) > 0 {
			for _, index := range toIntMap[fromClockInt] {
				idNodeMap[index].toSideMap[i] = idNodeMap[i]
				idNodeMap[i].fromSideMap[index] = idNodeMap[index]
			}
		}
		toClockInt := commonLogic.ClockInt(court.ToTime, commonLogicDomain.MINUTE_TIME_TYPE)
		if len(fromIntMap[toClockInt]) > 0 {
			for _, index := range fromIntMap[toClockInt] {
				idNodeMap[index].fromSideMap[i] = idNodeMap[i]
				idNodeMap[i].toSideMap[index] = idNodeMap[index]
			}
		}

		if toIntMap[toClockInt] == nil {
			toIntMap[toClockInt] = make([]int, 0)
		}
		toIntMap[toClockInt] = append(toIntMap[toClockInt], i)
		if fromIntMap[fromClockInt] == nil {
			fromIntMap[fromClockInt] = make([]int, 0)
		}
		fromIntMap[fromClockInt] = append(fromIntMap[fromClockInt], i)
	}

	nodes := make([]*node, 0)
	for _, m := range idNodeMap {
		m.sort()
		nodes = append(nodes, m)
	}

	sort.Slice(nodes, func(i, j int) bool {
		im := nodes[i]
		iValue := im.totalValue().Value()
		jm := nodes[j]
		jValue := jm.totalValue().Value()
		return iValue > jValue
	})

	timeCountMap := make(map[string]int)
	timeCourtMap := make(map[string]*clubCourtLogicDomain.ActivityCourt)
	for _, m := range nodes {
		maxNodes := m.takeMax()
		var newCourt *clubCourtLogicDomain.ActivityCourt
		if len(maxNodes) == 1 {
			newCourt = maxNodes[0]
		} else if len(maxNodes) == 0 {
			continue
		} else {
			newCourt = &clubCourtLogicDomain.ActivityCourt{
				FromTime:     maxNodes[0].FromTime,
				ToTime:       maxNodes[len(maxNodes)-1].ToTime,
				Count:        1,
				PricePerHour: maxNodes[0].PricePerHour,
			}
		}

		fromClockInt := commonLogic.ClockInt(newCourt.FromTime, commonLogicDomain.MINUTE_TIME_TYPE)
		toClockInt := commonLogic.ClockInt(newCourt.ToTime, commonLogicDomain.MINUTE_TIME_TYPE)
		timeKey := strconv.Itoa(fromClockInt)
		timeKey += strconv.Itoa(toClockInt)

		if timeCourtMap[timeKey] != nil {
			timeCountMap[timeKey]++
		} else {
			timeCourtMap[timeKey] = newCourt
			timeCountMap[timeKey] = 1
		}

	}

	for k, court := range timeCourtMap {
		court.Count = int16(timeCountMap[k])
		newCourts = append(newCourts, court)
	}

	return newCourts
}
