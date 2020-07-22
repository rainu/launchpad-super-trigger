package config

import (
	"encoding/base64"
	"github.com/rainu/launchpad-super-trigger/actor"
	"github.com/rainu/launchpad-super-trigger/config"
	"github.com/rainu/launchpad-super-trigger/pad"
	"io"
	"net/http"
	"os"
	"strings"
)

func ConfigureDispatcher(configReader io.Reader) (*pad.TriggerDispatcher, error) {
	parsedConfig, err := config.ReadConfig(configReader)
	if err != nil {
		return nil, err
	}

	dispatcher := &pad.TriggerDispatcher{}
	restHandler := map[string]actor.Actor{}

	for actorName, restActor := range parsedConfig.Actors.Rest {
		handler := &actor.RestActionHandler{
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

		restHandler[actorName] = handler
	}

	for pageNumber, page := range parsedConfig.Layout.Pages {
		handler := &pageHandler{
			pageNumber: pad.PageNumber(pageNumber),
			page:       page,
		}
		handler.Init(restHandler)
		dispatcher.AddPageHandler(handler, handler.pageNumber)
	}

	return dispatcher, nil
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
