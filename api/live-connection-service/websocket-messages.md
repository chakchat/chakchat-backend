# WebSocket messages

## Ping

WebSocket implements ping pong messages that should provide server and client with information about connection live.

After receiving information from messaging service, live-connection service sends response to client in the same format as kafka-messages:

```json
{
  "type": "update",
  "data": {
    "chat_id": "1342ffe5-26da-4289-b942-9a3219355b7e",
    "update_id": 123,
    "type": "text_message",
    "sender_id": "a1b7a452-6ef4-4d56-9d65-cd80f207b157",
    "created_at": "2025-12-04T17:00:00",
    "content": {} // Specified content
    }
}
```