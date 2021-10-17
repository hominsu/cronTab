package config

import (
	"log"
	"os"
	"reflect"
	"strconv"
)

var (
	Basic = &basic{}

	env = &envDef{
		basicWebPort: "CRON_TAB_BASIC_WEB_PORT",
	}
)

func init() {
	typeOfEnv := reflect.TypeOf(*env)
	valueOfEnv := reflect.ValueOf(*env)

	for i := 0; i < typeOfEnv.NumField(); i++ {
		if value := os.Getenv(valueOfEnv.Field(i).String()); value != "" {
			setEnv(typeOfEnv.Field(i).Name, value)
		} else {
			envError("Env should not be empty")
		}
	}
}

func setEnv(key, value string) {
	switch key {
	case "basicWebPort":
		Basic.webPort = value
	default:
		envError("Unknown Env Name")
	}
}

func strToInt(value string) int {
	ret, err := strconv.Atoi(value)
	if err != nil {
		envError(err.Error())
	}
	return ret
}

func envError(err string) {
	log.Fatalln(err)
}

type envDef struct {
	basicWebPort string
}

type basic struct {
	webPort string
}

func (b basic) WebPort() string {
	return b.webPort
}
