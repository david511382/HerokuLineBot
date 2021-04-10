package common

type BaseList struct {
	Base
}

func (k *BaseList) LRange(start, stop int64) ([]string, error) {
	dp := k.Read.LRange(k.Key, start, stop)

	if err := dp.Err(); err != nil {
		return nil, err
	}

	result, err := dp.Result()
	return result, err
}

func (k *BaseList) RPush(value interface{}) (int64, error) {
	dp := k.Write.RPush(k.Key, value)

	if err := dp.Err(); err != nil {
		return 0, err
	}

	result, err := dp.Result()
	return result, err
}
