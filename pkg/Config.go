package pkg

import (
	"errors"

	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

func GetAWSConfig() (map[string]string, error) {
	awsConf := viper.GetStringMapString("aws")

	// For now just support an aws cli profile name
	// In the future can add support for other credential methods
	if _, ok := awsConf["profile"]; !ok {
		log.Error().Msg("AWS Profile Not Found in Config")
		return nil, errors.New("aws profile not found in config")
	}
	if _, ok := awsConf["hostedzoneid"]; !ok {
		log.Error().Msg("Route53 Hosted Zone ID Not Found in Config")
		return nil, errors.New("route53 hosted zone id not found in config")
	}

	log.Info().Msgf("Parsed Config: %v", awsConf)
	return awsConf, nil
}

func GetRecords() ([]string, error) {
	records := viper.GetStringSlice("records")
	if len(records) < 1 {
		return nil, errors.New("no records found in config")
	}
	return records, nil
}
