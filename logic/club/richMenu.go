package club

import (
	"fmt"
	"heroku-line-bot/logic/club/domain"
	"heroku-line-bot/logic/redis/lineuser"
	"heroku-line-bot/service/linebot"
	linebotDomain "heroku-line-bot/service/linebot/domain"
	"heroku-line-bot/service/linebot/domain/model"
	linebotReqs "heroku-line-bot/service/linebot/domain/model/reqs"
	"heroku-line-bot/storage/database"
	dbReqs "heroku-line-bot/storage/database/domain/model/reqs"
	"math"
	"strconv"
)

type richMenu struct {
	context    domain.ICmdHandlerContext `json:"-"`
	Method     domain.RichMenuMethod     `json:"method"`
	Role       domain.ClubRole           `json:"role"`
	RichMenuID string                    `json:"rich_menu_id"`
}

func (b *richMenu) Init(context domain.ICmdHandlerContext) error {
	*b = richMenu{
		context: context,
	}

	return nil
}

func (b *richMenu) GetSingleParam(attr string) string {
	switch attr {
	case "rich_menu_id":
		return b.RichMenuID
	default:
		return ""
	}
}

func (b *richMenu) LoadSingleParam(attr, text string) error {
	switch attr {
	case "rich_menu_id":
		b.RichMenuID = text
	default:
	}

	return nil
}

func (b *richMenu) GetInputTemplate(requireRawParamAttr string) interface{} {
	switch requireRawParamAttr {
	case "role":
		buttons := []interface{}{}
		roles := []domain.ClubRole{
			domain.ADMIN_CLUB_ROLE,
			domain.CADRE_CLUB_ROLE,
			domain.GUEST_CLUB_ROLE,
		}
		for _, role := range roles {
			pathValueMap := map[string]interface{}{
				"ICmdLogic.method": domain.NEW_RICH_MENU_METHOD,
				"ICmdLogic.role":   role,
			}
			if js, err := b.context.
				GetCmdInputMode(nil).
				GetCancelInputMode().
				GetKeyValueInputMode(pathValueMap).
				GetSignal(); err != nil {
				return err
			} else {
				action := linebot.GetPostBackAction(strconv.Itoa(int(role)), js)
				button := linebot.GetButtonComponent(
					0,
					action,
					&domain.NormalButtonOption,
				)
				buttons = append(buttons, button)
			}
		}
		return linebot.GetFlexMessage(
			"RichMenu Methods",
			linebot.GetFlexMessageBubbleContent(
				linebot.GetFlexMessageBoxComponent(
					linebotDomain.VERTICAL_MESSAGE_LAYOUT,
					&model.FlexMessageBoxComponentOption{
						JustifyContent: linebotDomain.SPACE_EVENLY_JUSTIFY_CONTENT,
					},
					buttons...,
				),
				nil,
			),
		)
	default:
		return nil
	}
}

