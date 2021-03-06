package config

import (
	"errors"
	"fmt"
	"gopkg.in/go-playground/validator.v9"
	"reflect"
	"regexp"
	"strings"
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
	err = configValidator.RegisterValidation("coords", func(fl validator.FieldLevel) bool {
		sVal := fl.Field().String()
		_, err := Coordinates(sVal).Coordinates()

		return err == nil
	})
	if err != nil {
		panic(err)
	}
	err = configValidator.RegisterValidation("color", func(fl validator.FieldLevel) bool {
		sVal := fl.Field().String()
		if sVal == "" {
			return true
		}

		_, err := Color(sVal).Color()

		return err == nil
	})
	if err != nil {
		panic(err)
	}
	err = configValidator.RegisterValidation("pagenumber", func(fl validator.FieldLevel) bool {
		sVal := PageNumber(fl.Field().String())
		return sVal.isValid()
	})
	if err != nil {
		panic(err)
	}
	err = configValidator.RegisterValidation("regexPattern", func(fl validator.FieldLevel) bool {
		sVal := fl.Field().String()
		_, err := regexp.Compile(sVal)
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

	nameRegex := regexp.MustCompile(`^[0-9a-zA-Z_-]+$`)
	err = configValidator.RegisterValidation("component_name", func(fl validator.FieldLevel) bool {
		name := fl.Field().String()

		return nameRegex.MatchString(name)
	})
	if err != nil {
		panic(err)
	}
	err = configValidator.RegisterValidation("connection_mqtt", func(fl validator.FieldLevel) bool {
		connectionName := fl.Field().String()
		knownConnections := cfg.Connections.MQTT

		if _, found := knownConnections[connectionName]; found {
			return true
		}
		return false
	})
	if err != nil {
		panic(err)
	}
	err = configValidator.RegisterValidation("datapoint", func(fl validator.FieldLevel) bool {
		dpPath := Datapoint(fl.Field().String())

		if !dpPath.IsValid() {
			return false
		}

		knownDatapointPaths := map[string]bool{}

		//meta datapoints:
		refConnections := reflect.ValueOf(cfg.Connections)
		for connectionField := 0; connectionField < refConnections.NumField(); connectionField++ {
			if refConnections.Field(connectionField).Kind() == reflect.Map {
				//connections.mqtt
				for _, connectionName := range refConnections.Field(connectionField).MapKeys() {
					knownDatapointPaths[fmt.Sprintf("__connection[%s].status", connectionName.String())] = true
				}
			}
		}

		refSensors := reflect.ValueOf(cfg.Sensors)
		for sensorField := 0; sensorField < refSensors.NumField(); sensorField++ {
			if refSensors.Field(sensorField).Kind() == reflect.Map {
				//sensors.mqtt
				for _, sensorTypeName := range refSensors.Field(sensorField).MapKeys() {
					sensorValue := refSensors.Field(sensorField).MapIndex(sensorTypeName)
					for sensorValueField := 0; sensorValueField < sensorValue.NumField(); sensorValueField++ {
						//sensors.mqtt.data
						if sensorValue.Field(sensorValueField).Kind() == reflect.ValueOf(DataPoints{}).Kind() {
							for sensorValueDataPoint := 0; sensorValueDataPoint < sensorValue.Field(sensorValueField).NumField(); sensorValueDataPoint++ {
								if sensorValue.Field(sensorValueField).Field(sensorValueDataPoint).Kind() == reflect.Map {
									//sensors.mqtt.data.gjson
									for _, dpName := range sensorValue.Field(sensorValueField).Field(sensorValueDataPoint).MapKeys() {

										dpPath := fmt.Sprintf("%s.%s", sensorTypeName.String(), dpName.String())
										knownDatapointPaths[dpPath] = true
									}
								}
							}

							completeField := sensorValue.Field(sensorValueField).FieldByName("Complete")
							if completeField.String() != "" {
								dpPath := fmt.Sprintf("%s.%s", sensorTypeName.String(), completeField.String())
								knownDatapointPaths[dpPath] = true
							}
						}
					}
				}
			}
		}

		return knownDatapointPaths[dpPath.Path()]
	})
	if err != nil {
		panic(err)
	}

	///////
	//call validation
	validationError := configValidator.Struct(cfg)
	joinedErrors := strings.Builder{}

	if validationError != nil {
		joinedErrors.WriteString(validationError.Error())
	}

	///////
	// do dome extra validation

	//unique actors
	refActor := reflect.ValueOf(cfg.Actors)
	knownActors := map[string]bool{}

	for i := 0; i < refActor.NumField(); i++ {
		if refActor.Field(i).Kind() == reflect.Map {
			for _, a := range refActor.Field(i).MapKeys() {
				if knownActors[a.String()] {
					joinedErrors.WriteString(fmt.Sprintf(`actor '%s' already known`, a.String()))
				} else {
					knownActors[a.String()] = true
				}
			}
		}
	}

	//unique sensors
	refSensors := reflect.ValueOf(cfg.Sensors)
	knownSensors := map[string]map[string]bool{}

	for sensorField := 0; sensorField < refSensors.NumField(); sensorField++ {
		if refSensors.Field(sensorField).Kind() == reflect.Map {
			for _, sensorTypeName := range refSensors.Field(sensorField).MapKeys() {
				if _, found := knownSensors[sensorTypeName.String()]; found {
					joinedErrors.WriteString(fmt.Sprintf(`sensor '%s' already known`, sensorTypeName.String()))
				} else {
					knownSensors[sensorTypeName.String()] = map[string]bool{}
				}

				sensorValue := refSensors.Field(sensorField).MapIndex(sensorTypeName)
				for sensorValueField := 0; sensorValueField < sensorValue.NumField(); sensorValueField++ {
					if sensorValue.Field(sensorValueField).Kind() == reflect.ValueOf(DataPoints{}).Kind() {
						for sensorValueDataPoint := 0; sensorValueDataPoint < sensorValue.Field(sensorValueField).NumField(); sensorValueDataPoint++ {
							if sensorValue.Field(sensorValueField).Field(sensorValueDataPoint).Kind() == reflect.Map {
								for _, dpName := range sensorValue.Field(sensorValueField).Field(sensorValueDataPoint).MapKeys() {

									if _, found := knownSensors[sensorTypeName.String()][dpName.String()]; found {
										joinedErrors.WriteString(fmt.Sprintf(`sensor data point '%s' already known`, dpName.String()))
									}

									knownSensors[sensorTypeName.String()][dpName.String()] = true
								}
							}
						}
					}
				}
			}
		}
	}

	if joinedErrors.Len() > 0 {
		return errors.New(joinedErrors.String())
	}

	return validationError
}
