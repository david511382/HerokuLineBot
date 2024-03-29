package club

import (
	"fmt"
	"heroku-line-bot/src/logic/club/domain"
	"heroku-line-bot/src/pkg/service/linebot"
	linebotDomain "heroku-line-bot/src/pkg/service/linebot/domain"
	"heroku-line-bot/src/pkg/util"
	errUtil "heroku-line-bot/src/pkg/util/error"
	"heroku-line-bot/src/repo/database"
	"heroku-line-bot/src/repo/database/database/clubdb/member"
	"heroku-line-bot/src/repo/redis"

	"github.com/rs/zerolog"
)

type confirmRegister struct {
	context  domain.ICmdHandlerContext `json:"-"`
	MemberID uint                      `json:"member_id"`
	domain.TimePostbackParams
	User *confirmRegisterUser `json:"user"`
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

func (b *confirmRegister) GetRequireAttr() (requireAttr string, warnMessage interface{}, resultErrInfo errUtil.IError) {
	return
}

func (b *confirmRegister) GetRequireAttrInfo(rawAttr string) (attrNameText string, valueText string, isNotRequireChecking bool) {
	switch rawAttr {
	default:
		return
	}
}

func (b *confirmRegister) GetInputTemplate(attr string) (messages interface{}) {
	switch attr {
	}
	return
}

func (b *confirmRegister) LoadRequireInputTextParam(attr, text string) (resultErrInfo errUtil.IError) {
	switch attr {
	default:
	}

	return nil
}

func (b *confirmRegister) LoadUsers(arg member.Reqs) (confirmRegisterUsers []*confirmRegisterUser, resultErr error) {
	confirmRegisterUsers = make([]*confirmRegisterUser, 0)
	if dbDatas, err := database.Club().Member.Select(
		arg,
		member.COLUMN_Name,
		member.COLUMN_Role,
		member.COLUMN_Department,
		member.COLUMN_LineID,
		member.COLUMN_CompanyID,
	); err != nil {
		return nil, err
	} else {
		for _, v := range dbDatas {
			confirmRegisterUser := &confirmRegisterUser{
				Department: Department(v.Department),
				Name:       v.Name,
				CompanyID:  v.CompanyID,
				Role:       domain.ClubRole(v.Role),
				LineID:     *v.LineID,
			}

			confirmRegisterUsers = append(confirmRegisterUsers, confirmRegisterUser)
		}
	}
	return confirmRegisterUsers, nil
}

func (b *confirmRegister) ConfirmDb() (resultErrInfo errUtil.IError) {
	isChangeRole := b.isMemberAble() && b.User.Role == domain.GUEST_CLUB_ROLE

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

	dateP := b.Date.TimeP()
	arg := member.UpdateReqs{
		Reqs: member.Reqs{
			ID: &b.MemberID,
		},
		JoinDate: &dateP,
	}
	if isChangeRole {
		arg.Role = util.PointerOf(int16(domain.MEMBER_CLUB_ROLE))
	}
	if err := db.Member.Update(arg); err != nil {
		resultErrInfo = errUtil.NewError(err)
		return
	}

	if isChangeRole {
		if _, errInfo := redis.Badminton().LineUser.HDel(b.User.LineID); errInfo != nil {
			errInfo.SetLevel(zerolog.WarnLevel)
			resultErrInfo = errUtil.Append(resultErrInfo, errInfo)
		}
	}

	return
}

func (b *confirmRegister) Do(text string) (resultErrInfo errUtil.IError) {
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

	if b.User == nil {
		arg := member.Reqs{
			ID: &b.MemberID,
		}
		if confirmRegisterUsers, err := b.LoadUsers(arg); err != nil {
			resultErrInfo = errUtil.NewError(err)
			return
		} else if len(confirmRegisterUsers) == 0 {
			errInfo := errUtil.New("查無用戶")
			resultErrInfo = errUtil.Append(resultErrInfo, errInfo)
			return
		} else {
			v := confirmRegisterUsers[0]
			b.User = v
		}
	}

	if b.context.IsConfirmed() {
		if errInfo := b.ConfirmDb(); errInfo != nil {
			resultErrInfo = errUtil.Append(resultErrInfo, errInfo)
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
		resultErrInfo = errUtil.Append(resultErrInfo, errInfo)
		return
	}

	replyMessges, errInfo := b.getTemplateMessage()
	if errInfo != nil {
		resultErrInfo = errUtil.Append(resultErrInfo, errInfo)
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

	if js, errInfo := NewSignal().
		GetBasePath("ICmdLogic").
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
				fmt.Sprintf("%s(%s)", b.TimePostbackParams.Date.Time().Format(util.DATE_FORMAT), util.GetWeekDayName(b.Date.Time().Weekday())),
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

	comfirmSignlJs, errInfo := NewSignal().
		GetConfirmMode().
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
