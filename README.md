# eh

**eh** is a command line tool the connects to and Event Hub and dumps the messages to stdout.

## Usage

```shell
EVENTHUB_CONNECTION_STRING='Endpoint=sb://whatever'
EVENTHUB_NAME="whatever"
EVENTHUB_NAMESPACE="whatever"

$ go run entrypoints/cli/main.go tail
event: [{"records": [{ "LogicalServerName": "DummyValue", "SubscriptionId": "00000000-0000-0000-0000-000000000000", "ResourceGroup": "DummyValue", "time": "2022-05-03T09:09:56.5496795Z", ... }]}]
```

It is a simple tool used to examine the messages in an Event Hub during development.


## Refs

- [Event Hub Go client library](https://github.com/Azure/azure-event-hubs-go)
- [Event Hub service](https://azure.microsoft.com/en-us/services/event-hubs/)
