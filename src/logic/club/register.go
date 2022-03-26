package club

import (
	"fmt"
	"heroku-line-bot/src/logger"
	"heroku-line-bot/src/logic/account"
	"heroku-line-bot/src/logic/club/domain"
	dbModel "heroku-line-bot/src/model/database"
	"heroku-line-bot/src/pkg/global"
	"heroku-line-bot/src/pkg/service/linebot"
	linebotDomain "heroku-line-bot/src/pkg/service/linebot/domain"
	"heroku-line-bot/src/pkg/service/linebot/domain/model"
	"heroku-line-bot/src/pkg/util"
	errUtil "heroku-line-bot/src/pkg/util/error"
	"heroku-line-bot/src/repo/database"
	memberDb "heroku-line-bot/src/repo/database/database/clubdb/member"
)

type register struct {
	context                   domain.ICmdHandlerContext `json:"-"`
	Department1               *domain.Department        `json:"department_1"`
	Department2               *string                   `json:"department_2"`
	Department3               *string                   `json:"department_3"`
	Name                      string                    `json:"name"`
	CompanyID                 *string                   `json:"company_id"`
	IsRequireDbCheckCompanyID bool                      `json:"is_require_db_check_company_id"`
	//IsCompany                 *bool                     `json:"is_company"`
	Role     domain.ClubRole `json:"role"`
	MemberID *int            `json:"member_id"`
}

func (b *register) Init(context domain.ICmdHandlerContext) (resultErrInfo errUtil.IError) {
	*b = register{
		context: context,
		Role:    domain.GUEST_CLUB_ROLE,
	}

	return nil
}

func (b *register) GetRequireAttr() (requireAttr string, resultErrInfo errUtil.IError) {
	// if b.IsCompany == nil {
	// 	requireAttr = "is_company"
	// 	return
	// }

	// isCompany := *b.IsCompany
	// if !isCompany {
	// 	return
	// }

	if b.CompanyID == nil {
		requireAttr = "company_id"
		return
	}

	if b.IsRequireDbCheckCompanyID {
		requireAttr = ""

		dbDatas, err := database.Club.Member.Select(
			dbModel.ReqsClubMember{
				CompanyID: b.CompanyID,
			},
			memberDb.COLUMN_ID,
			memberDb.COLUMN_Name,
			memberDb.COLUMN_Role,
			memberDb.COLUMN_Department,
		)
		if err != nil {
			resultErrInfo = errUtil.NewError(err)
			return
		} else if len(dbDatas) == 0 {
			if 處 := b.Department1; 處 == nil {
				// 是第一次輸入而不是修改 CompanyID 時
				if errInfo := b.context.CacheParams(); errInfo != nil {
					resultErrInfo = errInfo
					return
				}

				requireAttr = "處"
			}
			b.MemberID = nil
		} else {
			v := dbDatas[0]
			if v.Name != "" {
				b.Name = v.Name
			}

			b.MemberID = util.GetIntP(v.ID)
			處, 部, 組 := Department(v.Department).Split()
			b.Department1 = &處
			b.Department2 = &部
			b.Department3 = &組
			b.Role = domain.ClubRole(v.Role)
		}

		b.IsRequireDbCheckCompanyID = false
		if errInfo := b.context.CacheParams(); errInfo != nil {
			resultErrInfo = errInfo
			return
		}
		return
	}

	{
		if b.Department1 == nil {
			requireAttr = "處"
			return
		} else if *b.Department1 != "" {
			if b.Department2 == nil {
				requireAttr = "部"
				return
			} else if *b.Department2 != "" {
				if b.Department3 == nil {
					requireAttr = "組"
					return
				}
			}
		}
	}

	return
}

func (b *register) GetRequireAttrInfo(rawAttr string) (attrNameText string, valueText string, isNotRequireChecking bool) {
	switch rawAttr {
	case "處":
		attrNameText = "處"

		處 := b.get處()
		valueText = string(處)
	case "部":
		attrNameText = "部門"

		valueText = b.get部()
	case "組":
		attrNameText = "組"

		valueText = b.get組()
	case "company_id":
		attrNameText = "員工編號"
	case "name":
		attrNameText = "暱稱"
		valueText = b.Name
	}
	return
}

