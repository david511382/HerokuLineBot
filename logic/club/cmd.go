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
	pathValueMap map[string]interface{}
}

func (b *CmdHandler) ReadParam(jsonBytes []byte) error {
	return json.Unmarshal(jsonBytes, b)
}

func (b *CmdHandler) IsInputMode() bool {
	return b.RequireRawParamAttr != ""
}

func (b *CmdHandler) SetRequireInputMode(attr, attrText string, isInputImmediately bool) {
	b.RequireRawParamAttr = attr
	b.RequireRawParamAttrText = attrText
	b.IsInputImmediately = isInputImmediately
}

func (b *CmdHandler) LoadSingleParamValue(valueText string) error {
	return b.ICmdLogic.LoadSingleParam(b.RequireRawParamAttr, valueText)
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

func (b *CmdHandler) GetInputTemplate(requireRawParamAttr string) interface{} {
	const altText = "請確認或輸入"
	valueText := b.ICmdLogic.GetSingleParam(requireRawParamAttr)
	var text = fmt.Sprintf("%s %s ,確認或請輸入數值", b.RequireRawParamAttrText, valueText)

	cancelRequireInputJs, err := b.GetCancelInputMode().GetSignal()
	if err != nil {
		return err
	}

	return linebot.GetFlexMessage(
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
					&domain.NormalButtonOption,
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
	)
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
		if b.IsSingleParamMode {
			if err := b.LoadSingleParamValue(text); err != nil {
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

		requireRawParamAttr := b.RequireRawParamAttr
		if b.IsInputImmediately {
			b.RequireRawParamAttr = ""
		}

		if err := b.CacheParams(); err != nil {
			return err
		}

		if !b.IsInputImmediately {
			replyMessge := b.ICmdLogic.GetInputTemplate(requireRawParamAttr)
			if replyMessge == nil {
				replyMessge = b.GetInputTemplate(requireRawParamAttr)
			}
			replyMessges := []interface{}{
				replyMessge,
			}
			if err := b.IContext.Reply(replyMessges); err != nil {
				return err
			}

			return nil
		}
	}

	return b.ICmdLogic.Do(text)
}
