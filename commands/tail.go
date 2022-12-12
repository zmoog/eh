package commands

import (
	"context"
	"fmt"

	eh "github.com/Azure/azure-event-hubs-go/v3"
)

type TailCommand struct{}

func (c TailCommand) ExecuteWith(ctx context.Context, uow UnitOfWork) error {
	if err := uow.Receiver.Receive(ctx, func(c context.Context, event *eh.Event) error {
		fmt.Println(string(event.Data))
		return nil
	}); err != nil {
		return err
	}

	return nil
}
