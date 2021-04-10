package domain

import (
	"errors"
)

var (
	NOT_CHANGE  error = errors.New("redis not change")
	NOT_SUCCESS error = errors.New("redis not success")
	NO_DATA     error = errors.New("redis no data")
	NOT_EXIST   error = errors.New("redis: nil")
)
