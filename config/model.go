package config

type Config struct {
	Actors    Actors    `yaml:"actors"`
	Listeners Listeners `yaml:"listeners"`
	Layout    Layout    `yaml:"layout"`
}

type Actors struct {
	Rest     map[string]RestActor     `yaml:"rest,omitempty" validate:"dive"`
	Combined map[string]CombinedActor `yaml:"combined,omitempty" validate:"dive"`
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

type Listeners struct {
}

type Layout struct {
	Pages map[int]Page `yaml:"pages" validate:"dive,keys,gte=0,lte=255,endkeys,required"`
}

type Page struct {
	Trigger map[string]Trigger `yaml:"trigger" validate:"dive,keys,coord,endkeys,required"`
}

type Trigger struct {
	Actor         string         `yaml:"actor" validate:"required,actor"`
	ColorSettings *ColorSettings `yaml:"color" validate:"omitempty"`
}

type ColorSettings struct {
	Ready    string `yaml:"ready" validate:"color"`
	Progress string `yaml:"progress" validate:"color"`
	Success  string `yaml:"success" validate:"color"`
	Failed   string `yaml:"failed" validate:"color"`
}
