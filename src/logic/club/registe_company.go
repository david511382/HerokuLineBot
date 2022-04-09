package club

import (
	"fmt"
	"heroku-line-bot/src/logger"
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

	"github.com/rs/zerolog"
)

type registeCompany struct {
	context                   domain.ICmdHandlerContext `json:"-"`
	Department1               *domain.Department        `json:"department_1"`
	Department2               *string                   `json:"department_2"`
	Department3               *string                   `json:"department_3"`
	CompanyID                 *string                   `json:"company_id"`
	IsRequireDbCheckCompanyID bool                      `json:"is_require_db_check_company_id"`
	MemberID                  *int                      `json:"member_id"`
}

func (b *registeCompany) Init(context domain.ICmdHandlerContext) (resultErrInfo errUtil.IError) {
	*b = registeCompany{
		context: context,
	}

	return nil
}

func (b *registeCompany) GetRequireAttr() (requireAttr string, warnMessage interface{}, resultErrInfo errUtil.IError) {
	b.loadMemberInfo()

	if b.CompanyID == nil {
		requireAttr = "company_id"
		return
	}

	if b.IsRequireDbCheckCompanyID {
		if count, err := database.Club.Member.Count(
			dbModel.ReqsClubMember{
				CompanyID: b.CompanyID,
			},
		); err != nil {
			errInfo := errUtil.NewError(err)
			resultErrInfo = errUtil.Append(resultErrInfo, errInfo)
			return
		} else if count > 0 {
			requireAttr = "company_id"
			warnMessage = linebot.GetTextMessage("員工編號已被使用")
			b.IsRequireDbCheckCompanyID = false
			if errInfo := b.context.CacheParams(); errInfo != nil {
				resultErrInfo = errUtil.Append(resultErrInfo, errInfo)
				return
			}
			return
		}

		if 處 := b.Department1; 處 == nil {
			// 是第一次輸入而不是修改 CompanyID 時
			requireAttr = "處"
		}

		b.IsRequireDbCheckCompanyID = false
		if errInfo := b.context.CacheParams(); errInfo != nil {
			resultErrInfo = errUtil.Append(resultErrInfo, errInfo)
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

func (b *registeCompany) GetRequireAttrInfo(rawAttr string) (attrNameText string, valueText string, isNotRequireChecking bool) {
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
	}
	return
}

func (b *registeCompany) GetInputTemplate(attr string) (messages interface{}) {
	switch attr {
	case "company_id":
		const altText = "請確認或輸入"
		bodyComponents := []interface{}{
			linebot.GetTextMessage("成為社員必須要員工編號喔！"),
		}
		titleMessages := []interface{}{}

		cancelSignlJs, errInfo := NewSignal().
			GetCancelMode().
			GetSignal()
		if errInfo != nil {
			logger.LogError(logger.NAME_LINE, errInfo)
			return
		}

		text := "請輸入員工編號"
		if 已經輸入 := b.CompanyID != nil; 已經輸入 {
			text = fmt.Sprintf("確認員工編號為: %s", *b.CompanyID)

			pathValueMap := map[string]interface{}{
				"ICmdLogic.is_require_db_check_company_id": true,
			}
			checkCompanyIDJs, errInfo := NewSignal().
				GetCancelInputMode().
				GetKeyValueInputMode(pathValueMap).
				GetSignal()
			if errInfo != nil {
				logger.LogError(logger.NAME_LINE, errInfo)
				return
			}
			comfirmButtonComponent := GetConfirmComponent(
				linebot.GetPostBackAction(
					"取消",
					cancelSignlJs,
				),
				linebot.GetPostBackAction(
					"確認",
					checkCompanyIDJs,
				),
			)
			bodyComponents = append(bodyComponents, comfirmButtonComponent)
		} else {
			cancelButtonComponent := linebot.GetButtonComponent(
				linebot.GetPostBackAction(
					"取消",
					cancelSignlJs,
				),
				&domain.NormalButtonOption,
			)
			bodyComponents = append(bodyComponents, cancelButtonComponent)
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
				logger.LogError(logger.NAME_LINE, errInfo)
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
			logger.LogError(logger.NAME_LINE, errInfo)
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

func (b *registeCompany) LoadRequireInputTextParam(attr, text string) (resultErrInfo errUtil.IError) {
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
	default:
	}
	return nil
}

func (b *registeCompany) loadMemberInfo() (resultErrInfo errUtil.IError) {
	if b.MemberID == nil {
		lineID := b.context.GetUserID()
		if dbDatas, err := database.Club.Member.Select(
			dbModel.ReqsClubMember{
				LineID: &lineID,
			},
			memberDb.COLUMN_ID,
			memberDb.COLUMN_Department,
			memberDb.COLUMN_CompanyID,
		); err != nil {
			errInfo := errUtil.NewError(err)
			resultErrInfo = errUtil.Append(resultErrInfo, errInfo)
			return
		} else if len(dbDatas) == 0 {
			name := b.context.GetUserName()
			registerMember := NewRegisterMember(name, &lineID)
			if errInfo := registerMember.Registe(nil); errInfo != nil {
				resultErrInfo = errUtil.Append(resultErrInfo, errInfo)
				if resultErrInfo.IsError() {
					return
				}
			}

			mID, _, errInfo := registerMember.LoadMemberID()
			if errInfo != nil {
				resultErrInfo = errUtil.Append(resultErrInfo, errInfo)
				if resultErrInfo.IsError() {
					return
				}
			}

			b.MemberID = mID
			if errInfo := b.context.CacheParams(); errInfo != nil {
				resultErrInfo = errUtil.Append(resultErrInfo, errInfo)
				if resultErrInfo.IsError() {
					return
				}
			}
		} else {
			dbData := dbDatas[0]
			b.MemberID = &dbData.ID
			b.CompanyID = dbData.CompanyID
			if department := Department(dbData.Department); department != NewEmptyDepartment() {
				處, 部, 組 := department.Split()
				b.Department1 = &處
				b.Department2 = &部
				b.Department3 = &組
			}
		}
	}
	return
}

func (b *registeCompany) Do(text string) (resultErrInfo errUtil.IError) {
	b.loadMemberInfo()

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

		arg := dbModel.ReqsClubMember{
			ID: b.MemberID,
		}
		fields := map[string]interface{}{
			"department": string(b.getDepartment()),
			"company_id": b.CompanyID,
		}
		if err := db.Member.Update(arg, fields); err != nil {
			resultErrInfo = errUtil.NewError(err)
			return
		}

		{
			if dbDatas, err := db.Member.Select(
				dbModel.ReqsClubMember{
					ID: b.MemberID,
				},
				memberDb.COLUMN_Name,
			); err != nil {
				errInfo := errUtil.NewError(err, zerolog.WarnLevel)
				resultErrInfo = errUtil.Append(resultErrInfo, errInfo)
			} else if len(dbDatas) == 0 {
				errInfo := errUtil.New("資料異常")
				resultErrInfo = errUtil.Append(resultErrInfo, errInfo)
				return
			} else {
				name := dbDatas[0].Name
				if adminReplyContents, err := b.GetNotifyRegisterContents(name); err != nil {
					resultErrInfo = errUtil.NewError(err)
					return
				} else {
					adminReplyMessges := []interface{}{
						linebot.GetFlexMessage(
							"公司新人註冊",
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
					if err := b.context.PushAdmin(adminReplyMessges); err != nil {
						resultErrInfo = errUtil.NewError(err)
						return
					}
				}
			}
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
		resultErrInfo = errUtil.Append(resultErrInfo, errInfo)
		return
	}

	contents := []interface{}{}
	size := linebotDomain.MD_FLEX_MESSAGE_SIZE
	if b.CompanyID != nil {
		if js, errInfo := NewSignal().
			GetRequireInputMode("company_id").
			GetSignal(); errInfo != nil {
			resultErrInfo = errUtil.Append(resultErrInfo, errInfo)
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
		resultErrInfo = errUtil.Append(resultErrInfo, errInfo)
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
		resultErrInfo = errUtil.Append(resultErrInfo, errInfo)
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
		resultErrInfo = errUtil.Append(resultErrInfo, errInfo)
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

func (b *registeCompany) GetNotifyRegisterContents(name string) ([]interface{}, error) {
	contents := []interface{}{}

	contents = append(contents,
		GetKeyValueEditComponent(
			"暱稱",
			name,
			nil,
		),
	)
	if b.CompanyID != nil {
		contents = append(contents,
			GetKeyValueEditComponent(
				"員工編號",
				*b.CompanyID,
				nil,
			),
		)
	}

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

func (b *registeCompany) get處() (處 domain.Department) {
	if p := b.Department1; p == nil || *p == "" {
		處 = DEPARTMENT_NONE
	} else {
		處 = *p
	}
	return
}

func (b *registeCompany) get部() (部 string) {
	if p := b.Department2; p == nil || *p == "" {
		部 = DEPARTMENT_NONE
	} else {
		部 = *p
	}
	return
}

func (b *registeCompany) get組() (組 string) {
	if p := b.Department3; p == nil || *p == "" {
		組 = DEPARTMENT_NONE
	} else {
		組 = *p
	}
	return
}

func (b *registeCompany) getDepartment() Department {
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
