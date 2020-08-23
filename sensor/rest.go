package sensor

import (
	"context"
	"fmt"
	"github.com/rainu/launchpad-super-trigger/sensor/store"
	"go.uber.org/zap"
	"io"
	"io/ioutil"
	"net/http"
	"time"
)

type Rest struct {
	callbackHandler

	HttpClient   *http.Client
	Method       string
	Url          string
	Header       map[string][]string
	Body         func() io.Reader
	Interval     time.Duration
	MessageStore store.Store

	running bool
}

func (r *Rest) Run(ctx context.Context) error {
	if r.running {
		return fmt.Errorf("listerner is already running")
	}
	defer func() {
		r.running = false
	}()

	timer := time.NewTimer(r.Interval)

	r.running = true

	//first call
	if err := r.call(ctx); err != nil {
		zap.L().Error("Error while call rest sensor!", zap.Error(err))
	}

	for {
		select {
		case <-timer.C:
			if err := r.call(ctx); err != nil {
				zap.L().Error("Error while call rest sensor!", zap.Error(err))
			}

			//reset timer
			timer = time.NewTimer(r.Interval)

		//wait until context closed
		case <-ctx.Done():
			return nil
		}
	}
}

func (r *Rest) call(ctx context.Context) error {
	var body io.Reader
	if r.Body != nil {
		body = r.Body()
	}

	req, err := http.NewRequest(r.Method, r.Url, body)
	if err != nil {
		return fmt.Errorf("could not create new request: %w", err)
	}
	if r.Header != nil {
		for header, values := range r.Header {
			for _, value := range values {
				req.Header.Set(header, value)
			}
		}
	}
	req = req.WithContext(ctx)

	resp, err := r.HttpClient.Do(req)
	if err != nil {
		return fmt.Errorf("error while do request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 200 && resp.StatusCode <= 299 {
		zap.L().Info("Rest call was successful!")
		return nil
	}

	rawBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if err := r.MessageStore.Set(rawBody); err != nil {
		zap.L().Error("Could not save message into message store!", zap.Error(err))
	}

	r.callbackHandler.Call(r)
	return nil
}

func (r *Rest) LastMessage() []byte {
	data, err := r.MessageStore.Get()
	if err != nil {
		zap.L().Error("Could not get message from message store!", zap.Error(err))
		return nil
	}

	return data
}
