package actor

import (
	"github.com/rainu/launchpad-super-trigger/actor"
	"github.com/rainu/launchpad-super-trigger/config"
	"github.com/rainu/launchpad-super-trigger/config/expressions"
	configSensor "github.com/rainu/launchpad-super-trigger/config/sensor"
	"github.com/rainu/launchpad-super-trigger/sensor"
	"github.com/rainu/launchpad-super-trigger/sensor/data_extractor"
)

func buildConditional(actors map[string]actor.Actor, sensors map[string]configSensor.Sensor, conditionalActors map[string]config.ConditionalActor) {
	for actorName, conditionalActor := range conditionalActors {
		handler := &actor.Conditional{}

		for _, condition := range conditionalActor.Conditions {
			actor := actors[condition.Actor]
			sensor, extractor := buildDatapoint(condition.DataPoint, sensors)
			expr := buildExpression(condition.Expression)

			handler.AddCondition(actor, sensor, extractor, expr)
		}

		actors[actorName] = handler
	}
}

func buildDatapoint(dpPath config.Datapoint, sensors map[string]configSensor.Sensor) (sensor.Sensor, data_extractor.Extractor) {
	sensorName := dpPath.Sensor()
	dpName := dpPath.Name()

	return sensors[sensorName].Sensor, sensors[sensorName].Extractors[dpName]
}

func buildExpression(expr config.ConditionExpression) func([]byte) bool {
	if expr.Eq != nil {
		return expressions.BuildEqExpressionFn(*expr.Eq)
	}
	if expr.Ne != nil {
		return expressions.BuildNeExpressionFn(*expr.Ne)
	}
	if expr.Lt != nil {
		return expressions.BuildLtExpressionFn(*expr.Lt)
	}
	if expr.Lte != nil {
		return expressions.BuildLteExpressionFn(*expr.Lte)
	}
	if expr.Gt != nil {
		return expressions.BuildGtExpressionFn(*expr.Gt)
	}
	if expr.Gte != nil {
		return expressions.BuildGteExpressionFn(*expr.Gte)
	}
	if expr.Match != nil {
		return expressions.BuildMatchExpressionFn(*expr.Match)
	}
	if expr.NotMatch != nil {
		return expressions.BuildNotMatchExpressionFn(*expr.NotMatch)
	}
	if expr.Contains != nil {
		return expressions.BuildContainsExpressionFn(*expr.Contains)
	}
	if expr.NotContains != nil {
		return expressions.BuildNotContainsExpressionFn(*expr.NotContains)
	}

	panic("No expression given!")
}
