package club

import (
	"fmt"
	"heroku-line-bot/src/logic/account"
	"heroku-line-bot/src/logic/club/domain"
	clublinebotDomain "heroku-line-bot/src/logic/clublinebot/domain"
	"heroku-line-bot/src/pkg/service/linebot"
	"heroku-line-bot/src/pkg/util"
	errUtil "heroku-line-bot/src/pkg/util/error"
	"heroku-line-bot/src/repo/database"
	"heroku-line-bot/src/repo/database/database/clubdb"
	"heroku-line-bot/src/repo/database/database/clubdb/member"
)

type RegisteMember struct {
	context  clublinebotDomain.ILineBotContext
	name     string
	lineID   *string
	memberID *int
}

func NewRegisterMember(name string, lineID *string) *RegisteMember {
	return &RegisteMember{
		name:   name,
		lineID: lineID,
	}
}

func (b *RegisteMember) Init(context clublinebotDomain.ILineBotContext) {
	b.context = context
}

func (b *RegisteMember) LoadMemberID() (
	memberID *int,
	memberName *string,
	resultErrInfo errUtil.IError,
) {
	if b.lineID == nil {
		return
	}

	if dbDatas, err := database.Club().Member.Select(
		member.Reqs{
			LineID: b.lineID,
		},
		member.COLUMN_ID,
		member.COLUMN_Name,
	); err != nil {
		errInfo := errUtil.NewError(err)
		resultErrInfo = errUtil.Append(resultErrInfo, errInfo)
		return
	} else if len(dbDatas) > 0 {
		dbData := dbDatas[0]
		memberID = util.PointerOf(dbData.ID)
		memberName = util.PointerOf(dbData.Name)

		b.memberID = memberID
	}
	return
}

func (b *RegisteMember) NotifyAdmin() (resultErrInfo errUtil.IError) {
	if b.context == nil {
		return
	}

	adminReplyMessges := []interface{}{
		linebot.GetTextMessage(fmt.Sprintf("%s 註冊", b.name)),
	}
	if err := b.context.PushAdmin(adminReplyMessges); err != nil {
		errInfo := errUtil.NewError(err)
		resultErrInfo = errUtil.Append(resultErrInfo, errInfo)
		return
	}
	return
}

func (b *RegisteMember) Registe(db *clubdb.Database) (resultErrInfo errUtil.IError) {
	if b.memberID != nil {
		// 存在的用戶
		return
	}

	data := &member.Model{
		Department: string(NewEmptyDepartment()),
		Name:       b.name,
		Role:       int16(domain.GUEST_CLUB_ROLE),
	}
	if b.lineID != nil {
		data.LineID = b.lineID
	}
	if errInfo := account.Registe(db, data); errInfo != nil {
		resultErrInfo = errUtil.Append(resultErrInfo, errInfo)
		return
	}

	return
}
