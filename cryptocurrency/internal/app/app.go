package app

import (
	"bytes"
	"context"
	"encoding/json/v2"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/hashicorp/go-retryablehttp"
	"github.com/shopspring/decimal"
	"github.com/v1adhope/bar-indicators/cryptocurrency/internal/coinmarketcap"
	"github.com/v1adhope/bar-indicators/cryptocurrency/internal/entity"
)

const (
	_reqAttempts = 3
	_reqTimeout  = 10 * time.Second

	_pricePlaces = 3
)

// Run fetches and prints cryptocurrency quotes in Waybar JSON format.
func Run(ctx context.Context, apiToken, convertId, rawCoinIds string) error {
	res, err := getCryptocurrency(ctx, apiToken, convertId, rawCoinIds)
	if err != nil {
		return err
	}

	fmt.Println(res)
	return nil
}

// getCryptocurrency returns the Waybar cryptocurrency JSON for the requested coins.
func getCryptocurrency(ctx context.Context, apiToken, convertId, rawCoinIds string) (string, error) {
	u, err := url.Parse("https://pro-api.coinmarketcap.com/v2/cryptocurrency/quotes/latest")
	if err != nil {
		return "", fmt.Errorf("parse: %w", err)
	}
	query := &url.Values{}
	query.Add("id", rawCoinIds)
	query.Add("convert_id", convertId)
	u.RawQuery = query.Encode()

	headers := http.Header{}
	headers.Set("Accept", "application/json")
	headers.Set("X-CMC_PRO_API_KEY", apiToken)

	logger := slog.New(slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{
		Level: slog.LevelWarn,
	}))
	retryClient := retryablehttp.NewClient()
	retryClient.Logger = logger
	retryClient.ErrorHandler = retryablehttp.PassthroughErrorHandler
	retryClient.RetryMax = _reqAttempts
	retryClient.HTTPClient.Timeout = _reqTimeout

	data, err := fetchCryptocurrencyData(ctx, retryClient, u.String(), headers)
	if err != nil {
		return "", err
	}

	module := entity.Module{}
	coinIds := strings.Split(rawCoinIds, ",")
	for _, coinId := range coinIds {
		if !data.Data[coinId].Quote[convertId].Price.IsPositive() {
			return "", fmt.Errorf("price unavailable for coin %s in currency %s", coinId, convertId)
		}
	}

	coinIdFirst := coinIds[0]
	price := data.Data[coinIdFirst].Quote[convertId].Price
	if data.Data[coinIdFirst].IsFiat == 0 {
		module.Text = price.StringFixedBank(_pricePlaces)
	} else {
		module.Text = decimal.NewFromInt(1).Div(price).StringFixedBank(_pricePlaces)
	}

	buf := bytes.Buffer{}

	w := tabwriter.NewWriter(&buf, 0, 0, 1, ' ', tabwriter.AlignRight)

	fmt.Fprintf(w, "%-6s\t%s\t%s\t\n", "Ticker", "Price", "24h%")

	for _, coinId := range coinIds[1:] {
		fmt.Fprintf(w,
			"%-6s\t%s\t%0.2f\t\n",
			data.Data[coinId].Symbol,
			data.Data[coinId].Quote[convertId].Price.StringFixedBank(_pricePlaces),
			data.Data[coinId].Quote[convertId].PercentChange24h,
		)
	}
	if err := w.Flush(); err != nil {
		return "", fmt.Errorf("flush tooltip: %w", err)
	}

	t, err := time.Parse(time.RFC3339Nano, data.Status.Timestamp)
	if err != nil {
		return "", fmt.Errorf("parse timestamp: %w", err)
	}
	fmt.Fprint(&buf, "\n <b>Timestamp</b>\n "+t.Local().Format(time.DateTime))

	module.Tooltip = buf.String()
	module.Tooltip = strings.Replace(module.Tooltip, "Ticker", "<b>Ticker</b>", 1)
	module.Tooltip = strings.Replace(module.Tooltip, "Price", "<b>Price</b>", 1)
	module.Tooltip = strings.Replace(module.Tooltip, "24h%", "<b>24h%</b>", 1)

	json, err := json.Marshal(module)
	if err != nil {
		return "", fmt.Errorf("marshal: %w", err)
	}

	return string(json), nil
}

func fetchCryptocurrencyData(ctx context.Context, client *retryablehttp.Client, url string, headers http.Header) (coinmarketcap.Data, error) {
	data := coinmarketcap.Data{}
	req, err := retryablehttp.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return data, fmt.Errorf("request: %w", err)
	}

	req.Header = headers.Clone()

	// Decode inside the retry loop so body read errors trigger another attempt.
	decodeResponse := func(resp *http.Response) error {
		if resp.StatusCode != http.StatusOK {
			return nil
		}
		next := coinmarketcap.Data{}
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
