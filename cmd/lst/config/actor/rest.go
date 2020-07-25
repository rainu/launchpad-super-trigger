package actor

import (
	"encoding/base64"
	"github.com/rainu/launchpad-super-trigger/actor"
	"github.com/rainu/launchpad-super-trigger/config"
	"io"
	"net/http"
	"os"
	"strings"
)

func buildRest(actors map[string]actor.Actor, restActors map[string]config.RestActor) {
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
