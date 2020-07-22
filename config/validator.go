package config

import (
	"gopkg.in/go-playground/validator.v9"
	"strconv"
	"strings"
)

func validate(cfg *Config) error {
	configValidator := validator.New()
	err := configValidator.RegisterValidation("coord", func(fl validator.FieldLevel) bool {
		sVal := fl.Field().String()
		split := strings.Split(sVal, ",")
		if len(split) != 2 {
			return false
		}
		x, err := strconv.Atoi(split[0])
		if err != nil {
			return false
		}
		if x < 0 || x > 7 {
			return false
		}
		y, err := strconv.Atoi(split[1])
		if err != nil {
			return false
		}
		if y < 0 || y > 7 {
			return false
		}

		return true
	})
	if err != nil {
		panic(err)
	}
	err = configValidator.RegisterValidation("color", func(fl validator.FieldLevel) bool {
		sVal := fl.Field().String()
		split := strings.Split(sVal, ",")
		if len(split) != 2 {
			return false
		}
		r, err := strconv.Atoi(split[0])
		if err != nil {
			return false
		}
		if r < 0 || r > 3 {
			return false
		}
		g, err := strconv.Atoi(split[1])
		if err != nil {
			return false
		}
		if g < 0 || g > 3 {
			return false
		}

		return true
	})
	if err != nil {
		panic(err)
	}
	err = configValidator.RegisterValidation("actor", func(fl validator.FieldLevel) bool {
		actorName := fl.Field().String()
		_, found := cfg.Actors.Rest[actorName]

		return found
	})
	if err != nil {
		panic(err)
	}

	return configValidator.Struct(cfg)
}
