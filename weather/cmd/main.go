package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/v1adhope/bar-indicators/weather/internal/app"
)

const _runTimeout = time.Minute

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), _runTimeout)
	defer cancel()

	city := os.Getenv("CITY_WEATHER")

	if err := app.Run(ctx, city); err != nil {
		log.Fatal(err)
	}
}
