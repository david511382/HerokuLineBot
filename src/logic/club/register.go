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
	"strings"
)

type Department string

func (d Department) Split() (處 domain.Department, 部, 組 string) {
	ds := strings.Split(string(d), "/")
	if len(ds) >= 1 {
		處 = domain.Department(ds[0])
	}
	if len(ds) >= 2 {
		部 = ds[1]
	}
	if len(ds) >= 3 {
		組 = ds[2]
	}
	return
}

func (d Department) IsClubMember() bool {
	處, _, _ := d.Split()
	for _, clubMemberDepartment := range domain.ClubMemberDepartments {
		if 處 == clubMemberDepartment {
			return true
		}
	}
	return false
}

func (d *Department) Set處(data domain.Department) {
	_, 部, 組 := d.Split()
	d.set(data, 部, 組)
}

func (d *Department) Set部(data string) {
	處, _, 組 := d.Split()
	d.set(處, data, 組)
}

func (d *Department) Set組(data string) {
	處, 部, _ := d.Split()
	d.set(處, 部, data)
}

func (d *Department) set(處 domain.Department, 部, 組 string) {
	strs := []string{
		string(處), 部, 組,
	}
	*d = Department(strings.Join(strs, "/"))
}

type register struct {
	context                   domain.ICmdHandlerContext `json:"-"`
	Department                Department                `json:"department"`
	Name                      string                    `json:"name"`
	CompanyID                 *string                   `json:"company_id"`
	IsRequireDbCheckCompanyID bool                      `json:"is_require_db_check_company_id"`
	Role                      domain.ClubRole           `json:"role"`
	MemberID                  int                       `json:"member_id"`
}

func (b *register) Init(context domain.ICmdHandlerContext) (resultErrInfo errUtil.IError) {
	*b = register{
		context: context,
		Role:    domain.GUEST_CLUB_ROLE,
	}

	b.context.SetRequireInputMode(
		"company_id",
		"員工編號",
		false,
	)

	return nil
}

func (b *register) GetSingleParam(attr string) string {
	switch attr {
	case "部single", "部":
		_, 部, _ := b.Department.Split()
		if 部 == "" {
			return "無"
		}
		return 部
	case "組single", "組":
		_, _, 組 := b.Department.Split()
		if 組 == "" {
			return "無"
		}
		return 組
	case "name":
		return b.Name
	default:
		return ""
	}
}

func (b *register) LoadSingleParam(attr, text string) (resultErrInfo errUtil.IError) {
	switch attr {
	case "company_id":
		b.CompanyID = &text
	case "處single", "處":
		b.Department.Set處(domain.Department(text))
	case "部single", "部":
		b.Department.Set部(text)
	case "組single", "組":
		b.Department.Set組(text)
	case "name":
		b.Name = text
	default:
	}
	return nil
}

