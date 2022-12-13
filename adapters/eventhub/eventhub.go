package eventhub

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	eh "github.com/Azure/azure-event-hubs-go/v3"
)

type Receiver interface {
	Receive(context.Context, func(c context.Context, event *eh.Event) error) error
}

type EventHubAdapter struct {
	hub *eh.Hub
}

func (a EventHubAdapter) Receive(
	ctx context.Context,
	handler func(c context.Context, event *eh.Event) error) error {

	runtimeInfo, err := a.hub.GetRuntimeInformation(ctx)
	if err != nil {
		fmt.Println(err)
		return err
	}

	fmt.Println("number of partitions:", len(runtimeInfo.PartitionIDs))

	// start receiving on all partitions
	for _, partitionID := range runtimeInfo.PartitionIDs {
		fmt.Println("Receiving on partition:", partitionID)
		_, err := a.hub.Receive(ctx, partitionID, handler)
		if err != nil {
			fmt.Println(err)
			return err
		}
	}

	// wait for a signal to quit
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)
	<-signalChan

	// clean up
	err = a.hub.Close(context.Background())
	if err != nil {
		fmt.Println(err)
		return nil
	}

	fmt.Println("Done")

	return nil
}

// New creates a new EventHubAdapter using the environment variables EVENTHUB_*.
func New() (EventHubAdapter, error) {
	hub, err := eh.NewHubFromEnvironment()
	if err != nil {
		return EventHubAdapter{}, err
	}

	return EventHubAdapter{hub: hub}, nil
}
