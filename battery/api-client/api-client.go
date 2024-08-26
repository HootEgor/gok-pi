package apiclient

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"gok-pi/battery/entity"
	"gok-pi/internal/lib/sl"
	"io"
	"log/slog"
	"net/http"
	"strings"
	"time"
)

const (
	maxRetry  = 5
	retryStep = 3
)

var httpClient = &http.Client{}

type ApiClient struct {
	url   string
	token string
	log   *slog.Logger
}

func New(url, token string, log *slog.Logger) *ApiClient {
	log.With(
		slog.String("url", url),
		sl.Secret("token", token),
	).Info("creating api client")
	return &ApiClient{
		url:   url,
		token: token,
		log:   log.With(sl.Module("battery.client")),
	}
}

func (c *ApiClient) Status() (*entity.SystemStatus, error) {
	body, err := c.requestWithRetry(http.MethodGet, nil, c.url, "status")
	if err != nil {
		return nil, err
	}
	status, err := entity.ParseSystemStatus(body)
	if err != nil {
		return nil, fmt.Errorf("parsing status: %w", err)
	}
	return status, nil
}

func (c *ApiClient) StartDischarge(power int) error {
	_, err := c.requestWithRetry(http.MethodPost, nil, c.url, "setpoint", "discharge", fmt.Sprintf("%d", power))
	return err
}

func (c *ApiClient) StopDischarge() error {
	_, err := c.requestWithRetry(http.MethodPost, nil, c.url, "setpoint", "discharge", "0")
	return err
}

func (c *ApiClient) fullPath(params ...string) string {
	return strings.Join(params, "/")
}

func (c *ApiClient) requestWithRetry(method string, data interface{}, params ...string) ([]byte, error) {
	path := c.fullPath(params...)
	log := c.log.With(
		slog.String("url", path),
		slog.String("method", method),
	)
	var body []byte
	if data != nil {
		var err error
		body, err = json.Marshal(data)
		if err != nil {
			log.Error("marshalling body", sl.Err(err))
			return nil, fmt.Errorf("marshalling body: %w", err)
		}
	}

	for i := 0; i < maxRetry; i++ {
		responseBody, err := c.doRequest(method, path, body)
		if err == nil {
			return responseBody, nil
		}
		log.With(
			slog.Int("attempt", i+1),
		).Debug("retrying request")
		time.Sleep(time.Duration((i+1)*retryStep) * time.Second)
	}
	return nil, fmt.Errorf("request failed after %d retries", maxRetry)
}

func (c *ApiClient) doRequest(method, url string, body []byte) ([]byte, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var err error
	log := c.log.With(
		slog.String("url", url),
		sl.Secret("token", c.token),
		slog.String("method", method),
	)
	t1 := time.Now()
	defer func() {
		log = log.With(slog.Float64("duration", time.Since(t1).Seconds()))
		if err != nil {
			log.Error("api request", sl.Err(err))
		} else {
			log.Debug("api request")
		}
	}()

	req, err := http.NewRequestWithContext(ctx, method, url, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Auth-Token", c.token)

	resp, err := httpClient.Do(req)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			err = fmt.Errorf("request timeout")
		}
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)

	log = log.With(slog.Int("status", resp.StatusCode))
	if resp.StatusCode >= 400 {
		err = fmt.Errorf("received status code: %d", resp.StatusCode)
		return nil, err
	}

	body, err = io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}