func (b *register) GetInputTemplate(requireRawParamAttr string) interface{} {
	switch requireRawParamAttr {
	case "company_id":
		const altText = "請確認或輸入"
		bodyComponents := []interface{}{
			linebot.GetTextMessage("成為社員必須要員工編號喔！"),
		}
		titleMessages := []interface{}{}

		text1 := "請輸入員工編號"
		if b.CompanyID != nil {
			text1 = fmt.Sprintf("確認員工編號為: %s", *b.CompanyID)

			pathValueMap := map[string]interface{}{
				"ICmdLogic.is_require_db_check_company_id": true,
			}
			checkCompanyIDJs, errInfo := b.context.
				GetCancelInputMode().
				GetKeyValueInputMode(pathValueMap).
				GetSignal()
			if errInfo != nil {
				logger.Log("line bot club", errInfo)
				return nil
			}
			comfirmButtonComponent := linebot.GetButtonComponent(
				linebot.GetPostBackAction(
					"確認",
					checkCompanyIDJs,
				),
				&domain.NormalButtonOption,
			)
			bodyComponents = append(bodyComponents, comfirmButtonComponent)
		} else {
			lineID := b.context.GetUserID()
			arg := dbModel.ReqsClubMember{
				LineID: &lineID,
			}
			if count, err := database.Club.Member.Count(arg); err == nil && count > 0 {
				if err := b.context.DeleteParam(); err != nil {
					logger.Log("line bot club", errUtil.NewError(err))
					return nil
				}

				return linebot.GetFlexMessage(
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
			}
		}
		titleMessages = append(titleMessages, linebot.GetTextMessage(text1))
		if b.CompanyID != nil {
			titleMessages = append(titleMessages, linebot.GetTextMessage("或繼續輸入"))
		}

		return linebot.GetFlexMessage(
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
	case "處single", "處":
		const altText = "請選擇處"
		處, _, _ := b.Department.Split()

		inputButtons := []interface{}{}
		text := altText
		if 處 != "" {
			text = fmt.Sprintf("確認處為: %s 嗎？", 處)

			comfirmInputJs, errInfo := b.context.
				GetRequireInputMode("部", "部門", false).
				GetSignal()
			if errInfo != nil {
				logger.Log("line bot club", errInfo)
				return nil
			}
			if requireRawParamAttr == "處single" {
				comfirmInputJs, errInfo = b.context.
					GetCancelInputMode().
					GetSignal()
				if errInfo != nil {
					logger.Log("line bot club", errInfo)
					return nil
				}
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
			linebot.GetMessageAction("無"),
			&domain.NormalButtonOption,
		)
		inputButtons = append(inputButtons, noDepartmentButton)

		return linebot.GetFlexMessage(
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
	case "部":
		const altText = "請選擇部"
		_, 部, _ := b.Department.Split()

		text := "請輸入部門,若無請直接確認"
		if 部 != "" {
			text = fmt.Sprintf("確認部為: %s 嗎？", 部)
		}

		requireDepartmentInputJs, errInfo := b.context.
			GetRequireInputMode("組", "組", false).
			GetSignal()
		if errInfo != nil {
			logger.Log("line bot club", errInfo)
			return nil
		}
		comfirmButton := linebot.GetButtonComponent(
			linebot.GetPostBackAction(
				"確認",
				requireDepartmentInputJs,
			),
			&domain.NormalButtonOption,
		)

		return linebot.GetFlexMessage(
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
	default:
		return nil
	}
}

func (b *register) init() {
	if b.Name == "" {
		b.Name = b.context.GetUserName()
	}
}

func (b *register) Do(text string) (resultErrInfo errUtil.IError) {
	if b.IsRequireDbCheckCompanyID {
		b.IsRequireDbCheckCompanyID = false
		b.MemberID = 0
		if errInfo := b.context.CacheParams(); errInfo != nil {
			resultErrInfo = errInfo
			return
		}
		arg := dbModel.ReqsClubMember{
			CompanyID: b.CompanyID,
		}
		if dbDatas, err := database.Club.Member.Select(
			arg,
			memberDb.COLUMN_ID,
			memberDb.COLUMN_Name,
			memberDb.COLUMN_Role,
			memberDb.COLUMN_Department,
		); err != nil {
			resultErrInfo = errUtil.NewError(err)
			return
		} else if len(dbDatas) == 0 {
			if 處, _, _ := b.Department.Split(); 處 == "" {
				b.context.SetRequireInputMode("處", "處", false)
				if errInfo := b.context.CacheParams(); errInfo != nil {
					resultErrInfo = errInfo
					return
				}

				replyMessge := b.GetInputTemplate("處")
				replyMessges := []interface{}{
					replyMessge,
				}
				if err := b.context.Reply(replyMessges); err != nil {
					resultErrInfo = errUtil.NewError(err)
					return
				}

				return nil
			}
		} else {
			v := dbDatas[0]
			if v.Name != "" {
				b.Name = v.Name
			}

			b.MemberID = v.ID
			b.Department = Department(v.Department)
			b.Role = domain.ClubRole(v.Role)
		}
	}

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
		if b.MemberID > 0 {
			arg := dbModel.ReqsClubMember{
				ID: &b.MemberID,
			}
			fields := map[string]interface{}{
				"department": string(b.Department),
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
				Department: string(b.Department),
				Name:       b.Name,
				CompanyID:  b.CompanyID,
				Role:       int16(b.Role),
				LineID:     &lineID,
			}
			if errInfo := account.Registe(db, data); errInfo != nil {
				resultErrInfo = errUtil.Append(resultErrInfo, errInfo)
				return
			}
			b.MemberID = data.ID
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

		return nil
	}

	if errInfo := b.context.CacheParams(); errInfo != nil {
		resultErrInfo = errInfo
		return
	}

	contents := []interface{}{}
	size := linebotDomain.MD_FLEX_MESSAGE_SIZE
	if b.CompanyID != nil {
		if js, errInfo := b.context.
			GetRequireInputMode("company_id", "員工編號", false).
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

	處, 部, 組 := b.Department.Split()
	if 處 == "" {
		處 = "無"
	}
	if js, errInfo := b.context.
		GetRequireInputMode("處single", "處", false).
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
	if 部 == "" {
		部 = "無"
	}
	if js, errInfo := b.context.
		GetRequireInputMode("部single", "部門", false).
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
	if 組 == "" {
		組 = "無"
	}
	if js, errInfo := b.context.
		GetRequireInputMode("組single", "組", false).
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

	if js, errInfo := b.context.
		GetRequireInputMode("name", "暱稱", false).
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

	cancelSignlJs, errInfo := b.context.
		GetCancelMode().
		GetSignal()
	if errInfo != nil {
		resultErrInfo = errInfo
		return
	}
	comfirmSignlJs, errInfo := b.context.
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
	if js, errInfo := b.context.
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
