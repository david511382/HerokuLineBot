package club

import (
	"fmt"
	"heroku-line-bot/logic/club/domain"
	commonLogic "heroku-line-bot/logic/common"
	"heroku-line-bot/service/linebot"
	linebotDomain "heroku-line-bot/service/linebot/domain"
	"heroku-line-bot/service/linebot/domain/model"
	"heroku-line-bot/storage/database"
	memberDb "heroku-line-bot/storage/database/database/clubdb/table/member"
	dbReqs "heroku-line-bot/storage/database/domain/model/reqs"
	"heroku-line-bot/util"
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
	context       domain.ICmdHandlerContext `json:"-"`
	Department    Department                `json:"department"`
	Name          string                    `json:"name"`
	CompanyID     *string                   `json:"company_id"`
	isAlreadyExit bool
}

func (b *register) Init(context domain.ICmdHandlerContext, initCmdBaseF func(requireRawParamAttr, requireRawParamAttrText string, isInputImmediately bool)) error {
	lineID := context.GetUserID()
	isAlreadyExit := false
	arg := dbReqs.Member{
		LineID: &lineID,
	}
	if count, err := database.Club.Member.Count(arg); err != nil {
		return err
	} else if count > 0 {
		isAlreadyExit = true
	} else {
		initCmdBaseF(
			"company_id",
			"員工編號",
			false,
		)
	}

	*b = register{
		context:       context,
		isAlreadyExit: isAlreadyExit,
	}

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

func (b *register) LoadSingleParam(attr, text string) error {
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
		buttonComponents := []interface{}{}
		titleMessages := []interface{}{}

		inputDepartmentJs, err := b.context.GetRequireInputCmdText(nil, "處", "處", false)
		if err != nil {
			return err
		}

		text1 := "請輸入員工編號"
		if b.CompanyID != nil {
			text1 = fmt.Sprintf("確認員工編號為: %s ,或繼續輸入", *b.CompanyID)

			comfirmButtonComponent := linebot.GetButtonComponent(
				0,
				linebot.GetPostBackAction(
					"確認",
					inputDepartmentJs,
				),
			)
			buttonComponents = append(buttonComponents, comfirmButtonComponent)
		}
		titleMessages = append(titleMessages, linebot.GetTextMessage(text1))
		if b.CompanyID == nil {
			const text2 = "成為社員必須要員工編號喔！"
			titleMessages = append(titleMessages, linebot.GetTextMessage(text2))
		}

		buttonComponents = append(buttonComponents,
			linebot.GetButtonComponent(
				0,
				linebot.GetPostBackAction(
					"沒有員工編號",
					inputDepartmentJs,
				),
			),
		)

		return linebot.GetFlexMessage(
			altText,
			linebot.GetFlexMessageBubbleContent(
				linebot.GetFlexMessageBoxComponent(
					linebotDomain.VERTICAL_MESSAGE_LAYOUT,
					&model.FlexMessageBoxComponentOption{
						JustifyContent: linebotDomain.SPACE_EVENLY_JUSTIFY_CONTENT,
					},
					buttonComponents...,
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

			comfirmInputJs, err := b.context.GetRequireInputCmdText(nil, "部", "部門", false)
			if err != nil {
				return err
			}
			if requireRawParamAttr == "處single" {
				comfirmInputJs, err = b.context.GetCancelInpuingSignl()
				if err != nil {
					return err
				}
			}

			comfirmButton := linebot.GetButtonComponent(
				0,
				linebot.GetPostBackAction(
					"確認",
					comfirmInputJs,
				),
			)
			inputButtons = append(inputButtons, comfirmButton)
		}
		titleMessage := linebot.GetTextMessage(text)

		for _, clubMemberDepartment := range domain.ClubMemberDepartments {
			departmentButton := linebot.GetButtonComponent(
				0,
				linebot.GetMessageAction(string(clubMemberDepartment)),
			)
			inputButtons = append(inputButtons, departmentButton)
		}
		noDepartmentButton := linebot.GetButtonComponent(
			0,
			linebot.GetMessageAction("無"),
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

		requireDepartmentInputJs, err := b.context.GetRequireInputCmdText(nil, "組", "組", false)
		if err != nil {
			return err
		}
		comfirmButton := linebot.GetButtonComponent(
			0,
			linebot.GetPostBackAction(
				"確認",
				requireDepartmentInputJs,
			),
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

func (b *register) Do(text string) (resultErr error) {
	if b.isAlreadyExit {
		replyMessges := []interface{}{
			linebot.GetTextMessage("您已經註冊過了!"),
		}
		if resultErr = b.context.Reply(replyMessges); resultErr != nil {
			return resultErr
		}

		return nil
	}

	b.init()

	if b.context.IsComfirmed() {
		transaction := database.Club.Begin()
		if err := transaction.Error; err != nil {
			return err
		}
		defer func() {
			if resultErr == nil {
				if resultErr = transaction.Commit().Error; resultErr != nil {
					return
				}
			}

			if resultErr = transaction.Rollback().Error; resultErr != nil {
				return
			}
		}()

		role := domain.GUEST_CLUB_ROLE
		if b.CompanyID != nil && b.Department.IsClubMember() {
			role = domain.MEMBER_CLUB_ROLE
		}
		nowTime := commonLogic.TimeUtilObj.Now()
		data := &memberDb.MemberTable{
			JoinDate:   util.DateOf(nowTime),
			Department: string(b.Department),
			Name:       b.Name,
			CompanyID:  b.CompanyID,
			Role:       int16(role),
			LineID:     util.GetStringP(b.context.GetUserID()),
		}
		if resultErr = database.Club.Member.BaseTable.Insert(transaction, data); resultErr != nil {
			return
		}

		if resultErr = b.context.DeleteParam(); resultErr != nil {
			return
		}

		replyMessges := []interface{}{
			linebot.GetTextMessage("完成"),
		}
		if resultErr = b.context.Reply(replyMessges); resultErr != nil {
			return resultErr
		}

		return nil
	}

	if err := b.context.CacheParams(); err != nil {
		return err
	}

	contents := []interface{}{}
	if b.CompanyID != nil {
		if js, err := b.context.GetRequireInputCmdText(nil, "company_id", "員工編號", false); err != nil {
			return err
		} else {
			action := linebot.GetPostBackAction(
				"修改",
				js,
			)
			contents = append(contents,
				linebot.GetKeyValueEditComponent(
					"員工編號",
					*b.CompanyID,
					action,
					nil, nil,
				),
			)
		}
	}

	處, 部, 組 := b.Department.Split()
	if 處 == "" {
		處 = "無"
	}
	if js, err := b.context.GetRequireInputCmdText(nil, "處single", "處", false); err != nil {
		return err
	} else {
		action := linebot.GetPostBackAction(
			"修改",
			js,
		)
		contents = append(contents,
			linebot.GetKeyValueEditComponent(
				"處",
				string(處),
				action,
				nil, nil,
			),
		)
	}
	if 部 == "" {
		部 = "無"
	}
	if js, err := b.context.GetRequireInputCmdText(nil, "部single", "部門", false); err != nil {
		return err
	} else {
		action := linebot.GetPostBackAction(
			"修改",
			js,
		)
		contents = append(contents,
			linebot.GetKeyValueEditComponent(
				"部",
				部,
				action,
				nil, nil,
			),
		)
	}
	if 組 == "" {
		組 = "無"
	}
	if js, err := b.context.GetRequireInputCmdText(nil, "組single", "組", false); err != nil {
		return err
	} else {
		action := linebot.GetPostBackAction(
			"修改",
			js,
		)
		contents = append(contents,
			linebot.GetKeyValueEditComponent(
				"組",
				組,
				action,
				nil, nil,
			),
		)
	}

	if js, err := b.context.GetRequireInputCmdText(nil, "name", "暱稱", false); err != nil {
		return err
	} else {
		action := linebot.GetPostBackAction(
			"修改",
			js,
		)
		contents = append(contents,
			linebot.GetKeyValueEditComponent(
				"暱稱",
				b.Name,
				action,
				nil, nil,
			),
		)
	}

	cancelSignlJs, err := b.context.GetCancelSignl()
	if err != nil {
		return err
	}
	comfirmSignlJs, err := b.context.GetComfirmSignl()
	if err != nil {
		return err
	}
	contents = append(contents,
		linebot.GetComfirmComponent(
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
		return err
	}

	return nil
}