func (b *richMenu) Do(text string) (resultErr error) {
	if u, err := lineuser.Get(b.context.GetUserID()); err != nil {
		return err
	} else {
		if u.Role != domain.ADMIN_CLUB_ROLE {
			return domain.NO_AUTH_ERROR
		}
	}

	messages := []interface{}{}
	switch b.Method {
	case domain.LIST_RICH_MENU_METHOD:
	case domain.DELETE_RICH_MENU_METHOD:
	case domain.NEW_RICH_MENU_METHOD:
		createRichMenuArg := b.createRoleRichMenu(b.Role)
		createRichMenuResp, err := b.context.GetBot().CreateRichMenu(createRichMenuArg)
		if err != nil {
			return err
		}

		imgBs := b.getRoleRichMenuImg(b.Role)
		if _, err := b.context.GetBot().UploadRichMenuImage(createRichMenuResp.RichMenuID, imgBs); err != nil {
			return err
		}

		messages = append(messages, linebot.GetTextMessage(
			"create done",
		))
		messages = append(messages, linebot.GetTextMessage(
			"RichMenuID:",
		))
		messages = append(messages, linebot.GetTextMessage(
			createRichMenuResp.RichMenuID,
		))

		if b.Role == domain.ADMIN_CLUB_ROLE ||
			b.Role == domain.CADRE_CLUB_ROLE {
			arg := dbReqs.Member{
				Role: (*int16)(&b.Role),
			}
			if dbDatas, err := database.Club.Member.NameLineID(arg); err != nil {
				return err
			} else if len(dbDatas) > 0 {
				lineIDs := []string{}
				names := []string{}
				for _, v := range dbDatas {
					if v.LineID == nil {
						continue
					}

					lineIDs = append(lineIDs, *v.LineID)
					names = append(names, v.Name)
				}

				for _, lineID := range lineIDs {
					if _, err := b.context.GetBot().SetRichMenuTo(createRichMenuResp.RichMenuID, lineID); err != nil {
						return err
					}
				}

				messages = append(messages, linebot.GetTextMessage(
					fmt.Sprintf("set to %v done", names),
				))
			}
		} else {
			if _, err := b.context.GetBot().SetDefaultRichMenu(createRichMenuResp.RichMenuID); err != nil {
				return err
			}
			messages = append(messages, linebot.GetTextMessage(
				"set done",
			))
		}

		if resultErr = b.context.DeleteParam(); resultErr != nil {
			return
		}

	case domain.SET_DEFAULT_RICH_MENU_METHOD:
		if _, err := b.context.GetBot().SetDefaultRichMenu(b.RichMenuID); err != nil {
			return err
		}
		messages = append(messages, linebot.GetTextMessage(
			"set done",
		))

		if resultErr = b.context.DeleteParam(); resultErr != nil {
			return
		}
	default:
		inputButtons := []interface{}{}
		methodSignalMap := map[domain.RichMenuMethod]domain.ICmdHandlerSignal{
			domain.LIST_RICH_MENU_METHOD:        nil,
			domain.DELETE_RICH_MENU_METHOD:      nil,
			domain.NEW_RICH_MENU_METHOD:         b.context.GetRequireInputMode("role", "", false),
			domain.SET_DEFAULT_RICH_MENU_METHOD: b.context.GetRequireInputMode("rich_menu_id", "rich menu id", false),
			domain.SET_TO_RICH_MENU_METHOD:      nil,
		}
		for method, signal := range methodSignalMap {
			if signal == nil {
				continue
			}

			pathValueMap := map[string]interface{}{
				"ICmdLogic.method": method,
			}
			if js, err := signal.
				GetCmdInputMode(nil).
				GetKeyValueInputMode(pathValueMap).
				GetSignal(); err != nil {
				return err
			} else {
				action := linebot.GetPostBackAction(string(method), js)
				departmentButton := linebot.GetButtonComponent(
					0,
					action,
					&domain.NormalButtonOption,
				)
				inputButtons = append(inputButtons, departmentButton)
			}
		}

		message := linebot.GetFlexMessage(
			"RichMenu Methods",
			linebot.GetFlexMessageBubbleContent(
				linebot.GetFlexMessageBoxComponent(
					linebotDomain.VERTICAL_MESSAGE_LAYOUT,
					&model.FlexMessageBoxComponentOption{
						JustifyContent: linebotDomain.SPACE_EVENLY_JUSTIFY_CONTENT,
					},
					inputButtons...,
				),
				nil,
			),
		)
		messages = append(messages, message)
	}

	if err := b.context.Reply(messages); err != nil {
		return err
	}

	return nil
}

func (b *richMenu) createRichMenu(name string, width, height, row, col int, actions ...interface{}) *linebotReqs.CreateRichMenu {
	result := &linebotReqs.CreateRichMenu{
		Size: linebotReqs.CreateRichMenuSize{
			Width:  width,
			Height: height,
		},
		Selected:    false,
		Name:        name,
		ChatBarText: "←開鍵盤  ↓開選單",
		Areas:       make([]*linebotReqs.CreateRichMenuAreas, 0),
	}

	widthUnit := float64(width) / float64(col)
	heightUnit := float64(height) / float64(row)
	for index, action := range actions {
		c := float64(index % col)
		r := math.Floor(float64(index) / float64(col))
		x := widthUnit * c
		y := heightUnit * r
		result.Areas = append(result.Areas, &linebotReqs.CreateRichMenuAreas{
			Action: action,
			Bounds: linebotReqs.CreateRichMenuAreasBounds{
				X: int(x),
				Y: int(y),
				CreateRichMenuSize: linebotReqs.CreateRichMenuSize{
					Width:  int(widthUnit),
					Height: int(heightUnit),
				},
			},
		})
	}

	return result
}

func (b *richMenu) createRoleRichMenu(role domain.ClubRole) *linebotReqs.CreateRichMenu {
	switch role {
	case domain.GUEST_CLUB_ROLE:
		return b.createRichMenu(
			"guest",
			2498, 1147,
			2, 3,
			linebot.GetMessageAction("無功能"),
			linebot.GetMessageAction("無功能"),
			linebot.GetMessageAction("查看活動"),
		)
	case domain.CADRE_CLUB_ROLE:
		return b.createRichMenu(
			"cadre",
			2498, 1721,
			3, 3,
			linebot.GetMessageAction("社長好強"),
			linebot.GetMessageAction("新增活動"),
			linebot.GetMessageAction("查看活動"),
		)
	case domain.ADMIN_CLUB_ROLE:
		return b.createRichMenu(
			"admin",
			2498, 1721,
			3, 3,
			linebot.GetMessageAction("richmenu"),
			linebot.GetMessageAction("新增活動"),
			linebot.GetMessageAction("查看活動"),
		)
	default:
		return nil
	}
}

func (b *richMenu) getRoleRichMenuImg(role domain.ClubRole) []byte {
	switch role {
	case domain.ADMIN_CLUB_ROLE:
		return adminRichMenuImg
	case domain.CADRE_CLUB_ROLE:
		return cadreRichMenuImg
	default:
		return guestRichMenuImg
	}
}
