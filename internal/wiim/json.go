package wiim

import (
	"errors"
	"fmt"
)

type JsonBoolean bool

func (bit *JsonBoolean) UnmarshalJSON(data []byte) error {
	str := string(data)
	if str == "1" || str == "true" || str == "\"1\"" {
		*bit = true
	} else if str == "0" || str == "false" || str == "\"0\"" {
		*bit = false
	} else {
		return errors.New(fmt.Sprintf("Boolean unmarshal error: invalid input %s", str))
	}
	return nil
}
