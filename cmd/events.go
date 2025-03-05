/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// eventsCmd represents the events command
var eventsCmd = &cobra.Command{
	Use:   "events",
	Short: "events related commands",
	Long:  `events related commands`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("events called")
	},
}

func init() {
	rootCmd.AddCommand(eventsCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// eventsCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// eventsCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
