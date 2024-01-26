package main

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
)

func configureViper() {
	viper.SetConfigName("route53-ddns")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("/etc/")
	viper.AddConfigPath(".")
	viper.SetDefault("logdir", "./")
	viper.SetDefault("retries", 20)
	viper.SetDefault("waitseconds", 10)
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			fmt.Printf("Config File not found in /etc or current dir: %v\n", err)
			os.Exit(1)
		} else {
			fmt.Printf("Error reading config file but file was found: %v\n", err)
			os.Exit(1)
		}
	}

}
