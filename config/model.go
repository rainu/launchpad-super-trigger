package config

import (
	"time"
)

type Config struct {
	Connections Connections `yaml:"connections"`
	Actors      Actors      `yaml:"actors"`
	Listeners   Listeners   `yaml:"listeners"`
	Layout      Layout      `yaml:"layout"`
}

type Connections struct {
	MQTT map[string]MQTTConnection `yaml:"mqtt" validate:"dive"`
}

type MQTTConnection struct {
	Broker   string `yaml:"broker" validate:"required,url"`
	ClientId string `yaml:"clientId"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

type Actors struct {
	Rest     map[string]RestActor     `yaml:"rest,omitempty" validate:"dive"`
	Mqtt     map[string]MQTTActor     `yaml:"mqtt,omitempty" validate:"dive"`
	Command  map[string]CommandActor  `yaml:"command,omitempty" validate:"dive"`
	Combined map[string]CombinedActor `yaml:"combined,omitempty" validate:"dive"`

	GfxBlink map[string]GfxBlinkActor `yaml:"gfxBlink,omitempty" validate:"dive"`
	GfxWave  map[string]GfxWaveActor  `yaml:"gfxWave,omitempty" validate:"dive"`
}

type RestActor struct {
	Method   string              `yaml:"method"`
	URL      string              `yaml:"url" validate:"required,url"`
	Header   map[string][]string `yaml:"header"`
	BodyB64  string              `yaml:"bodyBase64" validate:"omitempty,base64"`
	BodyPath string              `yaml:"bodyPath" validate:"omitempty,file"`
	BodyRaw  string              `yaml:"body"`
}

type MQTTActor struct {
	Connection string `yaml:"connection" validate:"required,connection_mqtt"`
	Topic      string `yaml:"topic" validate:"required"`
	QOS        byte   `yaml:"qos" validate:"gte=0,lte=2"`
	Retained   bool   `yaml:"retained"`
	BodyB64    string `yaml:"bodyBase64" validate:"omitempty,base64"`
	BodyPath   string `yaml:"bodyPath" validate:"omitempty,file"`
	BodyRaw    string `yaml:"body"`
}

type CommandActor struct {
	Name          string   `yaml:"name" validate:"required"`
	Arguments     []string `yaml:"args"`
	AppendContext bool     `yaml:"appendContext"`
}

type CombinedActor struct {
	Actor    []string `yaml:"actors" validate:"gte=2,dive,actor"`
	Parallel bool     `yaml:"parallel"`
}

type GfxBlinkActor struct {
	ColorOn  Color         `yaml:"on" validate:"color,required"`
	ColorOff Color         `yaml:"off" validate:"color"`
	Interval time.Duration `yaml:"interval"`
	Duration time.Duration `yaml:"duration"`
}

type GfxWaveActor struct {
	Square bool          `yaml:"square"`
	Color  Color         `yaml:"color" validate:"color"`
	Delay  time.Duration `yaml:"delay"`
}

type Listeners struct {
}

type Layout struct {
	Pages map[int]Page `yaml:"pages" validate:"dive,keys,gte=0,lte=255,endkeys,required"`
}

type Page struct {
	Trigger map[Coordinates]Trigger `yaml:"trigger" validate:"dive,keys,coords,endkeys,required"`
}

type Trigger struct {
	Actor         string         `yaml:"actor" validate:"required,actor"`
	ColorSettings *ColorSettings `yaml:"color" validate:"omitempty"`
}

type ColorSettings struct {
	Ready    Color `yaml:"ready" validate:"color"`
	Progress Color `yaml:"progress" validate:"color"`
	Success  Color `yaml:"success" validate:"color"`
	Failed   Color `yaml:"failed" validate:"color"`
}