func (b *register) GetInputTemplate(attr string) (messages interface{}) {
	switch attr {
	case "company_id":
		const altText = "請確認或輸入"
		bodyComponents := []interface{}{
			linebot.GetTextMessage("成為社員必須要員工編號喔！"),
		}
		titleMessages := []interface{}{}

		text := "請輸入員工編號"
		if 還沒輸入 := b.CompanyID == nil; 還沒輸入 {
			lineID := b.context.GetUserID()
			arg := dbModel.ReqsClubMember{
				LineID: &lineID,
			}
			if count, err := database.Club.Member.Count(arg); err == nil && count > 0 {
				if err := b.context.DeleteParam(); err != nil {
					logger.Log("line bot club", errUtil.NewError(err))
					return
				}

				messages = linebot.GetFlexMessage(
					"通知",
					linebot.GetFlexMessageBubbleContent(
						linebot.GetFlexMessageBoxComponent(
							linebotDomain.VERTICAL_MESSAGE_LAYOUT,
							&model.FlexMessageBoxComponentOption{
								JustifyContent: linebotDomain.SPACE_EVENLY_JUSTIFY_CONTENT,
							},
							linebot.GetTextMessage("您已經註冊過了!"),
						),
						nil,
					),
				)
				return
			}
		} else {
			text = fmt.Sprintf("確認員工編號為: %s", *b.CompanyID)

			pathValueMap := map[string]interface{}{
				"ICmdLogic.is_require_db_check_company_id": true,
			}
			checkCompanyIDJs, errInfo := NewSignal().
				GetCancelInputMode().
				GetKeyValueInputMode(pathValueMap).
				GetSignal()
			if errInfo != nil {
				logger.Log("line bot club", errInfo)
				return
			}
			comfirmButtonComponent := linebot.GetButtonComponent(
				linebot.GetPostBackAction(
					"確認",
					checkCompanyIDJs,
				),
				&domain.NormalButtonOption,
			)
			bodyComponents = append(bodyComponents, comfirmButtonComponent)
		}
		titleMessages = append(titleMessages, linebot.GetTextMessage(text))
		if b.CompanyID != nil {
			titleMessages = append(titleMessages, linebot.GetTextMessage("或繼續輸入"))
		}

		messages = linebot.GetFlexMessage(
			altText,
			linebot.GetFlexMessageBubbleContent(
				linebot.GetFlexMessageBoxComponent(
					linebotDomain.VERTICAL_MESSAGE_LAYOUT,
					&model.FlexMessageBoxComponentOption{
						JustifyContent: linebotDomain.SPACE_EVENLY_JUSTIFY_CONTENT,
					},
					bodyComponents...,
				),
				&model.FlexMessagBubbleComponentOption{
					Header: linebot.GetFlexMessageBoxComponent(
						linebotDomain.VERTICAL_MESSAGE_LAYOUT,
						nil,
						titleMessages...,
					),
					Styles: &model.FlexMessagBubbleComponentStyle{
						Header: &model.Background{
							BackgroundColor: "#8DFF33",
						},
						Body: &model.Background{
							BackgroundColor: "#FFFFFF",
							SeparatorColor:  "#000000",
							Separator:       true,
						},
					},
				},
			),
		)
	case "處":
		const altText = "請選擇處"

		inputButtons := []interface{}{}
		text := altText
		if b.Department1 != nil {
			text = fmt.Sprintf("確認處為: %s 嗎？", string(b.get處()))

			comfirmInputJs, errInfo := NewSignal().
				GetCancelInputMode().
				GetSignal()
			if errInfo != nil {
				logger.Log("line bot club", errInfo)
				return
			}
			comfirmButton := linebot.GetButtonComponent(
				linebot.GetPostBackAction(
					"確認",
					comfirmInputJs,
				),
				&domain.NormalButtonOption,
			)
			inputButtons = append(inputButtons, comfirmButton)
		}
		titleMessage := linebot.GetTextMessage(text)

		for _, clubMemberDepartment := range domain.ClubMemberDepartments {
			departmentButton := linebot.GetButtonComponent(
				linebot.GetMessageAction(string(clubMemberDepartment)),
				&domain.NormalButtonOption,
			)
			inputButtons = append(inputButtons, departmentButton)
		}
		noDepartmentButton := linebot.GetButtonComponent(
			linebot.GetMessageAction(DEPARTMENT_NONE),
			&domain.NormalButtonOption,
		)
		inputButtons = append(inputButtons, noDepartmentButton)

		messages = linebot.GetFlexMessage(
			altText,
			linebot.GetFlexMessageBubbleContent(
				linebot.GetFlexMessageBoxComponent(
					linebotDomain.VERTICAL_MESSAGE_LAYOUT,
					&model.FlexMessageBoxComponentOption{
						JustifyContent: linebotDomain.SPACE_EVENLY_JUSTIFY_CONTENT,
					},
					inputButtons...,
				),
				&model.FlexMessagBubbleComponentOption{
					Header: linebot.GetFlexMessageBoxComponent(
						linebotDomain.VERTICAL_MESSAGE_LAYOUT,
						nil,
						titleMessage,
					),
					Styles: &model.FlexMessagBubbleComponentStyle{
						Header: &model.Background{
							BackgroundColor: "#8DFF33",
						},
						Body: &model.Background{
							BackgroundColor: "#FFFFFF",
							SeparatorColor:  "#000000",
							Separator:       true,
						},
					},
				},
			),
		)
	case "部", "組":
		pathValueMap := make(map[string]interface{})
		var p *string
		if attr == "部" {
			p = b.Department2
			pathValueMap["ICmdLogic.department_2"] = ""
		} else {
			p = b.Department3
			pathValueMap["ICmdLogic.department_3"] = ""
		}
		if p != nil {
			return
		}

		altText := "請選擇" + attr
		text := fmt.Sprintf("請輸入%s,若無請直接確認", attr)

		requireDepartmentInputJs, errInfo := NewSignal().
			GetCancelInputMode().
			GetKeyValueInputMode(pathValueMap).
			GetSignal()
		if errInfo != nil {
			logger.Log("line bot club", errInfo)
			return
		}
		comfirmButton := linebot.GetButtonComponent(
			linebot.GetPostBackAction(
				"確認",
				requireDepartmentInputJs,
			),
			&domain.NormalButtonOption,
		)

		messages = linebot.GetFlexMessage(
			altText,
			linebot.GetFlexMessageBubbleContent(
				linebot.GetFlexMessageBoxComponent(
					linebotDomain.VERTICAL_MESSAGE_LAYOUT,
					&model.FlexMessageBoxComponentOption{
						JustifyContent: linebotDomain.SPACE_EVENLY_JUSTIFY_CONTENT,
					},
					comfirmButton,
				),
				&model.FlexMessagBubbleComponentOption{
					Header: linebot.GetFlexMessageBoxComponent(
						linebotDomain.VERTICAL_MESSAGE_LAYOUT,
						nil,
						linebot.GetTextMessage(text),
					),
					Styles: &model.FlexMessagBubbleComponentStyle{
						Header: &model.Background{
							BackgroundColor: "#8DFF33",
						},
						Body: &model.Background{
							BackgroundColor: "#FFFFFF",
							SeparatorColor:  "#000000",
							Separator:       true,
						},
					},
				},
			),
		)
	}
	return
}

