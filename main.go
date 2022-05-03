package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"time"

	eventhub "github.com/Azure/azure-event-hubs-go/v3"
)

func main() {
	eventHub, err := eventhub.NewHubFromEnvironment()
	if err != nil {
		fmt.Println(err)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	handler := func(c context.Context, event *eventhub.Event) error {
		fmt.Printf("event: [%v]\n", string(event.Data))
		return nil
	}

	runtimeInfo, err := eventHub.GetRuntimeInformation(ctx)
	if err != nil {
		fmt.Println(err)
		return
	}

	// start receiving on all partitions
	for _, partitionID := range runtimeInfo.PartitionIDs {
		_, err := eventHub.Receive(ctx, partitionID, handler)
		if err != nil {
			fmt.Println(err)
			return
		}
	}

	// wait for a signal to quit
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, os.Kill)
	<-signalChan

	// clean up
	err = eventHub.Close(context.Background())
	if err != nil {
		fmt.Println(err)
		return
	}
}
