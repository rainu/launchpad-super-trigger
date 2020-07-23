package actor

import (
	"fmt"
	"go.uber.org/zap"
	"io"
	"net/http"
)

type RestActionHandler struct {
	HttpClient *http.Client
	Method     string
	Url        string
	Header     map[string][]string
	Body       func() io.Reader
}

func (r *RestActionHandler) Do(ctx Context) error {
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
	req = req.WithContext(ctx.Context)

	resp, err := r.HttpClient.Do(req)
	if err != nil {
		return fmt.Errorf("error while do request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 200 && resp.StatusCode <= 299 {
		zap.L().Info("Rest call was successful!")
		return nil
	}

	return fmt.Errorf("bad status code: %d", resp.StatusCode)
}
