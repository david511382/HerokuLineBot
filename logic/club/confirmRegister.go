package club

import (
	"fmt"
	"heroku-line-bot/logic/club/domain"
	commonLogic "heroku-line-bot/logic/common"
	commonLogicDomain "heroku-line-bot/logic/common/domain"
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

func (b *confirmRegister) Init(context domain.ICmdHandlerContext) error {
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

func (b *confirmRegister) LoadSingleParam(attr, text string) error {
	switch attr {
	case "date":
		t, err := time.Parse(commonLogicDomain.DATE_TIME_RFC3339_FORMAT, text)
		if err != nil {
			return err
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

func (b *confirmRegister) ComfirmDb() (resultErr error) {
	isChangeRole := b.isMemberAble() && b.User.Role == domain.GUEST_CLUB_ROLE

	transaction := database.Club.Begin()
	if err := transaction.Error; err != nil {
		return err
	}
	defer func() {
		if resultErr == nil {
			if resultErr = transaction.Commit().Error; resultErr != nil {
				return
			}

			if isChangeRole {
				if _, err := redis.LineUser.Del(b.User.LineID); err != nil {
					if resultErr == nil {
						resultErr = err
					}
					return
				}
			}
		}

		if err := transaction.Rollback().Error; err != nil {
			if resultErr == nil {
				resultErr = err
			}
			return
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
	if resultErr = database.Club.Member.Update(transaction, arg, fields); resultErr != nil {
		return
	}

	return nil
}

func (b *confirmRegister) Do(text string) (resultErr error) {
	lineID := b.context.GetUserID()
	if user, err := lineUserLogic.Get(lineID); err != nil {
		return err
	} else if user.Role != domain.ADMIN_CLUB_ROLE {
		return domain.NO_AUTH_ERROR
	}

	if b.User == nil {
		arg := dbReqs.Member{
			ID: &b.MemberID,
		}
		if confirmRegisterUsers, err := b.LoadUsers(arg); err != nil {
			return err
		} else if len(confirmRegisterUsers) == 0 {
			return fmt.Errorf("查無用戶")
		} else {
			v := confirmRegisterUsers[0]
			b.User = v
		}
	}

	if b.context.IsComfirmed() {
		if err := b.ComfirmDb(); err != nil {
			return err
		}

		replyMessges := []interface{}{
			linebot.GetTextMessage("完成"),
		}
		if resultErr = b.context.Reply(replyMessges); resultErr != nil {
			return resultErr
		}

		if resultErr = b.context.DeleteParam(); resultErr != nil {
			return
		}

		return nil
	}

	if err := b.context.CacheParams(); err != nil {
		return err
	}

	replyMessges, err := b.getTemplateMessage()
	if err != nil {
		return err
	}

	if err := b.context.Reply(replyMessges); err != nil {
		return err
	}

	return nil
}

func (b *confirmRegister) isMemberAble() bool {
	return b.User.CompanyID != nil &&
		b.User.Role == domain.GUEST_CLUB_ROLE &&
		b.User.Department.IsClubMember()
}

func (b *confirmRegister) getTemplateMessage() ([]interface{}, error) {
	if b.User == nil {
		return nil, nil
	}

	contents := []interface{}{}
	size := linebotDomain.MD_FLEX_MESSAGE_SIZE
	keyValueEditComponentOption := &domain.KeyValueEditComponentOption{
		SizeP: &size,
	}

	if js, err := b.context.
		GetDateTimeCmdInputMode(domain.DATE_POSTBACK_DATE_TIME_CMD, "date").
		GetSignal(); err != nil {
		return nil, err
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
				fmt.Sprintf("%s(%s)", b.Date.Format(commonLogicDomain.DATE_FORMAT), commonLogic.WeekDayName(b.Date.Weekday())),
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

	comfirmSignlJs, err := b.context.
		GetComfirmMode().
		GetSignal()
	if err != nil {
		return nil, err
	}
	contents = append(contents,
		linebot.GetButtonComponent(
			0,
			linebot.GetPostBackAction(
				"確認",
				comfirmSignlJs,
			),
			nil,
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
