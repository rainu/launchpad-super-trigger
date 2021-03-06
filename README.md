# launchpad-super-trigger
Trigger application for the Novation "Launchpad S" and "Launchpad MK2"

# Get the Binary

Before you can build the binary, you have to install the **rtmidi** library on your machine:
```bash
apt-get install librtmidi-dev
```

Now you can build it on your own (you will need [golang](https://golang.org/) installed):
```bash
go build -a -installsuffix cgo ./cmd/lst/
```

Then start the application:
```bash
./launchpad-super-trigger -config config.yaml
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
    0,1:
     actor: rest-test
     color:
      ready: 0,0
      progress: 1,1
      success: 2,2
      failed: 3,3
    1,0:
     actor: mqtt-test
```

## Config documentation

| config key | default | mandatory | description |
|---|---|---|---|
| general | - | false | Contains all general settings. |
| general.startPage | 0 | false | The page which should active when the application starts. |
| general.startBrightness | 100 | false | The brightness which should be used at start the application. It is only effective when the Launchpad MK2 with RGB-Colors is used. |
| general.navigationMode | 0 | false | The page navigation mode: 0 => binary mode; 1 => toggle mode |
| general.sensorStore | - | false | If the sensor data should be persisted, here you can set the file path to the [bolt-db](https://github.com/boltdb/bolt). Otherwise the sensor date will only keep in memory. |
| general.compressSensorData | false | false | Should the received sensor data be (gzip) compressed? This will work with all store types (memory or bolt). |
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
| actors.rest[*name*].bodyTemplate | - | false | The *body template* for the http request. |
| actors.mqtt | - | false | Contains all available mqtt actors. An mqtt actor will send a message to a given topic. |
| actors.mqtt[*name*].connection | - | false | The name of to mqtt connection which should be used for this mqtt actor. |
| actors.mqtt[*name*].topic | - | **true** | The topic name. |
| actors.mqtt[*name*].qos | 0 | false | The QualityOfService for the published message. |
| actors.mqtt[*name*].retained | false | false | Should the message be retained? |
| actors.mqtt[*name*].body | - | false | The body content for the mqtt message. |
| actors.mqtt[*name*].bodyBase64 | - | false | The body content for the mqtt message encoded in base64. |
| actors.mqtt[*name*].bodyPath | - | false | The body file for the mqtt message. |
| actors.mqtt[*name*].bodyTemplate | - | false | The *body template* for the mqtt message. |
| actors.combined | - | false | Contains all available combined actors. An combined actor will call other actors (sequential or in parallel). |
| actors.combined[*name*].actors | - | **true** | The list of underlying actor names. Must be greater or equal than 2! |
| actors.combined[*name*].parallel | false | false | How should the underlying actors be called. If true they will be called parallel. Otherwise the will be called sequential. |
| actors.conditional | - | false | Contains all available conditional actors. An conditional actor will call other actors depends on the given conditions. |
| actors.conditional[*name*].conditions | - | **true** | Contains the conditions. |
| actors.conditional[*name*].conditions[].actor | - | **true** | The actor which should be called if the condition is met. |
| actors.conditional[*name*].conditions[].datapoint | - | **true** | The reference for the underlying data point. |
| actors.conditional[*name*].conditions[].expression | - | **true** | The reference for the underlying data point. |
| actors.conditional[*name*].conditions[].expression.eq | - | false | This expression matches if the data point is *equal* than the given value. |
| actors.conditional[*name*].conditions[].expression.ne | - | false | This expression matches if the data point is *not equal* than the given value. |
| actors.conditional[*name*].conditions[].expression.lt | - | false | This expression matches if the data point is *less than* the given value. |
| actors.conditional[*name*].conditions[].expression.lte | - | false | This expression matches if the data point is *less or equal* than the given value. |
| actors.conditional[*name*].conditions[].expression.gt | - | false | This expression matches if the data point is *greater than* the given value. |
| actors.conditional[*name*].conditions[].expression.gte | - | false | This expression matches if the data point is *greater or equal than* the given value. |
| actors.conditional[*name*].conditions[].expression.match | - | false | This expression matches if the data point *matches* the given regular expression the given value. |
| actors.conditional[*name*].conditions[].expression.nmatch | - | false | This expression matches if the data point *not matches* the given regular expression the given value. |
| actors.conditional[*name*].conditions[].expression.contains | - | false | This expression matches if the data point *contains* the given value. |
| actors.conditional[*name*].conditions[].expression.ncontains | - | false | This expression matches if the data point *contains not* the given value. |
| actors.gfxBlink | - | false | Contains all available gfx blink actors. A blink actor will draw blinking pads. |
| actors.gfxBlink[*name*].on | - | true | The on *color*. |
| actors.gfxBlink[*name*].off | 0,0 | false | The off *color*. |
| actors.gfxBlink[*name*].interval | 1s | false | The interval of blink animation. |
| actors.gfxBlink[*name*].duration | until page leave | false | The duration of the blink interval. |
| actors.gfxWave | - | false | Contains all available gfx wave actors. A wave actor will draw waves on the pads. |
| actors.gfxWave[*name*].square | false | false | Should the waveform be square? |
| actors.gfxWave[*name*].color | 0,3 | false | The color of the wave. |
| actors.gfxWave[*name*].delay | 500ms | false | The delay between wave steps. |
| actors.metaSwitchPage | - | false | Contains all available page switch actors. A page switch actor can switch to a given page. |
| actors.metaSwitchPage[*name*].target | - | true | The page you want to switch to.  |
| actors.metaSwitchNavigationMode | - | false | Contains all available navigation mode switch actors. A navigation mode switch actor can switch to a given navigation mode. |
| actors.metaSwitchNavigationMode[*name*].mode | - | true | The page navigation mode: 0 => binary mode; 1 => toggle mode |
| actors.metaSwitchLock | - | false | Contains all available lock switch actors. A lock switch actor can switch to a given lock state (arm and unarmed). |
| actors.metaSwitchLock[*name*].lock | - | true | Should the launchpad be armed (true) or unarmed (false)? |
| connections | - | false | Contains all available connections, which can be used by different actors. |
| connections.mqtt | - | false | Contains all available MQTT connections. |
| connections.mqtt[*name*].broker | - | **true** | The url of the MQTT-Broker. |
| connections.mqtt[*name*].clientId | - | false | The client id which is send to the mqtt broker. |
| connections.mqtt[*name*].username | - | false | The username. |
| connections.mqtt[*name*].password | - | false | The password. |
| sensors | - | false | Contains all available sensors. A sensor will listen for (incoming) data. |
| sensors.command | - | false | Contains all available command actors. A command actor will runs a command on the local machine. |
| sensors.command[*name*].name | - | **true** | The name/path of the command to execute. |
| sensors.command[*name*].args | - | false | A list of arguments which will be send to the command. |
| sensors.command[*name*].interval | - | **true** | The interval. |
| sensors.command[*name*].data | - | false | Contains all data points. A data point is a part of the received message. |
| sensors.command[*name*].data.complete | - | false | The name of the data point for the complete message. |
| sensors.command[*name*].data.gjson | - | false | Contains all [gjson](https://github.com/tidwall/gjson) data points. |
| sensors.command[*name*].data.gjson[*name*] | - | false | The gjson path to use to extract data point. |
| sensors.command[*name*].data.split | - | false | Contains all split data points. |
| sensors.command[*name*].data.split[*name*].separator | - | true | The separator which should be used to split the whole data. |
| sensors.command[*name*].data.split[*name*].index | - | true | The index of the split element which should be used as data point. Must be greater or equal 0! |
| sensors.command[*name*].data.gojq | - | false | Contains all [gojq](https://github.com/itchyny/gojq) data points. |
| sensors.command[*name*].data.gojq[*name*] | - | true | The gojq query which should be used to extract the data point. |
| sensors.mqtt | - | false | Contains all available mqtt sensors. A mqtt sensor will listen for mqtt topics. |
| sensors.mqtt[*name*].connection | - | **true** | The name of to mqtt connection which should be used for this mqtt sensor. |
| sensors.mqtt[*name*].topic | - | **true** | The topic name. |
| sensors.mqtt[*name*].qos | 0 | false | The QualityOfService for the subscription message. |
| sensors.mqtt[*name*].data | - | false | Contains all data points. A data point is a part of the received message. |
| sensors.mqtt[*name*].data.complete | - | false | The name of the data point for the complete message. |
| sensors.mqtt[*name*].data.gjson | - | false | Contains all [gjson](https://github.com/tidwall/gjson) data points. |
| sensors.mqtt[*name*].data.gjson[*name*] | - | false | The gjson path to use to extract data point. |
| sensors.mqtt[*name*].data.split | - | false | Contains all split data points. |
| sensors.mqtt[*name*].data.split[*name*].separator | - | true | The separator which should be used to split the whole data. |
| sensors.mqtt[*name*].data.split[*name*].index | - | true | The index of the split element which should be used as data point. Must be greater or equal 0! |
| sensors.mqtt[*name*].data.gojq | - | false | Contains all [gojq](https://github.com/itchyny/gojq) data points. |
| sensors.mqtt[*name*].data.gojq[*name*] | - | true | The gojq query which should be used to extract the data point. |
| sensors.rest | - | false | Contains all available rest sensors. An rest sensor will call an rest service in an given interval. |
| sensors.rest[*name*].method | GET | false | The http method which should be used. |
| sensors.rest[*name*].url | - | **true** | The target url. |
| sensors.rest[*name*].header[*name*][]*value* | - | false | The http headers which should be used. |
| sensors.rest[*name*].body | - | false | The body content for the http request. |
| sensors.rest[*name*].bodyBase64 | - | false | The body content for the http request encoded in base64. |
| sensors.rest[*name*].bodyPath | - | false | The body file for the http request. |
| sensors.rest[*name*].interval | - | **true** | The interval. |
| sensors.rest[*name*].data | - | false | Contains all data points. A data point is a part of the received message. |
| sensors.rest[*name*].data.complete | - | false | The name of the data point for the complete message. |
| sensors.rest[*name*].data.gjson | - | false | Contains all [gjson](https://github.com/tidwall/gjson) data points. |
| sensors.rest[*name*].data.gjson[*name*] | - | false | The gjson path to use to extract data point. |
| sensors.rest[*name*].data.split | - | false | Contains all split data points. |
| sensors.rest[*name*].data.split[*name*].separator | - | true | The separator which should be used to split the whole data. |
| sensors.rest[*name*].data.split[*name*].index | - | true | The index of the split element which should be used as data point. Must be greater or equal 0! |
| sensors.rest[*name*].data.gojq | - | false | Contains all [gojq](https://github.com/itchyny/gojq) data points. |
| sensors.rest[*name*].data.gojq[*name*] | - | true | The gojq query which should be used to extract the data point. |
| layout | - | false | Contains all layout settings. |
| layout.pages | - | false | Contains all page settings. |
| layout.pages[*pageNumber*] | - | false | Contains a page setting. |
| layout.pages[*pageNumber*].trigger | - | false | Contains settings about the trigger on this page. |
| layout.pages[*pageNumber*].trigger[*coordinates*] | - | false | Contains settings about the trigger which should be called if the given button at *coordinate* is hit. |
| layout.pages[*pageNumber*].trigger[*coordinates*].actor | - | **true** | The name of the actor (see actor config) which should be called if the trigger is hit. |
| layout.pages[*pageNumber*].trigger[*coordinates*].color | - | false | Contains the *color* settings about the trigger. |
| layout.pages[*pageNumber*].trigger[*coordinates*].color.disabled | false | false | Should color be used? |
| layout.pages[*pageNumber*].trigger[*coordinates*].color.ready | - | false | The color which should be used if the trigger is ready. |
| layout.pages[*pageNumber*].trigger[*coordinates*].color.progress | - | false | The color which should be used as long as the actor is in progress. |
| layout.pages[*pageNumber*].trigger[*coordinates*].color.success | - | false | The color which should be used if the actor work was done successfully. |
| layout.pages[*pageNumber*].trigger[*coordinates*].color.failed | - | false | The color which should be used if the actor work was done wrong. |
| layout.pages[*pageNumber*].plotter | - | false | Contains settings about the plotters on this page. |
| layout.pages[*pageNumber*].plotter.progressbar | - | false | Contains all progressbar plotter for this page. |
| layout.pages[*pageNumber*].plotter.progressbar[].datapoint | - | **true** | The reference for the underlying data point. |
| layout.pages[*pageNumber*].plotter.progressbar[].min | 0 | false | The minimum value of the progressbar. |
| layout.pages[*pageNumber*].plotter.progressbar[].max | 100 | false | The maximum value of the progressbar. |
| layout.pages[*pageNumber*].plotter.progressbar[].vertical | false | false | Is the progressbar a vertical one? |
| layout.pages[*pageNumber*].plotter.progressbar[].quadrant | *whole pad* | false | In which mathematical quadrant should the progressbar be plotted. |
| layout.pages[*pageNumber*].plotter.progressbar[].x | - | false | If the progressbar *is vertical*, where should the progressbar be plotted. |
| layout.pages[*pageNumber*].plotter.progressbar[].y | - | false | If the progressbar *is horizontal*, where should the progressbar be plotted. |
| layout.pages[*pageNumber*].plotter.progressbar[].rtl | false | false | Should the progressbar be filled from right to left? |
| layout.pages[*pageNumber*].plotter.progressbar[].fill | 0,0 | false | The filled **color** which should be used. |
| layout.pages[*pageNumber*].plotter.progressbar[].empty | 0,3 | false | The empty **color** which should be used. |
| layout.pages[*pageNumber*].plotter.static | - | false | Contains all static plotter for this page. A static plotter will color on button if an expression on the data point matches. |
| layout.pages[*pageNumber*].plotter.static[].datapoint | - | **true** | The reference for the underlying data point. |
| layout.pages[*pageNumber*].plotter.static[].pos | - | **true** | The *coordinate* where should plot. |
| layout.pages[*pageNumber*].plotter.static[].defaultColor | - | false | The color which should be used if no expression matches. |
| layout.pages[*pageNumber*].plotter.static[].expressions.eq | - | false | This expression matches if the data point is *equal* than the given value. |
| layout.pages[*pageNumber*].plotter.static[].expressions.eq[].value | - | false | This expression value. |
| layout.pages[*pageNumber*].plotter.static[].expressions.eq[].color | - | false | The *color* which should be used if the expression matches. |
| layout.pages[*pageNumber*].plotter.static[].expressions.ne | - | false | This expression matches if the data point is *not equal* than the given value. |
| layout.pages[*pageNumber*].plotter.static[].expressions.ne[].value | - | false | This expression value. |
| layout.pages[*pageNumber*].plotter.static[].expressions.ne[].color | - | false | The *color* which should be used if the expression matches. |
| layout.pages[*pageNumber*].plotter.static[].expressions.lt | - | false | This expression matches if the data point is *less than* the given value. |
| layout.pages[*pageNumber*].plotter.static[].expressions.lt[].value | - | false | This expression value. |
| layout.pages[*pageNumber*].plotter.static[].expressions.lt[].color | - | false | The *color* which should be used if the expression matches. |
| layout.pages[*pageNumber*].plotter.static[].expressions.lte | - | false | This expression matches if the data point is *less or equal* than the given value. |
| layout.pages[*pageNumber*].plotter.static[].expressions.lte[].value | - | false | This expression value. |
| layout.pages[*pageNumber*].plotter.static[].expressions.lte[].color | - | false | The *color* which should be used if the expression matches. |
| layout.pages[*pageNumber*].plotter.static[].expressions.gt | - | false | This expression matches if the data point is *greater than* the given value. |
| layout.pages[*pageNumber*].plotter.static[].expressions.gt[].value | - | false | This expression value. |
| layout.pages[*pageNumber*].plotter.static[].expressions.gt[].color | - | false | The *color* which should be used if the expression matches. |
| layout.pages[*pageNumber*].plotter.static[].expressions.gte | - | false | This expression matches if the data point is *greater or equal than* the given value. |
| layout.pages[*pageNumber*].plotter.static[].expressions.gte[].value | - | false | This expression value. |
| layout.pages[*pageNumber*].plotter.static[].expressions.gte[].color | - | false | The *color* which should be used if the expression matches. |
| layout.pages[*pageNumber*].plotter.static[].expressions.match | - | false | This expression matches if the data point *matches* the given regular expression the given value. |
| layout.pages[*pageNumber*].plotter.static[].expressions.match[].value | - | false | This expression value. |
| layout.pages[*pageNumber*].plotter.static[].expressions.match[].color | - | false | The *color* which should be used if the expression matches. |
| layout.pages[*pageNumber*].plotter.static[].expressions.nmatch | - | false | This expression matches if the data point *not matches* the given regular expression the given value. |
| layout.pages[*pageNumber*].plotter.static[].expressions.nmatch[].value | - | false | This expression value. |
| layout.pages[*pageNumber*].plotter.static[].expressions.nmatch[].color | - | false | The *color* which should be used if the expression matches. |
| layout.pages[*pageNumber*].plotter.static[].expressions.contains | - | false | This expression matches if the data point *contains* the given value. |
| layout.pages[*pageNumber*].plotter.static[].expressions.contains[].value | - | false | This expression value. |
| layout.pages[*pageNumber*].plotter.static[].expressions.contains[].color | - | false | The *color* which should be used if the expression matches. |
| layout.pages[*pageNumber*].plotter.static[].expressions.ncontains | - | false | This expression matches if the data point *contains not* the given value. |
| layout.pages[*pageNumber*].plotter.static[].expressions.ncontains[].value | - | false | This expression value. |
| layout.pages[*pageNumber*].plotter.static[].expressions.ncontains[].color | - | false | The *color* which should be used if the expression matches. |

* *pageNumber*
    * The page number must be a number from **0** until **255**.
    * or it can be a string representing the top page line, where each active button is a "1" nad the others "0"
        * for example: 10000001    (the first an last page button must be active for this page)
* *coordinate*
    * "X,Y" -> X and Y must have a value of 0-7.
* *coordinates* (plural)
    * represents a list of coordiantes
        * "X,Y" -> single coordinate
        * "0-7,Y" -> all X from 0 until 7 and given Y
        * "X,0-7" -> all Y from 0 until 7 and given X
        * "0-7,0-7" -> all X from 0 until 7 and all Y from 0 until 7
* *color*
    * "r,g" (Launchpad S) -> **r** is the red value of the color and must be between 0 and 3. **g** is the green value of the color and must be between 0 and 3.
    * "r,g,b" (Launchpad MK2) -> **r** is the red value of the color and must be between 0 and 63. **g** is the green value of the color and must be between 0 and 63. **b** is the blue value of the color and must be between 0 and 63.
* *datapoint reference*
    * **sensorName**.**datapointName**
* *body template*
    * You can define a [go template](https://golang.org/pkg/text/template/) as body content. In this template there are helper functions available:
        * **.DataPoint "<sensor.datapoint>"**
        * **.DataPointOr "<sensor.datapoint>" "<defaultValue>"**
    * Example - print the data point **available** from sensor **notebook**:
        ```
        {{.DataPoint "notebook.available"}}  
        ```
    * Example - print the data point **available** from sensor **notebook** or "offline" if no data is available:
        ```
        {{.DataPointOr "notebook.available" "offline"}}  
        ```
    * The most public functions from go default lib are available: 
        * [strings](https://golang.org/pkg/strings/)
            * prefix: **strings_**
            * example: [ **strings_Count** ](https://golang.org/pkg/strings/#Count)
        * [math](https://golang.org/pkg/math/)
            * prefix: **math_**
            * example: [ **math_Abs** ](https://golang.org/pkg/math/#Abs)
        * [strconv](https://golang.org/pkg/strconv/)
            * prefix: **strconv_**
            * example: [ **strconv_ParseBool** ](https://golang.org/pkg/strconv/#ParseBool)
        * [rand](https://golang.org/pkg/rand/)
            * prefix: **rand_**
            * example: [ **rand_Int63** ](https://golang.org/pkg/rand/#Int63)
        * [time](https://golang.org/pkg/time/)
            * prefix: **time_**
            * example: [ **time_Now** ](https://golang.org/pkg/time/#Now)
        * [fmt_Sprintf](https://golang.org/pkg/fmt/#Sprintf)
    * Furthermore there exists functions for basic arithmetic:
        * **add** - addition
        * **sub** - subtraction
        * **mul** - multiplication
        * **div** - diversion
    * Complex example - **add 13* to data point **ram** from sensor **notebook**:
        ```
        {{ $sr := .DataPointOr "notebook.ram" "0" }}
        {{ $r := strconv_Atoi $sr }}
        {{ add $r 13 }}
        ```
 
### Special sensors

For each connection there is a "hidden" sensor which holds the connection status. This sensor can be used such like other
sensors. You can reference them and use there datapoint(s). The name of the sensor is **__connection[&lt;connectionName&gt;]**.
The name of the datapoint which holds the connection information is **status**. And it can have two states:
* **Connected** -> the corresponding connection is established
* **Disconnected** -> the corresponding connection is lost
 
Example:

```yaml
connections:
 mqtt:
  test-broker:
   broker: tcp://broker-host:1883
   clientId: launchpad
layout:
 pages:
  0:
    plotter:
      static:
        - datapoint: __connection[test-broker].status
          pos: 0,0
            expressions:
              eq:
                - value: Connected
                  color: 0,3
                - value: Disconnected
                  color: 3,0
```

You will see a green light on the left top corner if the mqtt-connection "test-broker" is established, otherwise there is a red light!

# Special keys

In the following you can read a description what do the special keys. The special keys are the round keys on the right side. 

## vol

With this key you can change the brightness. This feature is only available if you use a Launchpad MK2 with RGB-Colors.

## pan

With this key you can switch between the different page navigation modes:

|mode|color|description|
|---|---|---|
| Binary mode (default) | off | In this mode only the pressed page key is switched on/off. |
| Toggle mode | green | In this mode all page keys are switched off except the pressed page key. |

## arm

With this key you can switch between lock modes. 

|mode|color|description|
|---|---|---|
| Unlocked (default) | off | In this mode you can touch the pad keys and their underlying action will performed. |
| Locked | red | In this mode you can only touch the special or page keys. If you touch the pad keys, no action will performed. |