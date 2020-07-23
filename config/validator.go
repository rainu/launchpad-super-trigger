package config

import (
	"gopkg.in/go-playground/validator.v9"
	"reflect"
)

func validate(cfg *Config) error {
	configValidator := validator.New()
	err := configValidator.RegisterValidation("coord", func(fl validator.FieldLevel) bool {
		sVal := fl.Field().String()
		_, _, err := Coordinate(sVal).Coordinate()

		return err == nil
	})
	if err != nil {
		panic(err)
	}
	err = configValidator.RegisterValidation("color", func(fl validator.FieldLevel) bool {
		sVal := fl.Field().String()
		_, err := Color(sVal).Color()

		return err == nil
	})
	if err != nil {
		panic(err)
	}
	err = configValidator.RegisterValidation("actor", func(fl validator.FieldLevel) bool {
		actorName := fl.Field().String()

		refActor := reflect.ValueOf(cfg.Actors)
		knownActors := map[string]bool{}

		for i := 0; i < refActor.NumField(); i++ {
			if refActor.Field(i).Kind() == reflect.Map {
				for _, a := range refActor.Field(i).MapKeys() {
					knownActors[a.String()] = true
				}
			}
		}

		if _, found := knownActors[actorName]; found {
			return true
		}
		return false
	})
	if err != nil {
		panic(err)
	}

	return configValidator.Struct(cfg)
}
