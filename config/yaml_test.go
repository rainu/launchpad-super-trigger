package config

import (
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
	"time"
)

func TestReadConfigFile(t *testing.T) {
	tests := []struct {
		name           string
		content        string
		expectedResult Config
		expectedErrors []string
	}{
		{
			`simple rest`,
			`
actors:
	rest:
		test:
			method: POST
			url: "http://localhost:1312"
			header: 
				"Content-Type": 
					- "application/json"
			body: "Hello World"`,
			Config{
				Actors: Actors{
					Rest: map[string]RestActor{
						"test": {
							Method: "POST",
							URL:    "http://localhost:1312",
							Header: map[string][]string{
								"Content-Type": {"application/json"},
							},
							BodyRaw: `Hello World`,
						},
					},
				},
			},
			[]string{},
		},
		{
			`validation rest`,
			`
actors:
	rest:
		test:
			url: nonValidUrl
			bodyPath: invalidPath
			bodyBase64: "^"`,
			Config{},
			[]string{
				`Key: 'Config.Actors.Rest[test].URL'`,
				`Key: 'Config.Actors.Rest[test].BodyPath'`,
				`Key: 'Config.Actors.Rest[test].BodyB64'`,
			},
		},
		{
			`simple connection`,
			`
connections:
	mqtt:
		b1:
			broker: tcp://mqtt-broker:1883
			clientId: client
			username: user
			password: password`,
			Config{
				Connections: Connections{
					MQTT: map[string]MQTTConnection{
						"b1": {
							Broker:   "tcp://mqtt-broker:1883",
							ClientId: "client",
							Username: "user",
							Password: "password",
						},
					},
				},
			},
			[]string{},
		},
		{
			`validation connection`,
			`
connections:
	mqtt:
		b1:
			broker: abc`,
			Config{},
			[]string{
				`Key: 'Config.Connections.MQTT[b1].Broker'`,
			},
		},
		{
			`simple mqtt`,
			`
connections:
	mqtt:
		b1:
			broker: tcp://mqtt-broker:1883
actors:
	mqtt:
		test:
			connection: b1
			topic: test/topic
			qos: 1
			retained: true
			body: "Hello World"`,
			Config{
				Connections: Connections{
					MQTT: map[string]MQTTConnection{
						"b1": {
							Broker: "tcp://mqtt-broker:1883",
						},
					},
				},
				Actors: Actors{
					Mqtt: map[string]MQTTActor{
						"test": {
							Connection: "b1",
							Topic:      "test/topic",
							QOS:        1,
							Retained:   true,
							BodyRaw:    `Hello World`,
						},
					},
				},
			},
			[]string{},
		},
		{
			`validation mqtt`,
			`
actors:
	mqtt:
		test:
			connection: b1
			qos: 4
			bodyPath: invalidPath
			bodyBase64: "^"`,
			Config{},
			[]string{
				`Key: 'Config.Actors.Mqtt[test].Topic'`,
				`Key: 'Config.Actors.Mqtt[test].QOS'`,
				`Key: 'Config.Actors.Mqtt[test].BodyB64'`,
				`Key: 'Config.Actors.Mqtt[test].BodyPath'`,
			},
		},
		{
			`simple combined`,
			`
actors:
	rest:
		test:
			url: "http://localhost:1312"
	combined:
		c-test:
			actors:
				- test
				- test`,
			Config{
				Actors: Actors{
					Rest: map[string]RestActor{
						"test": {
							URL: "http://localhost:1312",
						},
					},
					Combined: map[string]CombinedActor{
						"c-test": {
							Actor:    []string{"test", "test"},
							Parallel: false,
						},
					},
				},
			},
			[]string{},
		},
		{
			`combined validation`,
			`
actors:
	rest:
		test:
			url: "http://localhost:1312"
	combined:
		c-test:
			actors:
				- test
				- doesNotExists
		c-test2:
			actors:
				- test`,
			Config{
				Actors: Actors{
					Rest: map[string]RestActor{
						"test": {
							URL: "http://localhost:1312",
						},
					},
					Combined: map[string]CombinedActor{
						"c-test": {
							Actor:    []string{"test", "test"},
							Parallel: false,
						},
					},
				},
			},
			[]string{
				`Key: 'Config.Actors.Combined[c-test].Actor[1]'`,
				`Key: 'Config.Actors.Combined[c-test2].Actor'`,
			},
		},
		{
			`simple gfx blink`,
			`
actors:
	gfxBlink:
		bl:
			on: 1,1
			off: 3,3
			interval: 500ms
			duration: 10s`,
			Config{
				Actors: Actors{
					GfxBlink: map[string]GfxBlinkActor{
						"bl": {
							ColorOn:  "1,1",
							ColorOff: "3,3",
							Interval: 500 * time.Millisecond,
							Duration: 10 * time.Second,
						},
					},
				},
			},
			[]string{},
		},
		{
			`validate gfx blink`,
			`
actors:
	gfxBlink:
		bl:
			on: 1,
			off: 3,9`,
			Config{},
			[]string{
				`Key: 'Config.Actors.GfxBlink[bl].ColorOn'`,
				`Key: 'Config.Actors.GfxBlink[bl].ColorOff'`,
			},
		},
		{
			`simple gfx wave`,
			`
actors:
	gfxWave:
		wv:
			square: true
			color: 1,1
			delay: 50ms`,
			Config{
				Actors: Actors{
					GfxWave: map[string]GfxWaveActor{
						"wv": {
							Square: true,
							Color:  "1,1",
							Delay:  50 * time.Millisecond,
						},
					},
				},
			},
			[]string{},
		},
		{
			`validate gfx blink`,
			`
actors:
	gfxWave:
		wv:
			color: 1,`,
			Config{},
			[]string{
				`Key: 'Config.Actors.GfxWave[wv].Color'`,
			},
		},
		{
			`simple command`,
			`
actors:
	command:
		cmd:
			name: /bin/test.sh
			args: 
				- -h
			appendContext: true`,
			Config{
				Actors: Actors{
					Command: map[string]CommandActor{
						"cmd": {
							Name:          "/bin/test.sh",
							Arguments:     []string{"-h"},
							AppendContext: true,
						},
					},
				},
			},
			[]string{},
		},
		{
			`validate command`,
			`
actors:
	command:
		cmd:`,
			Config{},
			[]string{
				`Key: 'Config.Actors.Command[cmd].Name'`,
			},
		},
		{
			`simple layout`,
			`
actors:
	rest:
		test:
			url: "http://localhost:1312"
layout:
	pages:
		0:
			trigger:
				"1,2":
					actor: test
					color:
						ready: 0,0
						progress: 1,1
						success: 2,2
						failed: 3,3`,
			Config{
				Actors: Actors{
					Rest: map[string]RestActor{
						"test": {
							URL: "http://localhost:1312",
						},
					},
				},
				Layout: Layout{
					Pages: map[int]Page{
						0: {
							Trigger: map[Coordinates]Trigger{
								"1,2": {
									Actor: "test",
									ColorSettings: &ColorSettings{
										Ready:    "0,0",
										Progress: "1,1",
										Success:  "2,2",
										Failed:   "3,3",
									},
								},
							},
						},
					},
				}},
			[]string{},
		},
		{
			`layout validation`,
			`
layout:
	pages:
		999:
			trigger:
				"-1,2":
					actor: test
					color:
						ready: 0
						progress: 1,4
						success: 4,1
						failed: -1,3`,
			Config{},
			[]string{
				`Key: 'Config.Layout.Pages[999]'`,
				`Key: 'Config.Layout.Pages[999].Trigger[-1,2]'`,
				`Key: 'Config.Layout.Pages[999].Trigger[-1,2].Actor'`,
				`Key: 'Config.Layout.Pages[999].Trigger[-1,2].ColorSettings.Ready'`,
				`Key: 'Config.Layout.Pages[999].Trigger[-1,2].ColorSettings.Progress'`,
				`Key: 'Config.Layout.Pages[999].Trigger[-1,2].ColorSettings.Success'`,
				`Key: 'Config.Layout.Pages[999].Trigger[-1,2].ColorSettings.Failed'`,
			},
		},
	}
	for _, tt := range tests {
		t.Run("TestReadConfigFile_"+tt.name, func(t *testing.T) {
			tt.content = strings.ReplaceAll(tt.content, "\t", " ")

			parsedConf, err := ReadConfig(strings.NewReader(tt.content))

			if len(tt.expectedErrors) > 0 {
				assert.Error(t, err)

				for _, expectedError := range tt.expectedErrors {
					assert.Contains(t, err.Error(), expectedError)
				}
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedResult, *parsedConf)
			}
		})
	}
}
