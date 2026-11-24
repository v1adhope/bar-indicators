package app

import (
	"bytes"
	"context"
	"encoding/json/v2"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/shopspring/decimal"
	"github.com/v1adhope/bar-indicators/cryptocurrency/internal/coinmarketcap"
	"github.com/v1adhope/bar-indicators/cryptocurrency/internal/entity"
	"github.com/v1adhope/bar-indicators/lib/httpclient"
)

const _pricePlaces = 3

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

	client := httpclient.New()

	data, err := httpclient.FetchData[coinmarketcap.Data](ctx, client, u.String(), headers)
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
