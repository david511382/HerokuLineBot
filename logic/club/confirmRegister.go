package club

import (
	"fmt"
	"heroku-line-bot/logic/club/domain"
	commonLogic "heroku-line-bot/logic/common"
	commonLogicDomain "heroku-line-bot/logic/common/domain"
	errLogic "heroku-line-bot/logic/error"
	lineUserLogic "heroku-line-bot/logic/redis/lineuser"
	"heroku-line-bot/service/linebot"
	linebotDomain "heroku-line-bot/service/linebot/domain"
	"heroku-line-bot/storage/database"
	dbReqs "heroku-line-bot/storage/database/domain/model/reqs"
	"heroku-line-bot/storage/redis"
	"time"
)

type confirmRegister struct {
	context  domain.ICmdHandlerContext `json:"-"`
	MemberID int                       `json:"member_id"`
	Date     time.Time                 `json:"date"`
	User     *confirmRegisterUser      `json:"user"`
}

type confirmRegisterUser struct {
	Department Department      `json:"department"`
	Name       string          `json:"name"`
	CompanyID  *string         `json:"company_id"`
	Role       domain.ClubRole `json:"role"`
	LineID     string          `json:"line_id"`
}

func (b *confirmRegister) Init(context domain.ICmdHandlerContext) (resultErrInfo errLogic.IError) {
	*b = confirmRegister{
		context: context,
	}

	return nil
}

func (b *confirmRegister) GetSingleParam(attr string) string {
	switch attr {
	default:
		return ""
	}
}

func (b *confirmRegister) LoadSingleParam(attr, text string) (resultErrInfo errLogic.IError) {
	switch attr {
	case "date":
		t, err := time.Parse(commonLogicDomain.DATE_TIME_RFC3339_FORMAT, text)
		if err != nil {
			resultErrInfo = errLogic.NewError(err)
			return
		}
		b.Date = t
	}

	return nil
}

func (b *confirmRegister) GetInputTemplate(requireRawParamAttr string) interface{} {
	switch requireRawParamAttr {
	default:
		return nil
	}
}

func (b *confirmRegister) LoadUsers(arg dbReqs.Member) (confirmRegisterUsers []*confirmRegisterUser, resultErr error) {
	confirmRegisterUsers = make([]*confirmRegisterUser, 0)
	if dbDatas, err := database.Club.Member.NameRoleDepartmentLineIDCompanyID(arg); err != nil {
		return nil, err
	} else {
		for _, v := range dbDatas {
			confirmRegisterUser := &confirmRegisterUser{
				Department: Department(v.Department),
				Name:       v.Name,
				CompanyID:  v.CompanyID,
				Role:       domain.ClubRole(v.Role.Role),
				LineID:     *v.LineID,
			}

			confirmRegisterUsers = append(confirmRegisterUsers, confirmRegisterUser)
		}
	}
	return confirmRegisterUsers, nil
}

func (b *confirmRegister) ComfirmDb() (resultErrInfo errLogic.IError) {
	isChangeRole := b.isMemberAble() && b.User.Role == domain.GUEST_CLUB_ROLE

	transaction := database.Club.Begin()
	if err := transaction.Error; err != nil {
		resultErrInfo = errLogic.NewError(err)
		return
	}
	defer database.CommitTransaction(transaction, resultErrInfo)

	arg := dbReqs.Member{
		ID: &b.MemberID,
	}
	fields := map[string]interface{}{
		"join_date": b.Date,
	}
	if isChangeRole {
		fields["role"] = int16(domain.MEMBER_CLUB_ROLE)
	}
	if err := database.Club.Member.Update(transaction, arg, fields); err != nil {
		resultErrInfo = errLogic.NewError(err)
		return
	}

	if isChangeRole {
		if _, err := redis.LineUser.Del(b.User.LineID); err != nil {
			resultErrInfo = errLogic.NewError(err)
			return
		}
	}

	return nil
}

