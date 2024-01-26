package cmd

import (
	"os"

	"github.com/dpmcgarry/route53-ddns/pkg"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

// getipCmd represents the getip command
var getipCmd = &cobra.Command{
	Use:   "getip",
	Short: "Gets your Public IP and Logs It",
	Long:  `Gets your Public IP and Logs It`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Info().Msg("Get IP Called")
		ip, err := pkg.GetIP()
		if err != nil {
			log.Fatal().Msgf("Error Getting IP: %v", err)
			os.Exit(1)
		}
		log.Info().Msgf("Got Public IP Address: %v", ip)
	},
}

func init() {
	rootCmd.AddCommand(getipCmd)
}
