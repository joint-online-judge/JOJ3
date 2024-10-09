package stage

import (
	"fmt"

	"github.com/mcuadros/go-defaults"
	"github.com/mitchellh/mapstructure"
)

func DecodeConf[T any](confAny any) (*T, error) {
	conf := new(T)
	defaults.SetDefaults(conf)
	err := mapstructure.Decode(confAny, conf)
	if err != nil {
		return nil, fmt.Errorf("failed to decode conf: %w", err)
	}
	return conf, nil
}
