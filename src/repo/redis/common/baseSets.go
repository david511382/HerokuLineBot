package common

type BaseSets struct {
	Base
}

func (k *BaseSets) SMembers() ([]string, error) {
	conn, err := k.connection.GetMaster()
	if err != nil {
		return nil, err
	}

	dp := conn.SMembers(k.Key)
	if err := dp.Err(); err != nil {
		return nil, err
	}

	result, err := dp.Result()
	return result, err
}

func (k *BaseSets) SAdd(values ...string) (int64, error) {
	conn, err := k.connection.GetMaster()
	if err != nil {
		return 0, err
	}

	dp := conn.SAdd(k.Key, values)
	if err := dp.Err(); err != nil {
		return 0, err
	}

	result, err := dp.Result()
	return result, err
}
