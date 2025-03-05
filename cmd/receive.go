/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// receiveCmd represents the receive command
var receiveCmd = &cobra.Command{
	Use:   "receive",
	Short: "receive events from event hub",
	Long:  `receive events from event hub`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("receive called")
		connectionString := viper.GetString("connection-string")
		eventhubNamespace := viper.GetString("namespace")
		eventhubName := viper.GetString("name")
		fmt.Printf("connectionString: %s\n", connectionString)
		fmt.Printf("eventhubNamespace: %s\n", eventhubNamespace)
		fmt.Printf("eventhubName: %s\n", eventhubName)
	},
}

func init() {
	eventsCmd.AddCommand(receiveCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// receiveCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// receiveCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
