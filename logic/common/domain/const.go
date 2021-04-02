package domain

type TimeType uint

const (
	YEAR_TYPE  TimeType = 1
	MONTH_TYPE TimeType = 2
	DAY_TYPE   TimeType = 3
	HOUR_TYPE  TimeType = 4
)

const IANA_ZONE = "Asia/Taipei"

const (
	FLOAT_EXPONENT                         = -8
	FRONT_END_RTP_SYMBOL_EXPONENT          = -4
	FRONT_END_FLOAT_EXPONENT               = -2
	FRONT_END_FISH_FLOAT_EXPONENT          = -3
	FRONT_END_FISH_EXPECTED_FLOAT_EXPONENT = -5
)

const (
	DATE_TIME_FORMAT         = "2006-01-02 15:04:05"
	DATE_FORMAT              = "2006-01-02"
	TIME_FORMAT              = "15:04:05"
	TIME_HOUR_MIN_FORMAT     = "15:04"
	DATE_TIME_RFC3339_FORMAT = "2006-01-02T15:04:05Z07:00"
)
