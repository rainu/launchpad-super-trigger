package sensor

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	"io"
	"io/ioutil"
	"net/http"
	"sync"
	"time"
)

type Rest struct {
	callbackHandler

	HttpClient *http.Client
	Method     string
	Url        string
	Header     map[string][]string
	Body       func() io.Reader
	Interval   time.Duration

	running     bool
	mux         sync.RWMutex
	lastMessage []byte
}

func (r *Rest) Run(ctx context.Context) error {
	if r.running {
		return fmt.Errorf("listerner is already running")
	}

	ticker := time.NewTicker(r.Interval)
	defer ticker.Stop()

	r.running = true

	//first call
	if err := r.call(ctx); err != nil {
		zap.L().Error("Error while call rest sensor!", zap.Error(err))
	}

	for {
		select {
		case <-ticker.C:
			if err := r.call(ctx); err != nil {
				zap.L().Error("Error while call rest sensor!", zap.Error(err))
			}

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

	r.mux.Lock()
	r.lastMessage = rawBody
	r.mux.Unlock()

	r.callbackHandler.Call(r)
	return nil
}

func (r *Rest) LastMessage() []byte {
	r.mux.RLock()
	defer r.mux.RUnlock()

	return r.lastMessage
}
