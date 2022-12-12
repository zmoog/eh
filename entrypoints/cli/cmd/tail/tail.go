package tail

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/zmoog/eh/adapters/eventhub"
	"github.com/zmoog/eh/commands"
)

func NewCommand() *cobra.Command {
	cmd := cobra.Command{
		Use:   "tail",
		Short: "tail the messages from an event hub",
		RunE:  runTailCommand,
	}

	return &cmd
}

func runTailCommand(cmd *cobra.Command, args []string) error {
	fmt.Println("tail command")

	eventHubAdapter, err := eventhub.New()
	if err != nil {
		return err
	}

	uow := commands.UnitOfWork{
		Receiver: eventHubAdapter,
	}

	tailCmd := commands.TailCommand{}

	if err := tailCmd.ExecuteWith(context.Background(), uow); err != nil {
		return err
	}

	return nil
}
