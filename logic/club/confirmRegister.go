package club

import (
	"fmt"
	"heroku-line-bot/logic/club/domain"
	clubLineuserLogic "heroku-line-bot/logic/club/lineuser"
	"heroku-line-bot/service/linebot"
	linebotDomain "heroku-line-bot/service/linebot/domain"
	"heroku-line-bot/storage/database"
	dbReqs "heroku-line-bot/storage/database/domain/model/reqs"
	"heroku-line-bot/storage/redis"
	"heroku-line-bot/util"
	errUtil "heroku-line-bot/util/error"
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

func (b *confirmRegister) Init(context domain.ICmdHandlerContext) (resultErrInfo errUtil.IError) {
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

func (b *confirmRegister) LoadSingleParam(attr, text string) (resultErrInfo errUtil.IError) {
	switch attr {
	case "date":
		t, err := time.Parse(util.DATE_TIME_RFC3339_FORMAT, text)
		if err != nil {
			resultErrInfo = errUtil.NewError(err)
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

func (b *confirmRegister) ComfirmDb() (resultErrInfo errUtil.IError) {
	isChangeRole := b.isMemberAble() && b.User.Role == domain.GUEST_CLUB_ROLE

	transaction := database.Club.Begin()
	if err := transaction.Error; err != nil {
		resultErrInfo = errUtil.NewError(err)
		return
	}
	defer func() {
		if errInfo := database.CommitTransaction(transaction, resultErrInfo); errInfo != nil {
			resultErrInfo = errUtil.Append(resultErrInfo, errInfo)
		}
	}()

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
		resultErrInfo = errUtil.NewError(err)
		return
	}

	if isChangeRole {
		if errInfo := redis.LineUser.Del(b.User.LineID); errInfo != nil {
			errInfo.SetLevel(errUtil.WARN)
			resultErrInfo = errUtil.Append(resultErrInfo, errInfo)
		}
	}

	return
}

func (b *confirmRegister) Do(text string) (resultErrInfo errUtil.IError) {
	lineID := b.context.GetUserID()
	if user, err := clubLineuserLogic.Get(lineID); err != nil {
		resultErrInfo = errUtil.NewError(err)
		return
	} else if user.Role != domain.ADMIN_CLUB_ROLE {
		resultErrInfo = errUtil.NewError(domain.NO_AUTH_ERROR)
		return
	}

	if b.User == nil {
		arg := dbReqs.Member{
			ID: &b.MemberID,
		}
		if confirmRegisterUsers, err := b.LoadUsers(arg); err != nil {
			resultErrInfo = errUtil.NewError(err)
			return
		} else if len(confirmRegisterUsers) == 0 {
			errInfo := errUtil.New("查無用戶")
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
		resultErrInfo = errUtil.NewError(err)
		return
	}

	return
}

func (b *confirmRegister) isMemberAble() bool {
	return b.User.CompanyID != nil &&
		b.User.Role == domain.GUEST_CLUB_ROLE &&
		b.User.Department.IsClubMember()
}

func (b *confirmRegister) getTemplateMessage() ([]interface{}, errUtil.IError) {
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
			"修改",
			js,
			"",
			"",
			linebotDomain.DATE_TIME_ACTION_MODE,
		)
		contents = append(contents,
			GetKeyValueEditComponent(
				"日期",
				fmt.Sprintf("%s(%s)", b.Date.Format(util.DATE_FORMAT), util.GetWeekDayName(b.Date.Weekday())),
				&domain.KeyValueEditComponentOption{
					Action:     action,
					ValueSizeP: &size,
				},
			),
		)
	}

	處, 部, 組 := b.User.Department.Split()
	if 處 == "" {
		處 = "無"
	}
	contents = append(contents,
		GetKeyValueEditComponent(
			"處",
			string(處),
			keyValueEditComponentOption,
		),
	)
	if 部 == "" {
		部 = "無"
	}
	contents = append(contents,
		GetKeyValueEditComponent(
			"部",
			部,
			keyValueEditComponentOption,
		),
	)
	if 組 == "" {
		組 = "無"
	}
	contents = append(contents,
		GetKeyValueEditComponent(
			"組",
			組,
			keyValueEditComponentOption,
		),
	)

	companyID := "無"
	if b.User.CompanyID != nil {
		companyID = *b.User.CompanyID
	}
	contents = append(contents,
		GetKeyValueEditComponent(
			"員工編號",
			companyID,
			keyValueEditComponentOption,
		),
	)

	contents = append(contents,
		GetKeyValueEditComponent(
			"暱稱",
			b.User.Name,
			keyValueEditComponentOption,
		),
	)

	contents = append(contents,
		GetKeyValueEditComponent(
			"角色",
			b.User.Role.Name(),
			keyValueEditComponentOption,
		),
	)

	if b.User.Role == domain.GUEST_CLUB_ROLE &&
		b.isMemberAble() {
		contents = append(contents,
			linebot.GetTextMessage("具會員資格"),
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
				"確認",
				comfirmSignlJs,
			),
		),
	)

	return []interface{}{
		linebot.GetFlexMessage(
			"新人註冊",
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
