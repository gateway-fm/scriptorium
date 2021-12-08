package logger

import (
	"fmt"
	"strings"
)

type AppEnv int

const (
	Production AppEnv = iota
	Development
	Local
	Wrong
)

var Envs = [...]string{
	Production:  "prod",
	Development: "dev",
	Local:       "local",
}

func (s AppEnv) String() string {
	return Envs[s]
}

func EnvFromStr(s string) (AppEnv, error) {
	for i, r := range Envs {
		if strings.ToLower(s) == r {
			return AppEnv(i), nil
		}
	}
	return Wrong, fmt.Errorf("wrong env type where provided %q", s)
}
