package httpclient

import (
	"context"
	"encoding/json/v2"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/hashicorp/go-retryablehttp"
)

const (
	_reqAttempts = 3
	_reqTimeout  = 10 * time.Second
)

func New() *retryablehttp.Client {
	logger := slog.New(slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{
		Level: slog.LevelWarn,
	}))
	client := retryablehttp.NewClient()
	client.Logger = logger
	client.ErrorHandler = retryablehttp.PassthroughErrorHandler
	client.RetryMax = _reqAttempts
	client.HTTPClient.Timeout = _reqTimeout
	return client
}

func FetchData[T any](ctx context.Context, client *retryablehttp.Client, url string, headers http.Header) (T, error) {
	var data T
	req, err := retryablehttp.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return data, fmt.Errorf("request: %w", err)
	}

	req.Header = headers.Clone()

	// A 200 OK response can still time out while reading its body.
	// Decode inside the retry loop so body read errors trigger another attempt.
	decodeResponse := func(resp *http.Response) error {
		if resp.StatusCode != http.StatusOK {
			return nil
		}
		var next T
		if err := json.UnmarshalRead(resp.Body, &next); err != nil {
			return fmt.Errorf("decode: %w", err)
		}
		data = next
		return nil
	}
	req.SetResponseHandler(decodeResponse)

	resp, err := client.Do(req)
	if resp != nil {
		defer resp.Body.Close()
	}
	if resp != nil && resp.StatusCode != http.StatusOK {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return data, fmt.Errorf("get: unexpected status: %s; read body: %w", resp.Status, err)
		}
		return data, fmt.Errorf("get: unexpected status: %s; body: %s", resp.Status, body)
	}
	if err != nil {
		return data, fmt.Errorf("get: %w", err)
	}
	return data, nil
}
