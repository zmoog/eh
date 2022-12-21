package commands_test

import (
	"context"
	"testing"

	"github.com/zmoog/eh/commands"
)

func TestTailCommand(t *testing.T) {
	t.Run("TailCommand", func(t *testing.T) {
		uow := commands.UnitOfWork{}

		cmd := commands.TailCommand{}

		err := cmd.ExecuteWith(context.Background(), uow)
		if err != nil {
			t.Fail()
		}
	})
}
