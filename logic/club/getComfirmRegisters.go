package club

import (
	"heroku-line-bot/logic/club/domain"
	lineUserLogic "heroku-line-bot/logic/redis/lineuser"
	"heroku-line-bot/service/linebot"
	linebotDomain "heroku-line-bot/service/linebot/domain"
	linebotModel "heroku-line-bot/service/linebot/domain/model"
	"heroku-line-bot/storage/database"
	dbReqs "heroku-line-bot/storage/database/domain/model/reqs"
	"heroku-line-bot/util"
)

type GetComfirmRegisters struct {
	context               domain.ICmdHandlerContext `json:"-"`
	confirmRegistersUsers []*confirmRegistersUser
}

type confirmRegistersUser struct {
	MemberID int    `json:"member_id"`
	Name     string `json:"name"`
}

func (b *GetComfirmRegisters) Init(context domain.ICmdHandlerContext) error {
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

func (b *GetComfirmRegisters) LoadSingleParam(attr, text string) error {
	switch attr {
	default:
	}

	return nil
}

func (b *GetComfirmRegisters) GetInputTemplate(requireRawParamAttr string) interface{} {
	return nil
}

func (b *GetComfirmRegisters) LoadComfirmRegisterUsers() error {
	arg := dbReqs.Member{
		CompanyIDIsNull: util.GetBoolP(false),
		LineIDIsNull:    util.GetBoolP(false),
		JoinDateIsNull:  util.GetBoolP(true),
	}
	if dbDatas, err := database.Club.Member.IDName(arg); err != nil {
		return err
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

func (b *GetComfirmRegisters) Do(text string) (resultErr error) {
	lineID := b.context.GetUserID()
	if user, err := lineUserLogic.Get(lineID); err != nil {
		return err
	} else if user.Role != domain.ADMIN_CLUB_ROLE {
		return domain.NO_AUTH_ERROR
	}

	if err := b.LoadComfirmRegisterUsers(); err != nil {
		return err
	}

	replyMessge := b.GetConfirmRegisterUsersMessage("待確認入社社員")
	replyMessges := []interface{}{
		replyMessge,
	}
	if err := b.context.Reply(replyMessges); err != nil {
		return err
	}

	return nil
}

func (b *GetComfirmRegisters) GetConfirmRegisterUsersMessage(altText string) (replyMessge interface{}) {
	carouselContents := []*linebotModel.FlexMessagBubbleComponent{}
	for _, confirmRegistersUser := range b.confirmRegistersUsers {
		registerHandler := &register{
			context:  b.context,
			Name:     confirmRegistersUser.Name,
			MemberID: confirmRegistersUser.MemberID,
		}
		contents, err := registerHandler.GetNotifyRegisterContents()
		if err != nil {
			return err
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

	if len(carouselContents) == 0 {
		replyMessge = linebot.GetTextMessage("沒有資料")
	} else {
		replyMessge = linebot.GetFlexMessage(
			altText,
			linebot.GetFlexMessageCarouselContent(carouselContents...),
		)
	}

	return
}
