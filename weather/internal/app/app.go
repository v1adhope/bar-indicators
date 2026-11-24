package app

import (
	"bytes"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/url"
	"strings"
	"time"

	"github.com/v1adhope/bar-indicators/lib/database"
	"github.com/v1adhope/bar-indicators/lib/history"
	"github.com/v1adhope/bar-indicators/lib/httpclient"
	"github.com/v1adhope/bar-indicators/lib/waybar"
	"github.com/v1adhope/bar-indicators/weather/internal/queires"
	"github.com/v1adhope/bar-indicators/weather/internal/wttr"
)

const (
	_timeModel12h = "03:04 PM"
	_timeModel24h = "15:04"
)

// Run fetches and prints weather for city, using saved results if the fetch fails.
func Run(ctx context.Context, city string) error {
	res, fetchErr := getWeather(ctx, city)
	if fetchErr == nil {
		fmt.Println(res)
	}

	db, err := database.Open(ctx, "bar-indicators/weather_history.sqlite3")
	if err != nil {
		return errors.Join(fetchErr, err)
	}

	defer db.Close()

	if err := db.Init(ctx, queires.Schema); err != nil {
		return errors.Join(fetchErr, err)
	}

	status := history.StatusSuccess

	if fetchErr != nil {
		status, res = history.StatusError, fetchErr.Error()

		fallback, err := db.GetLatestSuccess(ctx, queires.GetLatestSuccess, city, history.StatusSuccess.String())
		if err != nil {
			recordErr := db.Record(ctx, queires.Record, time.Now().UTC().Format(time.RFC3339), city, status.String(), res)

			if errors.Is(err, sql.ErrNoRows) {
				return errors.Join(fetchErr, recordErr)
			}

			return errors.Join(fetchErr, fmt.Errorf("load fallback: %w", err), recordErr)
		}

		log.Printf("%v; using last successful weather result", fetchErr)
		fmt.Println(strings.Replace(fallback, "Fallback not used", "Fallback used", 1))
	}

	return db.Record(ctx, queires.Record, time.Now().UTC().Format(time.RFC3339), city, status.String(), res)
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
	client := httpclient.New()

	data, err := httpclient.FetchData[wttr.Data](ctx, client, u.String(), nil)
	if err != nil {
		return "", err
	}

	module := waybar.Module{}
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

	res, err := module.Encode()
	if err != nil {
		return "", fmt.Errorf("waybar: %w", err)
	}
	return res, nil
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
