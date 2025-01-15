package common

import (
	"github.com/mitchellh/mapstructure"
)

func Decode(input interface{}, output interface{}) error {
	decoder, err := mapstructure.NewDecoder(
		&mapstructure.DecoderConfig{
			ZeroFields: true,
			TagName:    "json",
			Result:     output,
		},
	)
	if err != nil {
		return err
	}

	return decoder.Decode(input)
}
