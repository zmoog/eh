# eh

**eh** is a command line tool to dump messages with Azure Event Hub.

```shell
EVENTHUB_CONNECTION_STRING='Endpoint=sb://whatever'
EVENTHUB_NAME="whatever"
EVENTHUB_NAMESPACE="whatever"

$ eh
event: [{"records": [{ "LogicalServerName": "DummyValue", "SubscriptionId": "00000000-0000-0000-0000-000000000000", "ResourceGroup": "DummyValue", "time": "2022-05-03T09:09:56.5496795Z", ... }]}]
```

I am currently using it to print the original events published to my test Event Hub.

## Refs

- [Event Hub Go client library](https://github.com/Azure/azure-event-hubs-go)
- [Event Hub service](https://azure.microsoft.com/en-us/services/event-hubs/)