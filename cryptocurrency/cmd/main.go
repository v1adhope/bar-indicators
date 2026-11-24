package main

import (
	"errors"
	"log"
	"os"

	"github.com/v1adhope/bar-indicators/cryptocurrency/internal/app"
)

const (
	_envApiToken     = "CRYPTO_CURRENCY_API"
	_envFavoriteCoin = "CRYPTO_CURRENCY_FV"
	_envWatchList    = "CRYPTO_CURRENCY_WL"
)

func main() {
	convertId, rawCoinIds := "2781", "1,1027,6636"

	apiToken := os.Getenv(_envApiToken)

	if err := fillVarByEnvKey(_envFavoriteCoin, &convertId); err != nil {
		log.Fatal(err)
	}

	if err := fillVarByEnvKey(_envWatchList, &rawCoinIds); err != nil {
		log.Fatal(err)
	}

	if err := app.Run(apiToken, convertId, rawCoinIds); err != nil {
		log.Fatal(err)
	}
}

func fillVarByEnvKey(key string, placeholder *string) error {
	if env := os.Getenv(key); env != "" {
		for _, v := range env {
			if byte(v) == 44 || (byte(v) < 58 && byte(v) > 47) {
				continue
			}

			return errors.New("incorrect by key: " + key)
		}
		*placeholder = env
	}

	return nil
}