func (b *register) LoadRequireInputTextParam(attr, text string) (resultErrInfo errUtil.IError) {
	switch attr {
	case "company_id":
		b.CompanyID = &text
	case "處":
		if text == DEPARTMENT_NONE {
			text = ""
		}
		d := domain.Department(text)
		b.Department1 = &d
	case "部":
		if text == DEPARTMENT_NONE {
			text = ""
		}
		b.Department2 = &text
	case "組":
		if text == DEPARTMENT_NONE {
			text = ""
		}
		b.Department3 = &text
	case "name":
		b.Name = text
	default:
	}
	return nil
}

func (b *register) init() {
	if b.Name == "" {
		b.Name = b.context.GetUserName()
	}
}

func (b *register) Do(text string) (resultErrInfo errUtil.IError) {
	b.init()

	if b.context.IsConfirmed() {
		db, transaction, err := database.Club.Begin()
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
		if b.MemberID != nil {
			arg := dbModel.ReqsClubMember{
				ID: b.MemberID,
			}
			fields := map[string]interface{}{
				"department": string(b.getDepartment()),
				"name":       b.Name,
				"company_id": b.CompanyID,
				"line_id":    &lineID,
			}
			if err := db.Member.Update(arg, fields); err != nil {
				resultErrInfo = errUtil.NewError(err)
				return
			}
		} else {
			data := &dbModel.ClubMember{
				Department: string(b.getDepartment()),
				Name:       b.Name,
				CompanyID:  b.CompanyID,
				Role:       int16(b.Role),
				LineID:     &lineID,
			}
			if errInfo := account.Registe(db, data); errInfo != nil {
				resultErrInfo = errUtil.Append(resultErrInfo, errInfo)
				return
			}
			b.MemberID = util.GetIntP(data.ID)
		}

		var adminReplyMessges []interface{}
		if adminReplyContents, err := b.GetNotifyRegisterContents(); err != nil {
			resultErrInfo = errUtil.NewError(err)
			return
		} else {
			adminReplyMessges = []interface{}{
				linebot.GetFlexMessage(
					"新人註冊",
					linebot.GetFlexMessageBubbleContent(
						linebot.GetFlexMessageBoxComponent(
							linebotDomain.VERTICAL_MESSAGE_LAYOUT,
							nil,
							adminReplyContents...,
						),
						nil,
					),
				),
			}
		}

		if err := b.context.PushAdmin(adminReplyMessges); err != nil {
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

	if errInfo := b.context.CacheParams(); errInfo != nil {
		resultErrInfo = errInfo
		return
	}

	contents := []interface{}{}
	size := linebotDomain.MD_FLEX_MESSAGE_SIZE
	if b.CompanyID != nil {
		if js, errInfo := NewSignal().
			GetRequireInputMode("company_id").
			GetSignal(); errInfo != nil {
			resultErrInfo = errInfo
			return
		} else {
			action := linebot.GetPostBackAction(
				"修改",
				js,
			)
			contents = append(contents,
				GetKeyValueEditComponent(
					"員工編號",
					*b.CompanyID,
					&domain.KeyValueEditComponentOption{
						Action: action,
						SizeP:  &size,
					},
				),
			)
		}
	}

	處 := b.get處()
	if js, errInfo := NewSignal().
		GetRequireInputMode("處").
		GetSignal(); errInfo != nil {
		resultErrInfo = errInfo
		return
	} else {
		action := linebot.GetPostBackAction(
			"修改",
			js,
		)
		contents = append(contents,
			GetKeyValueEditComponent(
				"處",
				string(處),
				&domain.KeyValueEditComponentOption{
					Action: action,
					SizeP:  &size,
				},
			),
		)
	}

	部 := b.get部()
	if js, errInfo := NewSignal().
		GetRequireInputMode("部").
		GetSignal(); errInfo != nil {
		resultErrInfo = errInfo
		return
	} else {
		action := linebot.GetPostBackAction(
			"修改",
			js,
		)
		contents = append(contents,
			GetKeyValueEditComponent(
				"部",
				部,
				&domain.KeyValueEditComponentOption{
					Action: action,
					SizeP:  &size,
				},
			),
		)
	}

	組 := b.get組()
	if js, errInfo := NewSignal().
		GetRequireInputMode("組").
		GetSignal(); errInfo != nil {
		resultErrInfo = errInfo
		return
	} else {
		action := linebot.GetPostBackAction(
			"修改",
			js,
		)
		contents = append(contents,
			GetKeyValueEditComponent(
				"組",
				組,
				&domain.KeyValueEditComponentOption{
					Action: action,
					SizeP:  &size,
				},
			),
		)
	}

	if js, errInfo := NewSignal().
		GetRequireInputMode("name").
		GetSignal(); errInfo != nil {
		resultErrInfo = errInfo
		return
	} else {
		action := linebot.GetPostBackAction(
			"修改",
			js,
		)
		contents = append(contents,
			GetKeyValueEditComponent(
				"暱稱",
				b.Name,
				&domain.KeyValueEditComponentOption{
					Action: action,
					SizeP:  &size,
				},
			),
		)
	}

	cancelSignlJs, errInfo := NewSignal().
		GetCancelMode().
		GetSignal()
	if errInfo != nil {
		resultErrInfo = errInfo
		return
	}
	comfirmSignlJs, errInfo := NewSignal().
		GetConfirmMode().
		GetSignal()
	if errInfo != nil {
		resultErrInfo = errInfo
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

	replyMessges := []interface{}{
		linebot.GetFlexMessage(
			"註冊",
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

	return nil
}

func (b *register) GetNotifyRegisterContents() ([]interface{}, error) {
	contents := []interface{}{}

	contents = append(contents,
		GetKeyValueEditComponent(
			"暱稱",
			b.Name,
			nil,
		),
	)

	pathValueMap := make(map[string]interface{})
	pathValueMap["ICmdLogic.member_id"] = b.MemberID
	pathValueMap["ICmdLogic.date"] = util.DateOf(global.TimeUtilObj.Now())
	cmd := domain.CONFIRM_REGISTER_TEXT_CMD
	if js, errInfo := NewSignal().
		GetCmdInputMode(&cmd).
		GetKeyValueInputMode(pathValueMap).
		GetSignal(); errInfo != nil {
		return nil, errInfo
	} else {
		contents = append(contents,
			linebot.GetClassButtonComponent(
				linebot.GetPostBackAction(
					"入社",
					js,
				),
			),
		)
	}

	return contents, nil
}

func (b *register) get處() (處 domain.Department) {
	if p := b.Department1; p == nil || *p == "" {
		處 = DEPARTMENT_NONE
	} else {
		處 = *p
	}
	return
}

func (b *register) get部() (部 string) {
	if p := b.Department2; p == nil || *p == "" {
		部 = DEPARTMENT_NONE
	} else {
		部 = *p
	}
	return
}

func (b *register) get組() (組 string) {
	if p := b.Department3; p == nil || *p == "" {
		組 = DEPARTMENT_NONE
	} else {
		組 = *p
	}
	return
}

func (b *register) getDepartment() Department {
	var 處 domain.Department = ""
	if b.Department1 != nil {
		處 = *b.Department1
	}
	部 := ""
	if b.Department2 != nil {
		部 = *b.Department2
	}
	組 := ""
	if b.Department3 != nil {
		組 = *b.Department3
	}
	return NewDepartment(處, 部, 組)
}
