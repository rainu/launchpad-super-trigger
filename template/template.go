package template

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/rainu/launchpad-super-trigger/sensor"
	"github.com/rainu/launchpad-super-trigger/sensor/data_extractor"
	"go.uber.org/zap"
	"strings"
	"text/template"
)

type Engine struct {
	templates    map[string]*template.Template
	templateData *templateData
}

type Sensor struct {
	Sensor     sensor.Sensor
	Extractors map[string]data_extractor.Extractor
}

type templateData struct {
	sensors map[string]Sensor
}

func NewEngine(sensors map[string]Sensor) *Engine {
	return &Engine{
		templates: map[string]*template.Template{},
		templateData: &templateData{
			sensors: sensors,
		},
	}
}

func (e *Engine) RegisterTemplate(name, templateContent string) error {
	tmpl, err := template.New(name).Parse(templateContent)
	if err != nil {
		return err
	}

	e.templates[name] = tmpl
	return nil
}

func (e *Engine) Execute(name string) ([]byte, error) {
	tmpl := e.templates[name]
	if tmpl == nil {
		return nil, errors.New("template not found")
	}

	buff := &bytes.Buffer{}
	err := tmpl.Execute(buff, e.templateData)

	return buff.Bytes(), err
}

func (t *templateData) DataPoint(dpPath string) string {
	split := strings.Split(dpPath, ".")
	if len(split) != 2 {
		zap.L().Error("Template: datapoint is invalid!")
		return ""
	}

	sensorName := split[0]
	extractorName := split[1]

	s, ok := t.sensors[sensorName]
	if !ok {
		zap.L().Error("Template: datapoint is invalid! Sensor was not found!")
		return ""
	}

	e, ok := s.Extractors[extractorName]
	if !ok {
		zap.L().Error(fmt.Sprintf("Template: datapoint is invalid! Datapoint for Sensor %s was not found!", sensorName))
		return ""
	}

	extracted, err := e.Extract(s.Sensor.LastMessage())
	if err != nil {
		zap.L().Warn("Template: datapoint extraction failed!", zap.Error(err))
		return ""
	}

	return string(extracted)
}
