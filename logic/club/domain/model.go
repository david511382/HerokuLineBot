package domain

type CmdBase struct {
	Cmd                     TextCmd `json:"cmd,omitempty"`
	RequireRawParamAttr     string  `json:"require_raw_param_attr"`
	IsInputImmediately      bool    `json:"is_input_immediately"`
	RequireRawParamAttrText string  `json:"require_raw_param_attr_text"`
	IsSingleParamMode       bool    `json:"-"`
	IsCancel                bool    `json:"is_cancel"`
	IsComfirm               bool    `json:"is_comfirm,omitempty"`
}

type TimePostbackParams struct {
	Date     string `json:"date"`
	DateTime string `json:"date_time"`
	Time     string `json:"time"`
}
