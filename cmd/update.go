package cmd

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/route53"
	"github.com/aws/aws-sdk-go-v2/service/route53/types"
	"github.com/dpmcgarry/route53-ddns/pkg"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// updateCmd represents the update command
var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Upserts DNS Entries with Public IP",
	Long:  `Upserts DNS Entries with Public IP`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Info().Msgf("Update called")
		awsConf, err := pkg.GetAWSConfig()
		if err != nil {
			log.Fatal().Msgf("Error getting AWS Config: %v", err)
			os.Exit(1)
		}

		records, err := pkg.GetRecords()
		if err != nil {
			log.Fatal().Msgf("Error getting DNS Records to Update: %v", err)
			os.Exit(1)
		}
		ip, err := pkg.GetIP()
		if err != nil {
			log.Fatal().Msgf("Error Getting IP: %v", err)
			os.Exit(1)
		}
		log.Info().Msgf("Got Public IP Address: %v", ip)
		log.Info().Msgf("AWS Config: %v", awsConf)
		log.Info().Msgf("Records: %v", records)
		log.Info().Msg("Creating AWS Client")
		cfg, err := config.LoadDefaultConfig(context.TODO(),
			config.WithSharedConfigProfile(awsConf["profile"]))
		if err != nil {
			log.Fatal().Msgf("error creating AWS client: %v", err)
			os.Exit(1)
		}

		ttl := viper.GetInt64("ttl")

		var changes []types.Change
		for _, record := range records {
			log.Info().Msgf("Creating changeset for: %v", record)
			resourcerecords := []types.ResourceRecord{}
			resourcerecord := types.ResourceRecord{}
			resourcerecord.Value = aws.String(ip)
			resourcerecords = append(resourcerecords, resourcerecord)
			change := types.Change{}
			change.Action = types.ChangeActionUpsert
			rrset := types.ResourceRecordSet{}
			rrset.Name = aws.String(record)
			rrset.Type = types.RRTypeA
			rrset.TTL = aws.Int64(ttl)
			rrset.ResourceRecords = resourcerecords
			change.ResourceRecordSet = &rrset
			changes = append(changes, change)
		}
		for _, change := range changes {
			log.Debug().Msg("Dumping changes")
			log.Debug().Msg(string(change.Action))
			log.Debug().Msg(*change.ResourceRecordSet.Name)
			log.Debug().Msg(string(change.ResourceRecordSet.Type))
			log.Debug().Msg(fmt.Sprint(*change.ResourceRecordSet.TTL))
			for _, rr := range change.ResourceRecordSet.ResourceRecords {
				log.Debug().Msg(*rr.Value)
			}

		}
		route53client := route53.NewFromConfig(cfg)
		r53resp, err := route53client.ChangeResourceRecordSets(context.TODO(), &route53.ChangeResourceRecordSetsInput{
			HostedZoneId: aws.String(awsConf["hostedzoneid"]),
			ChangeBatch: &types.ChangeBatch{
				Comment: aws.String(fmt.Sprintf("Updated by route53-ddns at %v", time.Now().UTC())),
				Changes: changes,
			},
		})
		if err != nil {
			log.Fatal().Msgf("Error Calling Route53: %v", err)
		}
		changeid := *r53resp.ChangeInfo.Id
		log.Info().Msgf("Got Change ID: %v", changeid)
		retries := viper.GetInt("retries")
		waitseconds := viper.GetInt("waitseconds")
		log.Debug().Msgf("Will check changeset for %v retries waiting %v seconds", retries, waitseconds)

		for i := 0; i < retries; i++ {
			time.Sleep(time.Duration(waitseconds) * time.Second)
			log.Info().Msg("Getting change status from route53")
			r53resp, err := route53client.GetChange(context.TODO(), &route53.GetChangeInput{
				Id: aws.String(changeid),
			})
			if err != nil {
				log.Warn().Msgf("Got Error Calling Route53: %v", err)
			}
			log.Debug().Msgf("Got Status %v", r53resp.ChangeInfo.Status)
			if r53resp.ChangeInfo.Status == types.ChangeStatusInsync {
				log.Info().Msg("Change is in sync! My work here is done.")
				i = retries
			} else {
				log.Info().Msg("Change is still pending. Will retry after sleeping.")
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(updateCmd)
}
