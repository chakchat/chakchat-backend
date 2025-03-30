# WebSocket messages

## Ping body

Every 5 seconds client sends ping to server

```json
{
  "timestamp": "2025-03-19T12:00:00Z"
}
```

After receiving information from messaging service, live-connection service sends response to client in the same format as kafka-messages