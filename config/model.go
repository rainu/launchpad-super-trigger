package config

import (
	"time"
)

type Config struct {
	Actors    Actors    `yaml:"actors"`
	Listeners Listeners `yaml:"listeners"`
	Layout    Layout    `yaml:"layout"`
}

type Actors struct {
	Rest     map[string]RestActor     `yaml:"rest,omitempty" validate:"dive"`
	Combined map[string]CombinedActor `yaml:"combined,omitempty" validate:"dive"`

	GfxBlink map[string]GfxBlinkActor `yaml:"gfxBlink,omitempty" validate:"dive"`
}

type RestActor struct {
	Method   string              `yaml:"method"`
	URL      string              `yaml:"url" validate:"required,url"`
	Header   map[string][]string `yaml:"header"`
	BodyB64  string              `yaml:"bodyBase64" validate:"omitempty,base64"`
	BodyPath string              `yaml:"bodyPath" validate:"omitempty,file"`
	BodyRaw  string              `yaml:"body"`
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

type Listeners struct {
}

type Layout struct {
	Pages map[int]Page `yaml:"pages" validate:"dive,keys,gte=0,lte=255,endkeys,required"`
}

type Page struct {
	Trigger map[Coordinate]Trigger `yaml:"trigger" validate:"dive,keys,coord,endkeys,required"`
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
