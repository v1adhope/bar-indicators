package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"regexp"
	"time"

	"github.com/v1adhope/bar-indicators/cryptocurrency/internal/app"
)

const _runTimeout = time.Minute

var (
	validID  = regexp.MustCompile(`^[0-9]+$`)
	validIDs = regexp.MustCompile(`^[0-9]+(,[0-9]+)*$`)
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), _runTimeout)
	defer cancel()

	apiToken := os.Getenv("CRYPTO_CURRENCY_API")

	convertId, err := parseEnv("CRYPTO_CURRENCY_FV", "2781", validID)
	if err != nil {
		log.Fatal(err)
	}

	rawCoinIds, err := parseEnv("CRYPTO_CURRENCY_WL", "1,20396", validIDs)
	if err != nil {
		log.Fatal(err)
	}

	if err := app.Run(ctx, apiToken, convertId, rawCoinIds); err != nil {
		log.Fatal(err)
	}
}

func parseEnv(key, defaultValue string, pattern *regexp.Regexp) (string, error) {
	env := os.Getenv(key)
	if env == "" {
		return defaultValue, nil
	}
	if !pattern.MatchString(env) {
		return "", fmt.Errorf("invalid value %q for %s", env, key)
	}

	return env, nil
}
