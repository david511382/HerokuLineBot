package club

import (
	"heroku-line-bot/src/logic/club/domain"
	commonLogic "heroku-line-bot/src/logic/common"
	"heroku-line-bot/src/pkg/service/linebot"
	linebotDomain "heroku-line-bot/src/pkg/service/linebot/domain"
	linebotModel "heroku-line-bot/src/pkg/service/linebot/domain/model"
	"heroku-line-bot/src/pkg/util"
	errUtil "heroku-line-bot/src/pkg/util/error"
	"heroku-line-bot/src/repo/database"
	"heroku-line-bot/src/repo/database/database/clubdb/member"
)

type GetConfirmRegisters struct {
	context               domain.ICmdHandlerContext `json:"-"`
	confirmRegistersUsers []*confirmRegistersUser
}

type confirmRegistersUser struct {
	MemberID int    `json:"member_id"`
	Name     string `json:"name"`
}

func (b *GetConfirmRegisters) Init(context domain.ICmdHandlerContext) (resultErrInfo errUtil.IError) {
	*b = GetConfirmRegisters{
		context:               context,
		confirmRegistersUsers: make([]*confirmRegistersUser, 0),
	}

	return nil
}

func (b *GetConfirmRegisters) GetRequireAttr() (requireAttr string, warnMessage interface{}, resultErrInfo errUtil.IError) {
	return
}

func (b *GetConfirmRegisters) GetRequireAttrInfo(rawAttr string) (attrNameText string, valueText string, isNotRequireChecking bool) {
	switch rawAttr {
	default:
		return
	}
}

func (b *GetConfirmRegisters) GetInputTemplate(attr string) (messages interface{}) {
	switch attr {
	}
	return
}

func (b *GetConfirmRegisters) LoadRequireInputTextParam(attr, text string) (resultErrInfo errUtil.IError) {
	switch attr {
	default:
	}

	return nil
}

func (b *GetConfirmRegisters) LoadConfirmRegisterUsers() (resultErrInfo errUtil.IError) {
	arg := member.Reqs{
		CompanyIDIsNull: util.PointerOf(false),
		LineIDIsNull:    util.PointerOf(false),
		JoinDateIsNull:  util.PointerOf(true),
	}
	if dbDatas, err := database.Club().Member.Select(
		arg,
		member.COLUMN_ID,
		member.COLUMN_Name,
	); err != nil {
		resultErrInfo = errUtil.NewError(err)
		return
	} else {
		for _, v := range dbDatas {
			confirmRegisterUser := &confirmRegistersUser{
				Name:     v.Name,
				MemberID: v.ID,
			}

			b.confirmRegistersUsers = append(b.confirmRegistersUsers, confirmRegisterUser)
		}
	}
	return nil
}

func (b *GetConfirmRegisters) Do(text string) (resultErrInfo errUtil.IError) {
	if user, isAutoRegiste, errInfo := autoRegiste(b.context); errInfo != nil {
		resultErrInfo = errUtil.Append(resultErrInfo, errInfo)
		if resultErrInfo.IsError() {
			return
		}
	} else if isAutoRegiste {
		replyMessges := autoRegisteMessage()
		if err := b.context.Reply(replyMessges); err != nil {
			resultErrInfo = errUtil.NewError(err)
			return
		}
	} else if user.Role != domain.ADMIN_CLUB_ROLE {
		resultErrInfo = errUtil.NewError(domain.NO_AUTH_ERROR)
		return
	}

	if err := b.LoadConfirmRegisterUsers(); err != nil {
		resultErrInfo = errUtil.NewError(err)
		return
	}

	replyMessges, err := b.GetConfirmRegisterUsersMessages("待確認入社社員")
	if err != nil {
		resultErrInfo = errUtil.NewError(err)
		return
	}

	if err := b.context.Reply(replyMessges); err != nil {
		resultErrInfo = errUtil.NewError(err)
		return
	}

	return nil
}

func (b *GetConfirmRegisters) GetConfirmRegisterUsersMessages(altText string) (replyMessges []interface{}, resultErr error) {
	replyMessges = make([]interface{}, 0)
	if len(b.confirmRegistersUsers) == 0 {
		replyMessges = append(replyMessges, linebot.GetTextMessage("沒有資料"))
		return
	}

	commonLogic.BatchDo(
		linebotDomain.CAROUSEL_CONTENTS_LIMIT,
		len(b.confirmRegistersUsers),
		func(i, last int) bool {
			carouselContents := []*linebotModel.FlexMessagBubbleComponent{}
			for _, confirmRegistersUser := range b.confirmRegistersUsers[i:last] {
				registerHandler := &registeCompany{
					context:  b.context,
					MemberID: &confirmRegistersUser.MemberID,
				}
				contents, err := registerHandler.GetNotifyRegisterContents(confirmRegistersUser.Name)
				if err != nil {
					resultErr = err
					return false
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

			replyMessges = append(replyMessges, linebot.GetFlexMessage(
				altText,
				linebot.GetFlexMessageCarouselContent(carouselContents...),
			))
			return true
		},
	)
	if resultErr != nil {
		return nil, resultErr
	}

	return
}
