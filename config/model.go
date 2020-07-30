package config

import (
	"github.com/rainu/launchpad-super-trigger/gfx"
	"time"
)

type Config struct {
	Connections Connections `yaml:"connections"`
	Actors      Actors      `yaml:"actors"`
	Sensors     Sensors     `yaml:"sensors"`
	Layout      Layout      `yaml:"layout"`
}

type Connections struct {
	MQTT map[string]MQTTConnection `yaml:"mqtt" validate:"dive,keys,component_name,endkeys,required"`
}

type MQTTConnection struct {
	Broker   string `yaml:"broker" validate:"required,url"`
	ClientId string `yaml:"clientId"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

type Actors struct {
	Rest        map[string]RestActor        `yaml:"rest,omitempty" validate:"dive,keys,component_name,endkeys,required"`
	Mqtt        map[string]MQTTActor        `yaml:"mqtt,omitempty" validate:"dive,keys,component_name,endkeys,required"`
	Conditional map[string]ConditionalActor `yaml:"conditional,omitempty" validate:"dive,keys,component_name,endkeys,required"`
	Command     map[string]CommandActor     `yaml:"command,omitempty" validate:"dive,keys,component_name,endkeys,required"`
	Combined    map[string]CombinedActor    `yaml:"combined,omitempty" validate:"dive,keys,component_name,endkeys,required"`

	GfxBlink map[string]GfxBlinkActor `yaml:"gfxBlink,omitempty" validate:"dive,keys,component_name,endkeys,required"`
	GfxWave  map[string]GfxWaveActor  `yaml:"gfxWave,omitempty" validate:"dive,keys,component_name,endkeys,required"`
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

type ConditionalActor struct {
	Conditions []Condition `yaml:"conditions" validate:"dive,required"`
}

type Condition struct {
	Actor      string              `yaml:"actor" validate:"required,actor"`
	DataPoint  Datapoint           `yaml:"datapoint" validate:"required,datapoint"`
	Expression ConditionExpression `yaml:"expression" validate:"required"`
}

type ConditionExpression struct {
	Eq          *string  `yaml:"eq"`
	Ne          *string  `yaml:"ne"`
	Lt          *float64 `yaml:"lt"`
	Lte         *float64 `yaml:"lte"`
	Gt          *float64 `yaml:"gt"`
	Gte         *float64 `yaml:"gte"`
	Match       *string  `yaml:"match" validate:"omitempty,regexPattern"`
	NotMatch    *string  `yaml:"nmatch" validate:"omitempty,regexPattern"`
	Contains    *string  `yaml:"contains"`
	NotContains *string  `yaml:"ncontains"`
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

type Sensors struct {
	Mqtt map[string]MQTTSensor `yaml:"mqtt,omitempty" validate:"dive,keys,component_name,endkeys,required"`
	Rest map[string]RestSensor `yaml:"rest,omitempty" validate:"dive,keys,component_name,endkeys,required"`
}

type MQTTSensor struct {
	Connection string     `yaml:"connection" validate:"required,connection_mqtt"`
	Topic      string     `yaml:"topic" validate:"required"`
	QOS        byte       `yaml:"qos" validate:"gte=0,lte=2"`
	DataPoints DataPoints `yaml:"data" validate:"required"`
}

type RestSensor struct {
	Method     string              `yaml:"method"`
	URL        string              `yaml:"url" validate:"required,url"`
	Header     map[string][]string `yaml:"header"`
	BodyB64    string              `yaml:"bodyBase64" validate:"omitempty,base64"`
	BodyPath   string              `yaml:"bodyPath" validate:"omitempty,file"`
	BodyRaw    string              `yaml:"body"`
	Interval   time.Duration       `yaml:"interval" validate:"required"`
	DataPoints DataPoints          `yaml:"data" validate:"required"`
}

type DataPoints struct {
	Complete string                    `yaml:"complete"`
	Gjson    map[string]string         `yaml:"gjson" validate:"dive,keys,component_name,endkeys,required"`
	Gojq     map[string]string         `yaml:"gojq" validate:"dive,keys,component_name,endkeys,required"`
	Split    map[string]SplitDataPoint `yaml:"split" validate:"dive,keys,component_name,endkeys,required"`
}

type SplitDataPoint struct {
	Separator string `yaml:"separator" validate:"required"`
	Index     int    `yaml:"index" validate:"gte=0"`
}

type Layout struct {
	Pages map[int]Page `yaml:"pages" validate:"dive,keys,gte=0,lte=255,endkeys,required"`
}

type Page struct {
	Trigger map[Coordinates]Trigger `yaml:"trigger" validate:"dive,keys,coords,endkeys,required"`
	Plotter Plotters                `yaml:"plotter" validate:"dive"`
}

type Trigger struct {
	Actor         string         `yaml:"actor" validate:"required,actor"`
	ColorSettings *ColorSettings `yaml:"color" validate:"omitempty"`
}

type ColorSettings struct {
	Disable  bool  `yaml:"disable"`
	Ready    Color `yaml:"ready" validate:"color"`
	Progress Color `yaml:"progress" validate:"color"`
	Success  Color `yaml:"success" validate:"color"`
	Failed   Color `yaml:"failed" validate:"color"`
}

type Plotters struct {
	Progressbar []Progressbar `yaml:"progressbar" validate:"dive"`
	Static      []Static      `yaml:"static" validate:"dive"`
}

type Progressbar struct {
	DataPoint   Datapoint    `yaml:"datapoint" validate:"required,datapoint"`
	X           int          `yaml:"x"`
	Y           int          `yaml:"y"`
	Min         float64      `yaml:"min"`
	Max         *float64     `yaml:"max"`
	Vertical    bool         `yaml:"vertical"`
	Quadrant    gfx.Quadrant `yaml:"quadrant"`
	RightToLeft bool         `yaml:"rtl"`
	Fill        Color        `yaml:"fill" validate:"color"`
	Empty       Color        `yaml:"empty" validate:"color"`
}

type Static struct {
	DataPoint    Datapoint         `yaml:"datapoint" validate:"required,datapoint"`
	Position     Coordinate        `yaml:"pos" validate:"coord,required"`
	Expressions  StaticExpressions `yaml:"expressions" validate:"required"`
	DefaultColor *Color            `yaml:"defaultColor" validate:"omitempty,color"`
}

type StaticExpressions struct {
	Eq          []StaticExpression        `yaml:"eq" validate:"dive"`
	Ne          []StaticExpression        `yaml:"ne" validate:"dive"`
	Lt          []StaticNumericExpression `yaml:"lt" validate:"dive"`
	Lte         []StaticNumericExpression `yaml:"lte" validate:"dive"`
	Gt          []StaticNumericExpression `yaml:"gt" validate:"dive"`
	Gte         []StaticNumericExpression `yaml:"gte" validate:"dive"`
	Match       []StaticMatchExpression   `yaml:"match" validate:"dive"`
	NotMatch    []StaticMatchExpression   `yaml:"nmatch" validate:"dive"`
	Contains    []StaticExpression        `yaml:"contains" validate:"dive"`
	NotContains []StaticExpression        `yaml:"ncontains" validate:"dive"`
}

type StaticExpression struct {
	Value           string `yaml:"value" validate:"required"`
	ActivationColor Color  `yaml:"color" validate:"color,required"`
}

type StaticMatchExpression struct {
	Value           string `yaml:"value" validate:"required,regexPattern"`
	ActivationColor Color  `yaml:"color" validate:"color,required"`
}

type StaticNumericExpression struct {
	Value           float64 `yaml:"value" validate:"required"`
	ActivationColor Color   `yaml:"color" validate:"color,required"`
}
