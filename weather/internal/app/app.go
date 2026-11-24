package app

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json/v2"
	"errors"
	"fmt"
	"io"
	"log"
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/hashicorp/go-retryablehttp"
	"github.com/v1adhope/bar-indicators/weather/internal/database"
	"github.com/v1adhope/bar-indicators/weather/internal/entity"
	"github.com/v1adhope/bar-indicators/weather/internal/wttr"
)

const (
	_reqAttempts = 3
	_reqTimeout  = 10 * time.Second

	_timeModel12h = "03:04 PM"
	_timeModel24h = "15:04"
)

// Run fetches and prints weather for city, using saved results if the fetch fails.
func Run(ctx context.Context, city string) error {
	res, fetchErr := getWeather(ctx, city)
	if fetchErr == nil {
		fmt.Println(res)
	}

	db, err := database.Open(ctx)
	if err != nil {
		return errors.Join(fetchErr, err)
	}

	defer db.Close()

	if err := db.Init(ctx); err != nil {
		return errors.Join(fetchErr, err)
	}

	status := entity.StatusSuccess

	if fetchErr != nil {
		status, res = entity.StatusError, fetchErr.Error()

		fallback, err := db.GetLatestSuccess(ctx, city)
		if err != nil {
			recordErr := db.Record(ctx, city, status, res)

			if errors.Is(err, sql.ErrNoRows) {
				return errors.Join(fetchErr, recordErr)
			}

			return errors.Join(fetchErr, fmt.Errorf("load fallback: %w", err), recordErr)
		}

		log.Printf("%v; using last successful weather result", fetchErr)
		fmt.Println(strings.Replace(fallback, "Fallback not used", "Fallback used", 1))
	}

	return db.Record(ctx, city, status, res)
}

// getWeather returns the Waybar weather JSON for city.
//
// Assume wttr.in always provides current conditions, three forecast days,
// astronomy, and weather descriptions. If that contract changes, the slice
// accesses below can panic before history recording and fallback.
func getWeather(ctx context.Context, city string) (string, error) {
	u, err := url.Parse("https://wttr.in/" + city + "?format=j1")
	if err != nil {
		return "", fmt.Errorf("parse: %w", err)
	}
	logger := slog.New(slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{
		Level: slog.LevelWarn,
	}))
	retryClient := retryablehttp.NewClient()
	retryClient.Logger = logger
	retryClient.ErrorHandler = retryablehttp.PassthroughErrorHandler
	retryClient.RetryMax = _reqAttempts
	retryClient.HTTPClient.Timeout = _reqTimeout

	data, err := fetchWeatherData(ctx, retryClient, u.String())
	if err != nil {
		return "", err
	}

	module := entity.Module{}
	module.Text = fmt.Sprintf(
		"%s°(%s°)",
		data.CurrentCondition[0].TempC,
		data.CurrentCondition[0].FeelsLikeC,
	)

	buf := bytes.Buffer{}
	for _, row := range []struct{ Title, Value string }{
		{"Current temp:", data.CurrentCondition[0].TempC + "°"},
		{"Feels like:", data.CurrentCondition[0].FeelsLikeC + "°"},
		{"Humidity:", data.CurrentCondition[0].Humidity + "%"},
		{"Pressure:", data.CurrentCondition[0].Pressure + "hPa"},
		{"Wind speed:", data.CurrentCondition[0].WindspeedKmph + "km/h"},
		{"Description:", data.CurrentCondition[0].WeatherDesc[0].Value},
	} {
		fmt.Fprintln(&buf, row.Title+" "+row.Value)
	}

	sunriseTime, err := timeConvertFrom12to24H(data.Weather[0].Astronomy[0].SunRise)
	if err != nil {
		return "", fmt.Errorf("sunrise convert: %w", err)
	}
	sunsetTime, err := timeConvertFrom12to24H(data.Weather[0].Astronomy[0].SunSet)
	if err != nil {
		return "", fmt.Errorf("sunset convert: %w", err)
	}
	fmt.Fprint(&buf, "\n<b>Solar cycle</b>\n")
	fmt.Fprintln(&buf, "Sunrise at "+sunriseTime)
	fmt.Fprintln(&buf, "Sunset at "+sunsetTime)

	days, nowHour := 3, time.Now().Hour()
	for day := range days {
		switch day {
		case 0:
			fmt.Fprintf(&buf, "\n<b>Today</b>\n")
		case 1:
			fmt.Fprintf(&buf, "\n<b>Tomorrow</b>\n")
		case 2:
			fmt.Fprintf(&buf, "\n<b>After a day</b>\n")
		}
		for idx, weather := range data.Weather[day].Hourly {
			tWindow := idx * 3
			// Skip if time is overdue
			if day == 0 && nowHour > tWindow+2 {
				continue
			}
			if tWindow < 10 {
				fmt.Fprintf(&buf, "At 0%d:00 %2s°(%2s°) %s\n", tWindow, weather.TempC, weather.FeelsLikeC, checkDescription(weather.WeatherDesc[0].Value))
			} else {
				fmt.Fprintf(&buf, "At %d:00 %2s°(%2s°) %s\n", tWindow, weather.TempC, weather.FeelsLikeC, checkDescription(weather.WeatherDesc[0].Value))
			}
		}
	}

	fmt.Fprint(&buf, "\n<b>Timestamp</b>\n")
	fmt.Fprint(&buf, time.Now().Format(time.DateTime))
	fmt.Fprint(&buf, "\nFallback not used")

	module.Tooltip = buf.String()

	json, err := json.Marshal(module)
	if err != nil {
		return "", fmt.Errorf("marshal: %w", err)
	}
	return string(json), nil
}

func fetchWeatherData(ctx context.Context, client *retryablehttp.Client, url string) (wttr.Data, error) {
	data := wttr.Data{}
	req, err := retryablehttp.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return data, fmt.Errorf("request: %w", err)
	}

	// A 200 OK response can still time out while reading its body.
	// Decode inside the retry loop so body read errors trigger another attempt.
	decodeResponse := func(resp *http.Response) error {
		if resp.StatusCode != http.StatusOK {
			return nil
		}
		next := wttr.Data{}
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

func checkDescription(target string) string {
	if strings.Contains(strings.ToLower(target), "rain") {
		return "r"
	}

	if strings.Contains(strings.ToLower(target), "snow") {
		return "s"
	}

	return ""
}

func timeConvertFrom12to24H(target string) (string, error) {
	time, err := time.Parse(_timeModel12h, target)
	if err != nil {
		return "", err
	}

	return time.Format(_timeModel24h), nil
}
