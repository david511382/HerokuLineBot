package club

import (
	accountLineuserLogic "heroku-line-bot/src/logic/account/lineuser"
	"heroku-line-bot/src/logic/club/domain"
	commonLogic "heroku-line-bot/src/logic/common"
	dbModel "heroku-line-bot/src/model/database"
	"heroku-line-bot/src/pkg/service/linebot"
	linebotDomain "heroku-line-bot/src/pkg/service/linebot/domain"
	linebotModel "heroku-line-bot/src/pkg/service/linebot/domain/model"
	"heroku-line-bot/src/pkg/util"
	errUtil "heroku-line-bot/src/pkg/util/error"
	"heroku-line-bot/src/repo/database"
	"heroku-line-bot/src/repo/database/database/clubdb/member"
)

type GetComfirmRegisters struct {
	context               domain.ICmdHandlerContext `json:"-"`
	confirmRegistersUsers []*confirmRegistersUser
}

type confirmRegistersUser struct {
	MemberID int    `json:"member_id"`
	Name     string `json:"name"`
}

func (b *GetComfirmRegisters) Init(context domain.ICmdHandlerContext) (resultErrInfo errUtil.IError) {
	*b = GetComfirmRegisters{
		context:               context,
		confirmRegistersUsers: make([]*confirmRegistersUser, 0),
	}

	return nil
}

func (b *GetComfirmRegisters) GetSingleParam(attr string) string {
	switch attr {
	default:
		return ""
	}
}

func (b *GetComfirmRegisters) LoadSingleParam(attr, text string) (resultErrInfo errUtil.IError) {
	switch attr {
	default:
	}

	return nil
}

func (b *GetComfirmRegisters) GetInputTemplate(requireRawParamAttr string) interface{} {
	return nil
}

func (b *GetComfirmRegisters) LoadComfirmRegisterUsers() (resultErrInfo errUtil.IError) {
	arg := dbModel.ReqsClubMember{
		CompanyIDIsNull: util.GetBoolP(false),
		LineIDIsNull:    util.GetBoolP(false),
		JoinDateIsNull:  util.GetBoolP(true),
	}
	if dbDatas, err := database.Club.Member.Select(
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

func (b *GetComfirmRegisters) Do(text string) (resultErrInfo errUtil.IError) {
	lineID := b.context.GetUserID()
	if user, err := accountLineuserLogic.Get(lineID); err != nil {
		resultErrInfo = errUtil.NewError(err)
		return
	} else if user.Role != domain.ADMIN_CLUB_ROLE {
		resultErrInfo = errUtil.NewError(domain.NO_AUTH_ERROR)
		return
	}

	if err := b.LoadComfirmRegisterUsers(); err != nil {
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

func (b *GetComfirmRegisters) GetConfirmRegisterUsersMessages(altText string) (replyMessges []interface{}, resultErr error) {
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
				registerHandler := &register{
					context:  b.context,
					Name:     confirmRegistersUser.Name,
					MemberID: confirmRegistersUser.MemberID,
				}
				contents, err := registerHandler.GetNotifyRegisterContents()
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
