package stage

import "github.com/mitchellh/mapstructure"

func DecodeConfig[T any](configAny any) (*T, error) {
	var config T
	err := mapstructure.Decode(configAny, &config)
	if err != nil {
		return nil, err
	}
	return &config, nil
}
