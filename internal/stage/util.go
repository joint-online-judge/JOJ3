package stage

import "github.com/mitchellh/mapstructure"

func DecodeConf[T any](confAny any) (*T, error) {
	var conf T
	err := mapstructure.Decode(confAny, &conf)
	if err != nil {
		return nil, err
	}
	return &conf, nil
}
