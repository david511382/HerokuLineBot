package club

import (
	"heroku-line-bot/src/logic/club/domain"
	errUtil "heroku-line-bot/src/pkg/util/error"

	"github.com/tidwall/sjson"
)

type Signal map[string]interface{}

func NewSignal() Signal {
	return make(Signal)
}

func (s Signal) setPathValue(key string, value interface{}) {
	if s == nil {
		return
	}
	s[key] = value
}

func (s Signal) GetKeyValueInputMode(pathValueMap map[string]interface{}) Signal {
	for path, value := range pathValueMap {
		s.setPathValue(path, value)
	}
	return s
}

func (s Signal) GetRequireInputMode(requireRawAttr string) Signal {
	s.setPathValue("require_raw_attr", requireRawAttr)
	return s
}

func (s Signal) GetCancelMode() Signal {
	s.setPathValue("is_cancel", true)
	return s
}

func (s Signal) GetConfirmMode() Signal {
	s.setPathValue("is_comfirm", true)
	return s
}

func (s Signal) GetRunOnceMode() Signal {
	s.setPathValue("is_not_cache", true)
	return s
}

func (s Signal) GetCancelInputMode() Signal {
	s.setPathValue("require_raw_attr", "")
	return s
}

func (s Signal) GetCmdInputMode(cmdP *domain.TextCmd) Signal {
	cmd := *cmdP
	s.setPathValue(string(domain.ATTR_CMD), cmd)
	return s
}

func (s Signal) GetBasePath(path string) Signal {
	s.setPathValue(string(domain.ATTR_CMD_BASE_PATH), path)
	return s
}

func (s Signal) GetSignal() (string, errUtil.IError) {
	js := "{}"
	for path, value := range s {
		var err error
		if js, err = sjson.Set(js, path, value); err != nil {
			return "", errUtil.NewError(err)
		}
	}

	return js, nil
}
