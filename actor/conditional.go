package actor

import (
	"github.com/rainu/launchpad-super-trigger/sensor"
	"github.com/rainu/launchpad-super-trigger/sensor/data_extractor"
)

type Conditional struct {
	conditions []condition
}

type condition struct {
	actor      Actor
	sensor     sensor.Sensor
	extractor  data_extractor.Extractor
	expression func([]byte) bool
}

func (c *Conditional) AddCondition(
	actor Actor,
	sensor sensor.Sensor,
	extractor data_extractor.Extractor,
	expression func([]byte) bool) {

	c.conditions = append(c.conditions, condition{
		actor:      actor,
		sensor:     sensor,
		extractor:  extractor,
		expression: expression,
	})
}

func (c *Conditional) Do(ctx Context) error {
	for _, condition := range c.conditions {
		lm := condition.sensor.LastMessage()
		dp, err := condition.extractor.Extract(lm)
		if err != nil {
			return err
		}

		//expression is fulfilled
		if condition.expression(dp) {
			//call corresponding actor
			return condition.actor.Do(ctx)
		}
	}

	//no condition was fulfilled -> nothing to do
	return nil
}
