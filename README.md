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

# ConfigFile

```yaml
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
      failed: 3,3
```

## Config documentation

| config key | default | mandatory | description |
|---|---|---|---|
| actors | - | false | Contains all available actors. |
| actors.rest | - | false | Contains all available rest actors. An rest actor will call an rest service. |
| actors.rest[*name*].method | GET | false | The http method which should be used. |
| actors.rest[*name*].url | - | **true** | The target url. |
| actors.rest[*name*].header[*name*][]*value* | - | false | The http headers which should be used. |
| actors.rest[*name*].body | - | false | The body content for the http request. |
| actors.rest[*name*].bodyBase64 | - | false | The body content for the http request encoded in base64. |
| actors.rest[*name*].bodyPath | - | false | The body file for the http request. |
| layout | - | false | Contains all layout settings. |
| layout.pages | - | false | Contains all page settings. |
| layout.pages[*pageNumber*] | - | false | Contains a page setting. |
| layout.pages[*pageNumber*].trigger | - | false | Contains settings about the trigger on this page. |
| layout.pages[*pageNumber*].trigger[*coordinate*] | - | false | Contains settings about the trigger which should be called if the given button at *coordinate* is hit. |
| layout.pages[*pageNumber*].trigger[*coordinate*].actor | - | **true** | The name of the actor (see actor config) which should be called if the trigger is hit. |
| layout.pages[*pageNumber*].trigger[*coordinate*].color | - | false | Contains the *color* settings about the trigger. |
| layout.pages[*pageNumber*].trigger[*coordinate*].color.ready | - | false | The color which should be used if the trigger is ready. |
| layout.pages[*pageNumber*].trigger[*coordinate*].color.progress | - | false | The color which should be used as long as the actor is in progress. |
| layout.pages[*pageNumber*].trigger[*coordinate*].color.success | - | false | The color which should be used if the actor work was done successfully. |
| layout.pages[*pageNumber*].trigger[*coordinate*].color.failed | - | false | The color which should be used if the actor work was done wrong. |

* *pageNumber*
    * The page number must be a number from **0** until **255**.
* *coordinate*
    * "X,Y" -> X and Y must have a value of 0-7.
* *color*
    * "r,g" -> **r** is the red value of the color and must be between 0 and 3. **g** is the green value of the color and must be between 0 and 3.