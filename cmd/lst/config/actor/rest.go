package actor

import (
	"bytes"
	"encoding/base64"
	"github.com/rainu/launchpad-super-trigger/actor"
	"github.com/rainu/launchpad-super-trigger/config"
	"github.com/rainu/launchpad-super-trigger/template"
	"go.uber.org/zap"
	"io"
	"net/http"
	"os"
	"strings"
)

func buildRest(actors map[string]actor.Actor, restActors map[string]config.RestActor, engine *template.Engine) {
	for actorName, restActor := range restActors {
		handler := &actor.Rest{
			HttpClient: http.DefaultClient,
			Method:     restActor.Method,
			Url:        restActor.URL,
			Header:     restActor.Header,
		}

		if restActor.BodyRaw != "" {
			handler.Body = rawBody(restActor.BodyRaw)
		} else if restActor.BodyB64 != "" {
			handler.Body = b64Body(restActor.BodyB64)
		} else if restActor.BodyPath != "" {
			handler.Body = fileBody(restActor.BodyPath)
		} else if restActor.BodyTemplate != "" {
			if err := engine.RegisterTemplate(actorName, restActor.BodyTemplate); err != nil {
				zap.L().Fatal("Failed to parse template!", zap.Error(err))
			}

			handler.Body = templateBody(actorName, engine)
		}

		actors[actorName] = handler
	}
}

func rawBody(body string) func() io.Reader {
	return func() io.Reader {
		return strings.NewReader(body)
	}
}

func b64Body(b64Body string) func() io.Reader {
	return func() io.Reader {
		return base64.NewDecoder(base64.StdEncoding, strings.NewReader(b64Body))
	}
}

func fileBody(bodyPath string) func() io.Reader {
	return func() io.Reader {
		file, err := os.Open(bodyPath)
		if err != nil {
			panic(err)
		}

		return file
	}
}

func templateBody(actorName string, engine *template.Engine) func() io.Reader {
	return func() io.Reader {
		value, err := engine.Execute(actorName)
		if err != nil {
			zap.L().Fatal("Failed to execute template!", zap.Error(err))
		}

		return bytes.NewReader(value)
	}
}
