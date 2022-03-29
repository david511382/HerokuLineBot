package club

import (
	"heroku-line-bot/src/logger"
	accountLineuserLogic "heroku-line-bot/src/logic/account/lineuser"
	"heroku-line-bot/src/logic/club/domain"
	dbModel "heroku-line-bot/src/model/database"
	"heroku-line-bot/src/pkg/service/linebot"
	linebotDomain "heroku-line-bot/src/pkg/service/linebot/domain"
	"heroku-line-bot/src/pkg/service/linebot/domain/model"
	linebotReqs "heroku-line-bot/src/pkg/service/linebot/domain/model/reqs"
	lineResp "heroku-line-bot/src/pkg/service/linebot/domain/model/resp"
	"heroku-line-bot/src/pkg/util"
	errUtil "heroku-line-bot/src/pkg/util/error"
	"heroku-line-bot/src/repo/database"
	"heroku-line-bot/src/repo/database/database/clubdb/member"
	"sort"
	"strconv"

	"github.com/rs/zerolog"
)

type richMenu struct {
	context domain.ICmdHandlerContext `json:"-"`
	Role    *domain.ClubRole          `json:"role"`
}

func (b *richMenu) Init(context domain.ICmdHandlerContext) (resultErrInfo errUtil.IError) {
	*b = richMenu{
		context: context,
	}

	return
}

func (b *richMenu) GetRequireAttr() (requireAttr string, warnMessage interface{}, resultErrInfo errUtil.IError) {
	return
}

func (b *richMenu) GetRequireAttrInfo(rawAttr string) (attrNameText string, valueText string, isNotRequireChecking bool) {
	switch rawAttr {
	}
	return
}

