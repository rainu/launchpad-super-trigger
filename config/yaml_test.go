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
