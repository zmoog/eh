/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azeventhubs"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/time/rate"
)

var (
	rateLimit rate.Limit
)

// sendCmd represents the send command
var sendCmd = &cobra.Command{
	Use:   "send",
	Short: "send events to event hub",
	Long:  `send events to event hub using parallel processing for maximum throughput`,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if viper.GetString("rate-limit") == "no" {
			// If the rate limit is not set, we set it to no limit
			rateLimit = rate.Inf
			return nil
		}

		// Parse the rate limit from the command line argument
		r, err := time.ParseDuration(viper.GetString("rate-limit"))
		if err != nil {
			return fmt.Errorf("failed to parse rate limit: %w", err)
		}

		rateLimit = rate.Every(r)

		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		var reader *bufio.Reader

		switch input := viper.GetString("input"); input {
		case "-":
			fmt.Println("input", input)
			reader = bufio.NewReader(os.Stdin)
		default:
			fmt.Println("input", input)
			file, err := os.Open(input)
			if err != nil {
				return fmt.Errorf("failed to open input file: %w", err)
			}
			defer file.Close()
			reader = bufio.NewReader(file)
		}

		connectionString := viper.GetString("connection-string")
		eventHubName := viper.GetString("name")

		// defaultAzureCred, err := azidentity.NewDefaultAzureCredential(nil)
		// if err != nil {
		// 	return fmt.Errorf("failed to create default azure credential: %w", err)
		// }

		// Can also use a connection string:
		//
		producerClient, err := azeventhubs.NewProducerClientFromConnectionString(connectionString, eventHubName, nil)
		//
		// producerClient, err := azeventhubs.NewProducerClient(eventHubNamespace, eventHubName, defaultAzureCred, nil)

		if err != nil {
			return fmt.Errorf("failed to create producer client: %w", err)
		}

		defer producerClient.Close(context.TODO())

		newBatchOptions := &azeventhubs.EventDataBatchOptions{
			// The options allow you to control the size of the batch, as well as the partition it will get sent to.

			// PartitionID can be used to target a specific partition ID.
			// specific partition ID.
			//
			// PartitionID: partitionID,

			// PartitionKey can be used to ensure that messages that have the same key
			// will go to the same partition without requiring your application to specify
			// that partition ID.
			//
			// PartitionKey: partitionKey,

			//
			// Or, if you leave both PartitionID and PartitionKey nil, the service will choose a partition.
		}

		fmt.Println("creating event data batch")

		// Creates an EventDataBatch, which you can use to pack multiple events together, allowing for efficient transfer.
		batch, err := producerClient.NewEventDataBatch(context.TODO(), newBatchOptions)
		if err != nil {
			return fmt.Errorf("failed to create event data batch: %w", err)
		}

		// Create a rate limiter that allows 1 event per second
		limiter := rate.NewLimiter(rateLimit, viper.GetInt("rate-burst"))

		lineCount := 0
		for {
			fmt.Println("reading line", lineCount)
			line, err := reader.ReadString('\n')
			if err != nil {
				if err == io.EOF {
					break
				}
				return fmt.Errorf("failed to read line: %w", err)
			}
			lineCount++

			fmt.Println("read line with length", len(line))
			fmt.Println("adding event data to batch", batch.NumEvents())

			// Wait for the rate limiter before adding the event to the batch
			if err := limiter.Wait(context.TODO()); err != nil {
				return fmt.Errorf("rate limiter error: %w", err)
			}

			err = batch.AddEventData(&azeventhubs.EventData{Body: []byte(line)}, nil)
			if errors.Is(err, azeventhubs.ErrEventDataTooLarge) {
				if batch.NumEvents() == 0 {
					// This one event is too large for this batch, even on its own. No matter what we do it
					// will not be sendable at its current size.
					return fmt.Errorf("failed to create event data batch: %w", err)
				}

				// This batch is full - we can send it and create a new one and continue
				fmt.Printf("batch is full, flushing %d events\n", batch.NumEvents())

				// Wait for the rate limiter before sending the batch
				if err := limiter.Wait(context.TODO()); err != nil {
					return fmt.Errorf("rate limiter error: %w", err)
				}

				if err := producerClient.SendEventDataBatch(context.TODO(), batch, nil); err != nil {
					return fmt.Errorf("failed to send event data batch: %w", err)
				}
				fmt.Println("✅ sent event data batch")

				// create the next batch we'll use for events, ensuring that we use the same options
				// each time so all the messages go the same target.
				fmt.Println("recreating batch")

				tmpBatch, err := producerClient.NewEventDataBatch(context.TODO(), newBatchOptions)
				if err != nil {
					return fmt.Errorf("failed to create event data batch: %w", err)
				}

				batch = tmpBatch

				// Wait for the rate limiter before adding the event to the batch
				if err := limiter.Wait(context.TODO()); err != nil {
					return fmt.Errorf("rate limiter error: %w", err)
				}

				fmt.Println("adding event data to batch", batch.NumEvents())
				err = batch.AddEventData(&azeventhubs.EventData{Body: []byte(line)}, nil)
				if err != nil {
					return fmt.Errorf("failed to add event data to batch: %w", err)
				}

			} else if err != nil {
				// This is a different error - we can't add this event to the batch, but we can
				return fmt.Errorf("failed to add event data to batch: %w", err)
			}
		}

		// if we have any events in the last batch, send it
		if batch.NumEvents() > 0 {
			fmt.Printf("flushing remaining %d events\n", batch.NumEvents())
			if err := producerClient.SendEventDataBatch(context.TODO(), batch, nil); err != nil {
				return fmt.Errorf("failed to flush remaining events: %w", err)
			}
		}

		if err := producerClient.Close(context.TODO()); err != nil {
			return fmt.Errorf("failed to close producer client: %w", err)
		}

		return nil
	},
}

func init() {
	eventsCmd.AddCommand(sendCmd)

	sendCmd.Flags().StringP("input", "i", "-", "Input JSON file (one JSON object per line)")
	// sendCmd.Flags().StringP("connection-string", "c", "", "Event Hub connection string")
	// sendCmd.Flags().IntP("workers", "w", 10, "Number of parallel workers")
	sendCmd.Flags().StringP("rate-limit", "r", "no", "Rate limit for sending events (e.g. 1s, 100ms, 1h)")
	sendCmd.Flags().IntP("rate-burst", "b", 1, "Rate burst for sending events (e.g. 1s, 100ms, 1h)")

	_ = viper.BindPFlag("input", sendCmd.Flags().Lookup("input"))
	_ = viper.BindPFlag("rate-limit", sendCmd.Flags().Lookup("rate-limit"))
	_ = viper.BindPFlag("rate-burst", sendCmd.Flags().Lookup("rate-burst"))

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// sendCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// sendCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
