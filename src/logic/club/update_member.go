package club

import (
	"heroku-line-bot/src/logic/club/domain"
	"heroku-line-bot/src/pkg/service/linebot"
	linebotDomain "heroku-line-bot/src/pkg/service/linebot/domain"
	"heroku-line-bot/src/pkg/util"
	errUtil "heroku-line-bot/src/pkg/util/error"
	"heroku-line-bot/src/repo/database"
	"heroku-line-bot/src/repo/database/database/clubdb/member"
)

type UpdateMember struct {
	context domain.ICmdHandlerContext `json:"-"`
	Name    *string                   `json:"name"`
}

func NewUpdateMember() *UpdateMember {
	return &UpdateMember{}
}

func (b *UpdateMember) Init(context domain.ICmdHandlerContext) (resultErrInfo errUtil.IError) {
	b.context = context

	if errInfo := b.LoadName(); errInfo != nil {
		resultErrInfo = errUtil.Append(resultErrInfo, errInfo)
		if errInfo.IsError() {
			return
		}
	}

	return
}

func (b *UpdateMember) GetRequireAttr() (requireAttr string, warnMessage interface{}, resultErrInfo errUtil.IError) {
	return
}

func (b *UpdateMember) GetRequireAttrInfo(rawAttr string) (attrNameText string, valueText string, isNotRequireChecking bool) {
	switch rawAttr {
	case "name":
		attrNameText = "暱稱"
		if b.Name == nil {
			if errInfo := b.LoadName(); errInfo != nil {
				if errInfo.IsError() {
					valueText = ""
					return
				}
			}
		}
		valueText = *b.Name
	}
	return
}

func (b *UpdateMember) GetInputTemplate(attr string) (messages interface{}) {
	switch attr {
	}
	return
}

func (b *UpdateMember) LoadRequireInputTextParam(attr, text string) (resultErrInfo errUtil.IError) {
	switch attr {
	case "name":
		b.Name = util.PointerOf(text)
	default:
	}
	return nil
}

func (b *UpdateMember) Do(text string) (resultErrInfo errUtil.IError) {
	if b.Name == nil {
		if errInfo := b.LoadName(); errInfo != nil {
			resultErrInfo = errUtil.Append(resultErrInfo, errInfo)
			if errInfo.IsError() {
				return
			}
		}
	}

	if b.Name == nil {
		replyMessges := []interface{}{
			linebot.GetTextMessage("您尚未註冊，請封鎖本帳號後再次加入"),
		}
		if err := b.context.Reply(replyMessges); err != nil {
			resultErrInfo = errUtil.NewError(err)
			return
		}

		if err := b.context.DeleteParam(); err != nil {
			resultErrInfo = errUtil.NewError(err)
			return
		}
		return
	}

	name := *b.Name
	if b.context.IsConfirmed() {
		db, transaction, err := database.Club().Begin()
		if err != nil {
			errInfo := errUtil.NewError(err)
			resultErrInfo = errUtil.Append(resultErrInfo, errInfo)
			return
		}
		defer func() {
			if errInfo := database.CommitTransaction(transaction, resultErrInfo); errInfo != nil {
				resultErrInfo = errUtil.Append(resultErrInfo, errInfo)
			}
		}()

		lineID := b.context.GetUserID()
		if err := db.Member.Update(member.UpdateReqs{
			Reqs: member.Reqs{
				LineID: &lineID,
			},
			Name: b.Name,
		}); err != nil {
			resultErrInfo = errUtil.NewError(err)
			return
		}

		replyMessges := []interface{}{
			linebot.GetTextMessage("完成"),
		}
		if err := b.context.Reply(replyMessges); err != nil {
			resultErrInfo = errUtil.NewError(err)
			return
		}

		if err := b.context.DeleteParam(); err != nil {
			resultErrInfo = errUtil.NewError(err)
			return
		}

		return
	}

	contents := make([]interface{}, 0)

	{
		size := linebotDomain.MD_FLEX_MESSAGE_SIZE
		js, errInfo := NewSignal().
			GetRequireInputMode("name").
			GetSignal()
		if errInfo != nil {
			resultErrInfo = errUtil.Append(resultErrInfo, errInfo)
			if resultErrInfo.IsError() {
				return
			}
		}

		action := linebot.GetPostBackAction(
			"修改",
			js,
		)
		contents = append(contents,
			GetKeyValueEditComponent(
				"暱稱",
				name,
				&domain.KeyValueEditComponentOption{
					Action: action,
					SizeP:  &size,
				},
			),
		)
	}

	{
		cancelSignlJs, errInfo := NewSignal().
			GetCancelMode().
			GetSignal()
		if errInfo != nil {
			resultErrInfo = errUtil.Append(resultErrInfo, errInfo)
			return
		}
		comfirmSignlJs, errInfo := NewSignal().
			GetConfirmMode().
			GetSignal()
		if errInfo != nil {
			resultErrInfo = errUtil.Append(resultErrInfo, errInfo)
			return
		}
		contents = append(contents,
			GetConfirmComponent(
				linebot.GetPostBackAction(
					"取消",
					cancelSignlJs,
				),
				linebot.GetPostBackAction(
					"確認",
					comfirmSignlJs,
				),
			),
		)
	}

	replyMessges := []interface{}{
		linebot.GetFlexMessage(
			"修改基本資訊",
			linebot.GetFlexMessageBubbleContent(
				linebot.GetFlexMessageBoxComponent(
					linebotDomain.VERTICAL_MESSAGE_LAYOUT,
					nil,
					contents...,
				),
				nil,
			),
		),
	}
	if err := b.context.Reply(replyMessges); err != nil {
		resultErrInfo = errUtil.NewError(err)
		return
	}

	return
}

func (b *UpdateMember) LoadName() (resultErrInfo errUtil.IError) {
	lineID := b.context.GetUserID()
	dbDatas, err := database.Club().Member.Select(
		member.Reqs{
			LineID: &lineID,
		},
		member.COLUMN_Name,
	)
	if err != nil {
		errInfo := errUtil.NewError(err)
		resultErrInfo = errUtil.Append(resultErrInfo, errInfo)
		return
	} else if len(dbDatas) == 0 {
		return
	}

	dbData := dbDatas[0]
	b.Name = &dbData.Name

	return
}
