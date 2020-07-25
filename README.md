# launchpad-super-trigger
Trigger application for the Novation Launchpad S

# Install
Portmidi is required to use this package.

```bash
$ apt-get install libportmidi-dev
# or
$ brew install portmidi
# or 
$ yay -S portmidi
```

# Get the Binary
You can build it on your own (you will need [golang](https://golang.org/) installed):
```bash
go build -a -installsuffix cgo ./cmd/lst/
```

# ConfigFile

```yaml
connections:
 mqtt:
  test-broker:
   broker: tcp://broker-host:1883
   clientId: launchpad
actors:
 rest:
  rest-test:
   url: "http://localhost:1312"
 mqtt:
  mqtt-test:
   connection: test-broker
   topic: my/topic
   body: Hello!
layout:
 pages:
  0:
   trigger:
    "0,1":
     actor: rest-test
     color:
      ready: 0,0
      progress: 1,1
      success: 2,2
      failed: 3,3
    "1,0":
     actor: mqtt-test
```

## Config documentation

| config key | default | mandatory | description |
|---|---|---|---|
| actors | - | false | Contains all available actors. |
| actors.command | - | false | Contains all available command actors. A command actor will runs a command on the local machine. |
| actors.command[*name*].name | - | **true** | The name/path of the command to execute. |
| actors.command[*name*].args | - | false | A list of arguments which will be send to the command. |
| actors.command[*name*].appendContext | false | false | Should the trigger context append in the command args? |
| actors.rest | - | false | Contains all available rest actors. An rest actor will call an rest service. |
| actors.rest[*name*].method | GET | false | The http method which should be used. |
| actors.rest[*name*].url | - | **true** | The target url. |
| actors.rest[*name*].header[*name*][]*value* | - | false | The http headers which should be used. |
| actors.rest[*name*].body | - | false | The body content for the http request. |
| actors.rest[*name*].bodyBase64 | - | false | The body content for the http request encoded in base64. |
| actors.rest[*name*].bodyPath | - | false | The body file for the http request. |
| actors.mqtt | - | false | Contains all available mqtt actors. An mqtt actor will send a message to a given topic. |
| actors.mqtt[*name*].connection | - | false | The name of to mqtt connection which should be used for this mqtt actor. |
| actors.mqtt[*name*].topic | - | **true** | The topic name. |
| actors.mqtt[*name*].qos | 0 | false | The QualityOfService for the published message. |
| actors.mqtt[*name*].retained | false | false | Should the message be retained? |
| actors.mqtt[*name*].body | - | false | The body content for the mqtt message. |
| actors.mqtt[*name*].bodyBase64 | - | false | The body content for the mqtt message encoded in base64. |
| actors.mqtt[*name*].bodyPath | - | false | The body file for the mqtt message. |
| actors.combined | - | false | Contains all available combined actors. An combined actor will call other actors (sequential or in parallel). |
| actors.combined[*name*].actors | - | **true** | The list of underlying actor names. Must be greater or equal than 2! |
| actors.combined[*name*].parallel | false | false | How should the underlying actors be called. If true they will be called parallel. Otherwise the will be called sequential. |
| actors.gfxBlink | - | false | Contains all available gfx blink actors. A blink actor will draw blinking pads. |
| actors.gfxBlink[*name*].on | - | true | The on *color*. |
| actors.gfxBlink[*name*].off | 0,0 | false | The off *color*. |
| actors.gfxBlink[*name*].interval | 1s | false | The interval of blink animation. |
| actors.gfxBlink[*name*].duration | until page leave | false | The duration of the blink interval. |
| actors.gfxWave | - | false | Contains all available gfx wave actors. A wave actor will draw waves on the pads. |
| actors.gfxWave[*name*].square | false | false | Should the waveform be square? |
| actors.gfxWave[*name*].color | 0,3 | false | The color of the wave. |
| actors.gfxWave[*name*].delay | 500ms | false | The delay between wave steps. |
| connections | - | false | Contains all available connections, which can be used by different actors. |
| connections.mqtt | - | false | Contains all available MQTT connections. |
| connections.mqtt[*name*].broker | - | **true** | The url of the MQTT-Broker. |
| connections.mqtt[*name*].clientId | - | false | The client id which is send to the mqtt broker. |
| connections.mqtt[*name*].username | - | false | The username. |
| connections.mqtt[*name*].password | - | false | The password. |
| layout | - | false | Contains all layout settings. |
| layout.pages | - | false | Contains all page settings. |
| layout.pages[*pageNumber*] | - | false | Contains a page setting. |
| layout.pages[*pageNumber*].trigger | - | false | Contains settings about the trigger on this page. |
| layout.pages[*pageNumber*].trigger[*coordinates*] | - | false | Contains settings about the trigger which should be called if the given button at *coordinate* is hit. |
| layout.pages[*pageNumber*].trigger[*coordinates*].actor | - | **true** | The name of the actor (see actor config) which should be called if the trigger is hit. |
| layout.pages[*pageNumber*].trigger[*coordinates*].color | - | false | Contains the *color* settings about the trigger. |
| layout.pages[*pageNumber*].trigger[*coordinates*].color.ready | - | false | The color which should be used if the trigger is ready. |
| layout.pages[*pageNumber*].trigger[*coordinates*].color.progress | - | false | The color which should be used as long as the actor is in progress. |
| layout.pages[*pageNumber*].trigger[*coordinates*].color.success | - | false | The color which should be used if the actor work was done successfully. |
| layout.pages[*pageNumber*].trigger[*coordinates*].color.failed | - | false | The color which should be used if the actor work was done wrong. |

* *pageNumber*
    * The page number must be a number from **0** until **255**.
* *coordinate*
    * "X,Y" -> X and Y must have a value of 0-7.
* *coordinates* (plural)
    * represents a list of coordiantes
        * "X,Y" -> single coordinate
        * "0-7,Y" -> all X from 0 until 7 and given Y
        * "X,0-7" -> all Y from 0 until 7 and given X
        * "0-7,0-7" -> all X from 0 until 7 and all Y from 0 until 7
* *color*
    * "r,g" -> **r** is the red value of the color and must be between 0 and 3. **g** is the green value of the color and must be between 0 and 3.