package club

import (
	"encoding/json"
	"fmt"
	"heroku-line-bot/logic/club/domain"
	clublinebotDomain "heroku-line-bot/logic/clublinebot/domain"
	"heroku-line-bot/service/linebot"
	linebotDomain "heroku-line-bot/service/linebot/domain"
	"heroku-line-bot/service/linebot/domain/model"
)

type CmdHandler struct {
	*domain.CmdBase
	*domain.TimePostbackParams
	clublinebotDomain.IContext `json:"-"`
	domain.ICmdLogic
}

func (b *CmdHandler) ReadParam(jsonBytes []byte) error {
	return json.Unmarshal(jsonBytes, b)
}

func (b *CmdHandler) IsInputMode() bool {
	return b.RequireRawParamAttr != ""
}

func (b *CmdHandler) GetSingleParamText() string {
	return b.ICmdLogic.GetSingleParam(b.RequireRawParamAttr)
}

func (b *CmdHandler) LoadSingleParamValue(valueText string) (resultValue interface{}, resultErr error) {
	return b.ICmdLogic.LoadSingleParam(b.RequireRawParamAttr, valueText)
}

func (b *CmdHandler) duplicate() *CmdHandler {
	nb := *b
	cb := *b.CmdBase
	nb.CmdBase = &cb
	return &nb
}

func (b *CmdHandler) GetCancelSignl() string {
	nb := b.duplicate()
	nb.IsCancel = true
	js, err := nb.GetRequireInputCmdText(nil, "", "", true)
	if err != nil {
		return ""
	}
	return js
}

func (b *CmdHandler) GetComfirmSignl() string {
	nb := b.duplicate()
	nb.IsComfirm = true
	js, err := nb.GetRequireInputCmdText(nil, "", "", true)
	if err != nil {
		return ""
	}
	return js
}

func (b *CmdHandler) GetRequireInputCmdText(cmd *domain.TextCmd, attr, attrText string, isInputImmediately bool) (string, error) {
	b.RequireRawParamAttr = attr
	b.RequireRawParamAttrText = attrText
	b.IsInputImmediately = isInputImmediately

	nb := b.duplicate()
	if cmd != nil {
		nb.Cmd = *cmd
	} else {
		nb.Cmd = ""
	}
	if jsBytes, err := json.Marshal(nb.CmdBase); err != nil {
		return "", err
	} else {
		js := string(jsBytes)
		return js, nil
	}
}

func (b *CmdHandler) CacheParams() error {
	if jsBytes, err := json.Marshal(b); err != nil {
		return err
	} else {
		js := string(jsBytes)
		if err := b.SaveParam(js); err != nil {
			return err
		}
	}
	return nil
}

func (b *CmdHandler) IsComfirmed() bool {
	return b.IsComfirm
}

func (b *CmdHandler) SetSingleParamMode() {
	b.IsSingleParamMode = true
}

func (b *CmdHandler) Do(text string) error {
	if b.IsCancel {
		if err := b.DeleteParam(); err != nil {
			return err
		}

		replyMessges := []interface{}{
			linebot.GetTextMessage("取消"),
		}
		if err := b.IContext.Reply(replyMessges); err != nil {
			return err
		}

		return nil
	}

	if b.IsInputMode() {
		const altText = "請確認或輸入"

		if b.IsSingleParamMode {
			_, err := b.LoadSingleParamValue(text)
			if err != nil {
				msg := fmt.Sprintf("參數格式錯誤:%s", err.Error())
				replyMessges := []interface{}{
					linebot.GetTextMessage(msg),
				}
				if err := b.Reply(replyMessges); err != nil {
					return err
				}
				return nil
			}
		}

		valueText := b.GetSingleParamText()
		var text = fmt.Sprintf("%s %s ,確認或請輸入數值", b.RequireRawParamAttrText, valueText)

		if b.IsInputImmediately {
			b.RequireRawParamAttr = ""
			b.RequireRawParamAttrText = ""
		}

		if err := b.CacheParams(); err != nil {
			return err
		}

		if !b.IsInputImmediately {
			cancelRequireInputJs, err := b.GetRequireInputCmdText(nil, "", "", false)
			if err != nil {
				return err
			}

			replyMessges := []interface{}{
				linebot.GetFlexMessage(
					altText,
					linebot.GetFlexMessageBubbleContent(
						linebot.GetFlexMessageBoxComponent(
							linebotDomain.VERTICAL_MESSAGE_LAYOUT,
							&model.FlexMessageBoxComponentOption{
								JustifyContent: linebotDomain.SPACE_EVENLY_JUSTIFY_CONTENT,
							},
							linebot.GetButtonComponent(
								0,
								linebot.GetPostBackAction(
									"確認",
									cancelRequireInputJs,
								),
							),
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
				),
			}
			if err := b.IContext.Reply(replyMessges); err != nil {
				return err
			}

			return nil
		}
	}

	return b.ICmdLogic.Do(text)
}
