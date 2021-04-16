package club

import (
	"heroku-line-bot/logic/club/domain"

	"github.com/tidwall/sjson"
)

func (b *CmdHandler) setPathValue(key string, value interface{}) {
	if b.pathValueMap == nil {
		b.pathValueMap = make(map[string]interface{})
	}
	b.pathValueMap[key] = value
}

func (b *CmdHandler) duplicate() *CmdHandler {
	nb := *b
	cb := *b.CmdBase
	nb.CmdBase = &cb
	nb.pathValueMap = make(map[string]interface{})
	for k, v := range b.pathValueMap {
		nb.pathValueMap[k] = v
	}
	return &nb
}

func (b *CmdHandler) GetKeyValueInputMode(pathValueMap map[string]interface{}) domain.ICmdHandlerSignal {
	nb := b.duplicate()
	for path, value := range pathValueMap {
		nb.setPathValue(path, value)
	}
	return nb
}

func (b *CmdHandler) GetCmdInputMode(cmdP *domain.TextCmd) domain.ICmdHandlerSignal {
	nb := b.duplicate()
	cmd := nb.Cmd
	if cmdP != nil {
		cmd = *cmdP
	}
	nb.setPathValue(string(domain.CMD_ATTR), cmd)
	return nb
}

func (b *CmdHandler) GetDateTimeCmdInputMode(timeCmd domain.DateTimeCmd, attr string) domain.ICmdHandlerSignal {
	nb := b.duplicate()
	nb.setPathValue(string(domain.DATE_TIME_CMD_ATTR), timeCmd)
	return nb.GetRequireInputMode(attr, "", true)
}

func (b *CmdHandler) GetCancelMode() domain.ICmdHandlerSignal {
	nb := b.duplicate()
	nb.setPathValue("is_cancel", true)
	return nb
}

func (b *CmdHandler) GetComfirmMode() domain.ICmdHandlerSignal {
	nb := b.duplicate()
	nb.setPathValue("is_comfirm", true)
	return nb
}

func (b *CmdHandler) GetCancelInputMode() domain.ICmdHandlerSignal {
	return b.GetRequireInputMode("", "", false)
}

func (b *CmdHandler) GetRequireInputMode(attr, attrText string, isInputImmediately bool) domain.ICmdHandlerSignal {
	nb := b.duplicate()
	nb.setPathValue("require_raw_param_attr", attr)
	nb.setPathValue("require_raw_param_attr_text", attrText)
	nb.setPathValue("is_input_immediately", isInputImmediately)
	return nb
}

func (b *CmdHandler) GetSignal() (string, error) {
	js := "{}"
	for path, value := range b.pathValueMap {
		var err error
		if js, err = sjson.Set(js, path, value); err != nil {
			return "", err
		}
	}

	return js, nil
}