func (b *confirmRegister) Do(text string) (resultErrInfo errLogic.IError) {
	lineID := b.context.GetUserID()
	if user, err := lineUserLogic.Get(lineID); err != nil {
		resultErrInfo = errLogic.NewError(err)
		return
	} else if user.Role != domain.ADMIN_CLUB_ROLE {
		resultErrInfo = errLogic.NewError(domain.NO_AUTH_ERROR)
		return
	}

	if b.User == nil {
		arg := dbReqs.Member{
			ID: &b.MemberID,
		}
		if confirmRegisterUsers, err := b.LoadUsers(arg); err != nil {
			resultErrInfo = errLogic.NewError(err)
			return
		} else if len(confirmRegisterUsers) == 0 {
			errInfo := errLogic.New("????????????")
			errInfo = errInfo.Trace()
			resultErrInfo = errInfo
			return
		} else {
			v := confirmRegisterUsers[0]
			b.User = v
		}
	}

	if b.context.IsComfirmed() {
		if errInfo := b.ComfirmDb(); errInfo != nil {
			resultErrInfo = errInfo
			return
		}

		replyMessges := []interface{}{
			linebot.GetTextMessage("??????"),
		}
		if err := b.context.Reply(replyMessges); err != nil {
			resultErrInfo = errLogic.NewError(err)
			return
		}

		if err := b.context.DeleteParam(); err != nil {
			resultErrInfo = errLogic.NewError(err)
			return
		}

		return nil
	}

	if errInfo := b.context.CacheParams(); errInfo != nil {
		resultErrInfo = errInfo
		return
	}

	replyMessges, errInfo := b.getTemplateMessage()
	if errInfo != nil {
		resultErrInfo = errInfo
		return
	}

	if err := b.context.Reply(replyMessges); err != nil {
		resultErrInfo = errLogic.NewError(err)
		return
	}

	return
}

func (b *confirmRegister) isMemberAble() bool {
	return b.User.CompanyID != nil &&
		b.User.Role == domain.GUEST_CLUB_ROLE &&
		b.User.Department.IsClubMember()
}

func (b *confirmRegister) getTemplateMessage() ([]interface{}, errLogic.IError) {
	if b.User == nil {
		return nil, nil
	}

	contents := []interface{}{}
	size := linebotDomain.MD_FLEX_MESSAGE_SIZE
	keyValueEditComponentOption := &domain.KeyValueEditComponentOption{
		SizeP: &size,
	}

	if js, errInfo := b.context.
		GetDateTimeCmdInputMode(domain.DATE_POSTBACK_DATE_TIME_CMD, "date").
		GetSignal(); errInfo != nil {
		return nil, errInfo
	} else {
		action := linebot.GetTimeAction(
			"??????",
			js,
			"",
			"",
			linebotDomain.DATE_TIME_ACTION_MODE,
		)
		contents = append(contents,
			GetKeyValueEditComponent(
				"??????",
				fmt.Sprintf("%s(%s)", b.Date.Format(commonLogicDomain.DATE_FORMAT), commonLogic.WeekDayName(b.Date.Weekday())),
				&domain.KeyValueEditComponentOption{
					Action:     action,
					ValueSizeP: &size,
				},
			),
		)
	}

	???, ???, ??? := b.User.Department.Split()
	if ??? == "" {
		??? = "???"
	}
	contents = append(contents,
		GetKeyValueEditComponent(
			"???",
			string(???),
			keyValueEditComponentOption,
		),
	)
	if ??? == "" {
		??? = "???"
	}
	contents = append(contents,
		GetKeyValueEditComponent(
			"???",
			???,
			keyValueEditComponentOption,
		),
	)
	if ??? == "" {
		??? = "???"
	}
	contents = append(contents,
		GetKeyValueEditComponent(
			"???",
			???,
			keyValueEditComponentOption,
		),
	)

	companyID := "???"
	if b.User.CompanyID != nil {
		companyID = *b.User.CompanyID
	}
	contents = append(contents,
		GetKeyValueEditComponent(
			"????????????",
			companyID,
			keyValueEditComponentOption,
		),
	)

	contents = append(contents,
		GetKeyValueEditComponent(
			"??????",
			b.User.Name,
			keyValueEditComponentOption,
		),
	)

	contents = append(contents,
		GetKeyValueEditComponent(
			"??????",
			b.User.Role.Name(),
			keyValueEditComponentOption,
		),
	)

	if b.User.Role == domain.GUEST_CLUB_ROLE &&
		b.isMemberAble() {
		contents = append(contents,
			linebot.GetTextMessage("???????????????"),
		)
	}

	comfirmSignlJs, errInfo := b.context.
		GetComfirmMode().
		GetSignal()
	if errInfo != nil {
		return nil, errInfo
	}
	contents = append(contents,
		linebot.GetClassButtonComponent(
			linebot.GetPostBackAction(
				"??????",
				comfirmSignlJs,
			),
		),
	)

	return []interface{}{
		linebot.GetFlexMessage(
			"????????????",
			linebot.GetFlexMessageBubbleContent(
				linebot.GetFlexMessageBoxComponent(
					linebotDomain.VERTICAL_MESSAGE_LAYOUT,
					nil,
					contents...,
				),
				nil,
			),
		),
	}, nil
}
