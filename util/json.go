package util

import (
	"encoding/json"

	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
)

func IsJSON(s string) bool {
	var js map[string]interface{}
	return json.Unmarshal([]byte(s), &js) == nil
}

type Json struct {
	json string
}

func NewJson(js string) *Json {
	return &Json{
		json: js,
	}
}

func (e *Json) Raw() string {
	return e.json
}

func (e *Json) GetAttrJson(attr string) *Json {
	result := &Json{
		json: gjson.Get(e.json, attr).Raw,
	}

	return result
}

func (e *Json) SetAttrJson(attr string, value interface{}) error {
	if newJs, err := sjson.Set(e.json, attr, value); err != nil {
		return err
	} else {
		e.json = newJs
	}
	return nil
}

func (e *Json) GetAttrValue(attr string) gjson.Result {
	return gjson.Get(e.json, attr)
}

func (e *Json) Parse(model interface{}) error {
	if err := json.Unmarshal([]byte(e.json), model); err != nil {
		return err
	}
	return nil
}
