package utils

import "github.com/kelseyhightower/envconfig"

type Env struct {
	Port        int64  `envconfig:"PORT"`
	DatabaseUrl string `envconfig:"DATABASE_URL"`
}

func GetEnv() Env {
	var env Env
	envconfig.MustProcess("", &env)
	return env
}
