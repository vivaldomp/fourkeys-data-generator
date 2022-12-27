package main

import (
	"fmt"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile       string
	eventTimespan int
	numEvents     int
	numIssues     int
	rootCmd       = &cobra.Command{
		Use:   "generator",
		Short: "generator is a command cli",
		Long:  `generator is a command cli`,
		Run: func(cmd *cobra.Command, args []string) {
			Start(GetConfig(), eventTimespan, numEvents, numIssues)
		},
	}
)

func Execute() {
	err := rootCmd.Execute()
	cobra.CheckErr(err)
}

func init() {

	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.env)")
	rootCmd.PersistentFlags().IntVarP(&eventTimespan, "event-timespan", "t", 604800, `time duration (in seconds) of timestamps of generated events
(from [Now-timespan] to [Now]); default=604800 (1 week)`)
	rootCmd.PersistentFlags().IntVarP(&numEvents, "num-events", "e", 40, "number of events to generate; default=40")
	rootCmd.PersistentFlags().IntVarP(&numIssues, "num-issues", "i", 2, "number of issues to generate; default=2")

	viper.SetDefault("WEBHOOK", "http://localhost:8080")
	viper.SetDefault("SECRET", "123456")
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		viper.AddConfigPath(".")
		viper.SetConfigType("env")
		viper.SetConfigName(".env")
	}

	replacer := strings.NewReplacer(".", "_")
	viper.SetEnvKeyReplacer(replacer)
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file: ", viper.ConfigFileUsed())
	} else {
		log.Errorf("Error: %v", err)
	}

}
