// Package docs GENERATED BY THE COMMAND ABOVE; DO NOT EDIT
// This file was generated by swaggo/swag
package docs

import (
	"bytes"
	"encoding/json"
	"strings"
	"text/template"

	"github.com/swaggo/swag"
)

var doc = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "contact": {},
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/badminton/activitys": {
            "get": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "尚未截止的活動",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Badminton"
                ],
                "summary": "活動列表",
                "parameters": [
                    {
                        "type": "string",
                        "description": "起始日期",
                        "name": "from_date",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "結束日期",
                        "name": "to_date",
                        "in": "query"
                    },
                    {
                        "type": "array",
                        "items": {
                            "type": "integer"
                        },
                        "description": "場館IDs",
                        "name": "place_ids",
                        "in": "query"
                    },
                    {
                        "type": "array",
                        "items": {
                            "type": "integer"
                        },
                        "description": "球隊IDs",
                        "name": "team_ids",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "description": "分頁每頁資料量",
                        "name": "page_size",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "integer",
                        "description": "分頁第幾頁，1開始",
                        "name": "page_index",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "資料",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/resp.Base"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "data": {
                                            "$ref": "#/definitions/resp.GetActivitys"
                                        }
                                    }
                                }
                            ]
                        }
                    }
                }
            }
        },
        "/badminton/rental-courts": {
            "get": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "租場狀況",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Badminton"
                ],
                "summary": "租場狀況",
                "parameters": [
                    {
                        "type": "string",
                        "default": "2013-08-02T00:00:00+08:00",
                        "description": "起始日期",
                        "name": "from_date",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "string",
                        "default": "2013-08-02T00:00:00+08:00",
                        "description": "結束日期",
                        "name": "to_date",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "integer",
                        "description": "球隊id",
                        "name": "team_id",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "資料",
                        "schema": {
                            "$ref": "#/definitions/resp.GetRentalCourts"
                        }
                    }
                }
            },
            "post": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "新增租場",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Badminton"
                ],
                "summary": "新增租場",
                "parameters": [
                    {
                        "description": "參數",
                        "name": "param",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/reqs.AddRentalCourt"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "資料",
                        "schema": {
                            "$ref": "#/definitions/resp.Base"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "reqs.AddRentalCourt": {
            "type": "object",
            "required": [
                "court_count",
                "court_from_time",
                "court_to_time",
                "from_date",
                "place_id",
                "price_per_hour",
                "team_id",
                "to_date"
            ],
            "properties": {
                "balance_date": {
                    "type": "string"
                },
                "balance_money": {
                    "type": "integer"
                },
                "court_count": {
                    "type": "integer"
                },
                "court_from_time": {
                    "type": "string"
                },
                "court_to_time": {
                    "type": "string"
                },
                "desposit_date": {
                    "type": "string"
                },
                "desposit_money": {
                    "type": "integer"
                },
                "every_weekday": {
                    "type": "integer"
                },
                "exclude_dates": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                },
                "from_date": {
                    "type": "string",
                    "default": "2013-08-02T00:00:00+08:00"
                },
                "place_id": {
                    "type": "integer"
                },
                "price_per_hour": {
                    "type": "integer"
                },
                "team_id": {
                    "type": "integer"
                },
                "to_date": {
                    "type": "string",
                    "default": "2013-08-02T00:00:00+08:00"
                }
            }
        },
        "resp.Base": {
            "type": "object",
            "properties": {
                "data": {},
                "message": {
                    "type": "string"
                }
            }
        },
        "resp.GetActivitys": {
            "type": "object",
            "properties": {
                "activitys": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/resp.GetActivitysActivity"
                    }
                },
                "data_count": {
                    "type": "integer"
                }
            }
        },
        "resp.GetActivitysActivity": {
            "type": "object",
            "properties": {
                "activity_id": {
                    "type": "integer"
                },
                "courts": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/resp.GetActivitysCourt"
                    }
                },
                "date": {
                    "type": "string"
                },
                "description": {
                    "type": "string"
                },
                "is_show_members": {
                    "type": "boolean"
                },
                "members": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/resp.GetActivitysMember"
                    }
                },
                "people_limit": {
                    "type": "integer"
                },
                "place_id": {
                    "type": "integer"
                },
                "place_name": {
                    "type": "string"
                },
                "price": {
                    "type": "integer"
                },
                "team_id": {
                    "type": "integer"
                },
                "team_name": {
                    "type": "string"
                }
            }
        },
        "resp.GetActivitysCourt": {
            "type": "object",
            "properties": {
                "count": {
                    "type": "integer"
                },
                "from_time": {
                    "type": "string"
                },
                "to_time": {
                    "type": "string"
                }
            }
        },
        "resp.GetActivitysMember": {
            "type": "object",
            "properties": {
                "id": {
                    "type": "integer"
                },
                "name": {
                    "type": "string"
                }
            }
        },
        "resp.GetRentalCourts": {
            "type": "object",
            "properties": {
                "not_pay_day_courts": {
                    "$ref": "#/definitions/resp.GetRentalCourtsPayInfo"
                },
                "not_refund_day_courts": {
                    "$ref": "#/definitions/resp.GetRentalCourtsPayInfo"
                },
                "total_day_courts": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/resp.GetRentalCourtsDayCourts"
                    }
                }
            }
        },
        "resp.GetRentalCourtsCourtInfo": {
            "type": "object",
            "properties": {
                "cost": {
                    "type": "number"
                },
                "count": {
                    "type": "integer"
                },
                "from_time": {
                    "type": "string"
                },
                "place": {
                    "type": "string"
                },
                "to_time": {
                    "type": "string"
                }
            }
        },
        "resp.GetRentalCourtsDayCourts": {
            "type": "object",
            "properties": {
                "courts": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/resp.GetRentalCourtsDayCourtsInfo"
                    }
                },
                "date": {
                    "type": "string"
                },
                "is_multiple_place": {
                    "type": "boolean"
                }
            }
        },
        "resp.GetRentalCourtsDayCourtsInfo": {
            "type": "object",
            "properties": {
                "cost": {
                    "type": "number"
                },
                "count": {
                    "type": "integer"
                },
                "from_time": {
                    "type": "string"
                },
                "place": {
                    "type": "string"
                },
                "reason_message": {
                    "type": "string"
                },
                "refund_time": {
                    "type": "string"
                },
                "status": {
                    "type": "integer"
                },
                "to_time": {
                    "type": "string"
                }
            }
        },
        "resp.GetRentalCourtsPayInfo": {
            "type": "object",
            "properties": {
                "cost": {
                    "type": "number"
                },
                "courts": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/resp.GetRentalCourtsPayInfoDay"
                    }
                }
            }
        },
        "resp.GetRentalCourtsPayInfoDay": {
            "type": "object",
            "properties": {
                "cost": {
                    "type": "number"
                },
                "courts": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/resp.GetRentalCourtsCourtInfo"
                    }
                },
                "date": {
                    "type": "string"
                }
            }
        }
    },
    "securityDefinitions": {
        "ApiKeyAuth": {
            "type": "apiKey",
            "name": "Authorization",
            "in": "header"
        }
    }
}`

type swaggerInfo struct {
	Version     string
	Host        string
	BasePath    string
	Schemes     []string
	Title       string
	Description string
}

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = swaggerInfo{
	Version:     "1.0",
	Host:        "",
	BasePath:    "/api/",
	Schemes:     []string{},
	Title:       "Heroku-Line-Bot",
	Description: "Line-Bot",
}

type s struct{}

func (s *s) ReadDoc() string {
	sInfo := SwaggerInfo
	sInfo.Description = strings.Replace(sInfo.Description, "\n", "\\n", -1)

	t, err := template.New("swagger_info").Funcs(template.FuncMap{
		"marshal": func(v interface{}) string {
			a, _ := json.Marshal(v)
			return string(a)
		},
		"escape": func(v interface{}) string {
			// escape tabs
			str := strings.Replace(v.(string), "\t", "\\t", -1)
			// replace " with \", and if that results in \\", replace that with \\\"
			str = strings.Replace(str, "\"", "\\\"", -1)
			return strings.Replace(str, "\\\\\"", "\\\\\\\"", -1)
		},
	}).Parse(doc)
	if err != nil {
		return doc
	}

	var tpl bytes.Buffer
	if err := t.Execute(&tpl, sInfo); err != nil {
		return doc
	}

	return tpl.String()
}

func init() {
	swag.Register("swagger", &s{})
}
