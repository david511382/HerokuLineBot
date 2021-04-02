package domain

type IContext interface {
	SaveParam(json string) error
	DeleteParam() error
	GetParam() (json string)
	Reply(replyMessges []interface{}) error
	GetUserID() string
}
