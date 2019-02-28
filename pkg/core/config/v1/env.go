package v1

import "strings"

func ParseEnv(str string) (*Env, error) {
	ss := strings.SplitN(str, "=", 2)
	if len(ss) == 2 {
		return &Env{
			Name:  ss[0],
			Value: ss[1],
		}, nil
	} else {
		return nil, &ErrInvalidEnvString{str}
	}
}
