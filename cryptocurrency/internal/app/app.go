package app

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/v1adhope/bar-indicators/cryptocurrency/internal/coinmarketcap"
	"github.com/v1adhope/bar-indicators/cryptocurrency/internal/entity"
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
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://pro-api.coinmarketcap.com/v2/cryptocurrency/quotes/latest", nil)
	if err != nil {
		return "", fmt.Errorf("request: %w", err)
	}

	query := &url.Values{}
	query.Add("id", rawCoinIds)
	query.Add("convert_id", convertId)

	req.Header.Set("Accept", "application/json")
	req.Header.Add("X-CMC_PRO_API_KEY", apiToken)
	req.URL.RawQuery = query.Encode()

	data := coinmarketcap.Data{}
	client := &http.Client{}
	resp, err := client.Do(req)
	if resp != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		return "", fmt.Errorf("get: %w", err)
	}

	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return "", fmt.Errorf("decode: %w", err)
	}

	module := entity.Module{}
	coinIds := strings.Split(rawCoinIds, ",")

	coinIdFirst := coinIds[0]
	if data.Data[coinIdFirst].IsFiat == 0 {
		module.Text = fmt.Sprintf("%0.3f", data.Data[coinIdFirst].Quote[convertId].Price)
	} else {
		module.Text = fmt.Sprintf("%0.3f", 1/data.Data[coinIdFirst].Quote[convertId].Price)
	}

	buf := bytes.Buffer{}

	func() {
		w := tabwriter.NewWriter(&buf, 0, 0, 1, ' ', tabwriter.AlignRight)
		defer w.Flush()

		fmt.Fprintf(w, "%-6s\t%s\t%s\t\n", "Ticker", "Price", "24h%")

		for _, coinId := range coinIds[1:] {
			fmt.Fprintf(w,
				"%-6s\t%0.3f\t%0.2f\t\n",
				data.Data[coinId].Symbol,
				data.Data[coinId].Quote[convertId].Price,
				data.Data[coinId].Quote[convertId].PercentChange24h,
			)
		}
	}()

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