func (b *richMenu) GetInputTemplate(attr string) (messages interface{}) {
	switch attr {
	case "role":
		buttons := []interface{}{}
		roles := []domain.ClubRole{
			domain.ADMIN_CLUB_ROLE,
			domain.CADRE_CLUB_ROLE,
			domain.GUEST_CLUB_ROLE,
		}
		for _, role := range roles {
			pathValueMap := map[string]interface{}{
				"ICmdLogic.role": role,
			}
			if js, errInfo := NewSignal().
				GetRunOnceMode().
				//	GetCancelInputMode().
				// if js, errInfo := b.context.
				// 	GetCmdInputMode(nil).
				GetKeyValueInputMode(pathValueMap).
				GetSignal(); errInfo != nil {
				logger.Log("line bot club", errInfo)
				return
			} else {
				action := linebot.GetPostBackAction(strconv.Itoa(int(role)), js)
				button := linebot.GetButtonComponent(
					action,
					&domain.NormalButtonOption,
				)
				buttons = append(buttons, button)
			}
		}
		messages = linebot.GetFlexMessage(
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
	}
	return
}

func (b *richMenu) LoadRequireInputTextParam(attr, text string) (resultErrInfo errUtil.IError) {
	switch attr {
	default:
	}

	return nil
}

func (b *richMenu) Do(text string) (resultErrInfo errUtil.IError) {
	if u, err := accountLineuserLogic.Get(b.context.GetUserID()); err != nil {
		resultErrInfo = errUtil.NewError(err)
		return
	} else {
		if u.Role != domain.ADMIN_CLUB_ROLE {
			resultErrInfo = errUtil.NewError(domain.NO_AUTH_ERROR)
			return
		}
	}

	messages := []interface{}{}
	if b.Role == nil {
		buttons := []interface{}{}
		roles := []domain.ClubRole{
			domain.ADMIN_CLUB_ROLE,
			domain.CADRE_CLUB_ROLE,
			domain.GUEST_CLUB_ROLE,
		}
		for _, role := range roles {
			pathValueMap := map[string]interface{}{
				"ICmdLogic.role": role,
			}
			if js, errInfo := NewSignal().
				//GetCmdInputMode(nil).
				GetKeyValueInputMode(pathValueMap).
				GetSignal(); errInfo != nil {
				logger.Log("line bot club", errInfo)
				return nil
			} else {
				action := linebot.GetPostBackAction(strconv.Itoa(int(role)), js)
				button := linebot.GetButtonComponent(
					action,
					&domain.NormalButtonOption,
				)
				buttons = append(buttons, button)
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
					buttons...,
				),
				nil,
			),
		)
		messages = append(messages, message)
	} else {
		role := *b.Role
		richMenuName := b.getRichMenuRoleName(role)

		originRichMenus := make([]*lineResp.ListRichMenuRichMenu, 0)
		{
			resp, err := b.context.GetBot().ListRichMenu()
			if err != nil {
				errInfo := errUtil.NewError(err)
				resultErrInfo = errUtil.Append(resultErrInfo, errInfo)
				return
			}

			for _, v := range resp.RichMenus {
				if v.Name != richMenuName {
					continue
				}
				originRichMenus = append(originRichMenus, v)
			}
		}

		createRichMenuArg := b.createRoleRichMenu(role)
		createRichMenuResp, err := b.context.GetBot().CreateRichMenu(createRichMenuArg)
		if err != nil {
			errInfo := errUtil.NewError(err)
			resultErrInfo = errUtil.Append(resultErrInfo, errInfo)
			return
		}

		imgBs := b.getRoleRichMenuImg(role)
		if _, err := b.context.GetBot().UploadRichMenuImage(createRichMenuResp.RichMenuID, imgBs); err != nil {
			resultErrInfo = errUtil.NewError(err)
			return
		}

		var updateMemberMessage interface{}
		{
			updateMemberNames := make([]string, 0)
			switch role {
			case domain.ADMIN_CLUB_ROLE, domain.CADRE_CLUB_ROLE:
				lineIDs := []string{}
				if dbDatas, err := database.Club.Member.Select(
					dbModel.ReqsClubMember{
						Role: util.GetInt16P(int16(role)),
					},
					member.COLUMN_Name,
					member.COLUMN_LineID,
				); err != nil {
					resultErrInfo = errUtil.NewError(err)
					return
				} else if len(dbDatas) > 0 {
					for _, v := range dbDatas {
						if v.LineID == nil {
							continue
						}
						lineIDs = append(lineIDs, *v.LineID)
						updateMemberNames = append(updateMemberNames, v.Name)
					}
				}

				for _, lineID := range lineIDs {
					if _, err := b.context.GetBot().SetRichMenuTo(createRichMenuResp.RichMenuID, lineID); err != nil {
						resultErrInfo = errUtil.NewError(err)
						return
					}
				}
			default:
				if _, err := b.context.GetBot().SetDefaultRichMenu(createRichMenuResp.RichMenuID); err != nil {
					resultErrInfo = errUtil.NewError(err)
					return
				}

				updateMemberNames = append(updateMemberNames, "default")
			}

			if len(updateMemberNames) > 0 {
				contents := make([]interface{}, 0)
				contents = append(contents, linebot.GetFlexMessageTextComponent(
					"Set To", nil,
				))
				for _, updateMemberName := range updateMemberNames {
					contents = append(contents, linebot.GetFlexMessageTextComponent(
						updateMemberName, nil,
					))
				}

				updateMemberMessage = linebot.GetFlexMessage(
					"更新的 RichMenu",
					linebot.GetFlexMessageBubbleContent(
						linebot.GetFlexMessageBoxComponent(
							linebotDomain.VERTICAL_MESSAGE_LAYOUT,
							&model.FlexMessageBoxComponentOption{
								JustifyContent: linebotDomain.SPACE_EVENLY_JUSTIFY_CONTENT,
							},
							contents...,
						),
						nil,
					),
				)
			}
		}

		var deletedMessage interface{}
		{
			richMenuIDs := make([]string, 0)
			for _, v := range originRichMenus {
				richMenuID := v.RichMenuID
				if richMenuID == createRichMenuResp.RichMenuID {
					continue
				}
				_, err := b.context.GetBot().DeleteRichMenu(richMenuID)
				if err != nil {
					errInfo := errUtil.NewError(err, zerolog.WarnLevel)
					resultErrInfo = errUtil.Append(resultErrInfo, errInfo)
					continue
				}

				richMenuIDs = append(richMenuIDs, richMenuID)
			}

			deletedContents := make([]interface{}, 0)
			sort.Slice(richMenuIDs, func(i, j int) bool {
				return richMenuIDs[i] < richMenuIDs[j]
			})
			for _, richMenuID := range richMenuIDs {
				deletedContents = append(deletedContents, linebot.GetFlexMessageTextComponent(
					richMenuID, nil,
				))
			}

			if len(deletedContents) > 0 {
				contents := make([]interface{}, 0)
				contents = append(contents, linebot.GetFlexMessageTextComponent(
					"刪除的 RichMenu", nil,
				))
				contents = append(contents, deletedContents...)

				deletedMessage = linebot.GetFlexMessage(
					"刪除的 RichMenu",
					linebot.GetFlexMessageBubbleContent(
						linebot.GetFlexMessageBoxComponent(
							linebotDomain.VERTICAL_MESSAGE_LAYOUT,
							&model.FlexMessageBoxComponentOption{
								JustifyContent: linebotDomain.SPACE_EVENLY_JUSTIFY_CONTENT,
							},
							contents...,
						),
						nil,
					),
				)
			}
		}

		messages = append(messages, linebot.GetTextMessage(
			"更新 RichMenu 完成",
		))
		if deletedMessage != nil {
			messages = append(messages, deletedMessage)
		}
		messages = append(messages, linebot.GetTextMessage(
			"建立的 RichMenuID:",
		))
		messages = append(messages, linebot.GetTextMessage(
			createRichMenuResp.RichMenuID,
		))
		if updateMemberMessage != nil {
			messages = append(messages, updateMemberMessage)
		}

		if err := b.context.DeleteParam(); err != nil {
			logger.Log("line bot club", errUtil.NewError(err))
		}
	}
	if err := b.context.Reply(messages); err != nil {
		resultErrInfo = errUtil.NewError(err)
		return
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

	widthUnit := util.NewFloat(float64(width)).Div(util.NewFloat(float64(col)))
	heightUnit := util.NewFloat(float64(height)).Div(util.NewFloat(float64(row)))
	for index, action := range actions {
		c := float64(index % col)
		r := util.NewFloat(float64(index)).
			DivFloat(float64(col)).
			Floor()
		x := widthUnit.MulFloat(c)
		y := heightUnit.Mul(r)
		result.Areas = append(result.Areas, &linebotReqs.CreateRichMenuAreas{
			Action: action,
			Bounds: linebotReqs.CreateRichMenuAreasBounds{
				X: int(x.ToInt()),
				Y: int(y.ToInt()),
				CreateRichMenuSize: linebotReqs.CreateRichMenuSize{
					Width:  int(widthUnit.ToInt()),
					Height: int(heightUnit.ToInt()),
				},
			},
		})
	}

	return result
}

func (b *richMenu) getRichMenuRoleName(role domain.ClubRole) string {
	switch role {
	case domain.GUEST_CLUB_ROLE:
		return "guest"
	case domain.CADRE_CLUB_ROLE:
		return "cadre"
	case domain.ADMIN_CLUB_ROLE:
		return "admin"
	default:
		return ""
	}
}

func (b *richMenu) createRoleRichMenu(role domain.ClubRole) *linebotReqs.CreateRichMenu {
	name := b.getRichMenuRoleName(role)
	switch role {
	case domain.GUEST_CLUB_ROLE:
		return b.createRichMenu(
			name,
			2498, 1147,
			2, 3,
			linebot.GetMessageAction("社長好強"),
			linebot.GetMessageAction("經理好棒"),
			linebot.GetMessageAction(string(domain.GET_ACTIVITIES_TEXT_CMD)),
		)
	case domain.CADRE_CLUB_ROLE:
		return b.createRichMenu(
			name,
			2498, 1721,
			3, 3,
			linebot.GetMessageAction("社長好強"),
			linebot.GetMessageAction(string(domain.NEW_ACTIVITY_TEXT_CMD)),
			linebot.GetMessageAction(string(domain.GET_ACTIVITIES_TEXT_CMD)),
		)
	case domain.ADMIN_CLUB_ROLE:
		return b.createRichMenu(
			name,
			2498, 1721,
			3, 3,
			linebot.GetMessageAction(string(domain.RICH_MENU_TEXT_CMD)),
			linebot.GetMessageAction(string(domain.NEW_ACTIVITY_TEXT_CMD)),
			linebot.GetMessageAction(string(domain.GET_ACTIVITIES_TEXT_CMD)),
			linebot.GetMessageAction(string(domain.GET_CONFIRM_REGISTER_TEXT_CMD)),
			linebot.GetMessageAction(string(domain.NEW_LOGISTIC_TEXT_CMD)),
			linebot.GetUriAction(liffUrl),
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
