package pkg

import (
	"fmt"
	"io"
	"net/http"

	"github.com/rs/zerolog/log"
)

func GetIP() (string, error) {
	res, err := http.Get(IPAPIUrl)
	if err != nil {
		log.Error().Msgf("Error with API Request: %v", err)
		return "", err
	}
	log.Debug().Msgf("HTTP Returned Headers: %v", res.Header)
	log.Debug().Msgf("HTTP Response Code: %v", res.StatusCode)
	if res.StatusCode > 299 {
		log.Error().Msgf("Got HTTP Response %v", res.Status)
		return "", fmt.Errorf("got unexpected HTTP Response %v", res.Status)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Error().Msgf("Error reading data from HTTP Body: %v", err)
		return "", err
	}
	res.Body.Close()

	return string(body), nil
}
