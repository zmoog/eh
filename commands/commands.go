package commands

import "github.com/zmoog/eh/adapters/eventhub"

type UnitOfWork struct {
	Receiver eventhub.Receiver
}
