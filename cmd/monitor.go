package cmd

import (
	"log"
	"os"
	"sync"

	"github.com/DisgoOrg/disgo/webhook"
	"github.com/DisgoOrg/snowflake"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

var monitorCmd = &cobra.Command{
	Use:   "monitor",
	Short: "Daemon to monitor validators",
	Long:  "Monitors validators and pushes alerts to Discord using the configuration in config.yaml",
	Run: func(cmd *cobra.Command, args []string) {
		dat, err := os.ReadFile(configFilePath)
		if err != nil {
			log.Fatalf("Error reading config.yaml: %v", err)
		}
		config := HalfLifeConfig{}
		err = yaml.Unmarshal(dat, &config)
		if err != nil {
			log.Fatalf("Error parsing config.yaml: %v", err)
		}
		writeConfigMutex := sync.Mutex{}
		discordClient := webhook.NewClient(snowflake.Snowflake(config.Discord.Webhook.ID), config.Discord.Webhook.Token)
		alertState := make(map[string]*ValidatorAlertState)
		for i, vm := range config.Validators {
			if i == len(config.Validators)-1 {
				runMonitor(&alertState, discordClient, &config, vm, &writeConfigMutex)
			} else {
				go runMonitor(&alertState, discordClient, &config, vm, &writeConfigMutex)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(monitorCmd)
}
