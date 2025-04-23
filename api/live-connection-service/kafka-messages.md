# Modification messages

All messages have such structure:

```json
{
  "receivers": [],
  "type": "",
  "data": {} // It will hold specific data
}
```

Possible types:

```
update

chat_created
chat_deleted
chat_blocked
chat_unblocked
chat_expiration_set
group_info_updated
group_members_added
group_members_removed

```

# Update

Information about new message should have format:
[Code](https://github.com/chakchat/chakchat-backend/tree/main/messaging-service/internal/rest/response/generic_update.go)

All update messages have structure:

```json
{
  "receivers": [],
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

For text messages structure is:

```json
{
  "receivers": ["57a85f64-5717-4562-b3fc-2c54636a123"],
  "type": "update",
  "data": {
    "chat_id": "b4c3f591-ef52-4b85-a04e-cf61ee243449",
    "update_id": 123,
    "type": "text_message",
    "sender_id": "a1b7a452-6ef4-4d56-9d65-cd80f207b157",
    "created_at": "2025-12-04T17:00:00",
    "content": {
      "text": "Hello, Anna!",
      "edited": {
        "chat_id": "b4c3f591-ef52-4b85-a04e-cf61ee243449",
        "update_id": 125,
        "type": "text_message_edited",
        "sender_id": "a1b7a452-6ef4-4d56-9d65-cd80f207b157",
        "created_at": "2025-12-04T17:01:00",
        "content": {
          "new_text": "Hey, Liza! (Sorry, Anna)",
          "message_id": 123
        },      
      }, // Can be null
      "reply_to": "2d49b5fa-8a26-48ce-ac8f-6463da227fc7",
      "reactions": [
        {
          "chat_id": "b4c3f591-ef52-4b85-a04e-cf61ee243449",
          "update_id": 124,
          "type": "reaction",
          "sender_id": "a1b7a452-6ef4-4d56-9d65-cd80f207b157",
          "created_at": "12025-12-04T17:00:30",
          "content": {
            "reaction": "black_smile",
            "message_id": 123
          }
        },
        {
          "chat_id": "b4c3f591-ef52-4b85-a04e-cf61ee243449",
          "update_id": 126,
          "type": "reaction",
          "sender_id": "a1b7a452-6ef4-4d56-9d65-cd80f207b157",
          "created_at": "12.32.2323/12:234:122Z+2",
          "content": {
            "reaction": "black_smile",
            "message_id": 123
          }
        }
      ] // Can be null
    }
  }
}
```

All file messages have such structure:

```json
{
  "receivers": ["57a85f64-5717-4562-b3fc-2c54636a123"],
  "type": "update",
  "data": {
    "chat_id": "b4c3f591-ef52-4b85-a04e-cf61ee243449",
    "update_id": 126,
    "type": "file",
    "sender_id": "a1b7a452-6ef4-4d56-9d65-cd80f207b157",
    "created_at": "12.32.2323/12:234:122Z+2",
    "content": {
      "file_id": "7bed2b32-01ac-43a6-abd7-fe037de495c4",
      "file_name": "photo",
      "mime_type": "string",
      "file_size": 0,
      "file_url": "https://s3.our.com/file-bucket/8c7dcdb3-f671-4cfd-96b5-10d350da13ee",
      "message_id": 123
    }
  }
}
```

All reactions have such structure:
```json
{
  
  "receivers": ["57a85f64-5717-4562-b3fc-2c54636a123"],
  "type": "update",
  "data": {
    "chat_id": "b4c3f591-ef52-4b85-a04e-cf61ee243449",
    "update_id": 126,
    "type": "reaction",
    "sender_id": "a1b7a452-6ef4-4d56-9d65-cd80f207b157",
    "created_at": "12.32.2323/12:234:122Z+2",
    "content": {
      "reaction": "black_smile",
      "message_id": 123
    }
  }
}
```

All deleting message have such structure:
```json
{
  "receivers": ["57a85f64-5717-4562-b3fc-2c54636a123"],
  "type": "update",
  "data": {
    "chat_id": "b4c3f591-ef52-4b85-a04e-cf61ee243449",
    "update_id": 128,
    "type": "delete",
    "sender_id": "a1b7a452-6ef4-4d56-9d65-cd80f207b157",
    "created_at": "12.32.2323/12:234:122Z+2",
    "content": {
      "deleted_id": 123,
      "deleted_mode": "all/only_me"
    }
  }
}
```

# Chat

Information about chat updates should have format:
[Code](https://github.com/chakchat/chakchat-backend/tree/main/messaging-service/internal/rest/response/generic_chat.go)

All update chat messages have the same structure as update messages(previous).

## Created chat

```json
{

  "receivers": [ "57a85f64-5717-4562-b3fc-2c54636a123"],
  "type": "chat_created",
  "data":{
    "sender_id": "32f8c01b-673c-4b3b-a42a-b84fbaf10bff",
    "chat": {
      "chat_id": "566bfca7-3ab0-4242-98b2-61d459acd879",
      "type": "group",
      "members": [],
      "created_at": "RFC3339",
      "info": {
        "admin_id": "32f8c01b-673c-4b3b-a42a-b84fbaf10bff",
        "name": "Ann",
        "description": "It is tayga chat",
        "group_photo": null,
      }
    }
  }
}
```

## Delete/Block/Unblock chat
```json
{
  "receivers": ["57a85f64-5717-4562-b3fc-2c54636a123"],
  "type": "chat_deleted/chat_blocked/chat_unblocked",
  "data": {
    "sender_id": "9994d052-3fc3-42de-be9c-0d692b6a0e39",
    "chat_id": "566bfca7-3ab0-4242-98b2-61d459acd879"
  }
}
```
## Chat expiration set

```json
{
  "receivers": ["57a85f64-5717-4562-b3fc-2c54636a123"],
  "type": "chat_expiration_set",
  "data": {
    "sender_id": "9994d052-3fc3-42de-be9c-0d692b6a0e39",
    "chat_id": "566bfca7-3ab0-4242-98b2-61d459acd879",
    "expiration": null
  }
}
```

## Update group info

Update only info about group not about members.

All fields are passed.

```json
{
  "receivers": ["57a85f64-5717-4562-b3fc-2c54636a123"],
  "type": "group_info_updated",
  "data": {
    "sender_id": "9994d052-3fc3-42de-be9c-0d692b6a0e39",
    "chat_id": "566bfca7-3ab0-4242-98b2-61d459acd879",
    "name": "Super Group name",
    "description": "", // Description removed
    "group_photo": "https://s3.our.com/file-bucket/8c7dcdb3-f671-4cfd-96b5-10d350da13ee"
  }
}
```

## Update group members


```json
{
  "receivers": ["57a85f64-5717-4562-b3fc-2c54636a123"],
  "type": "group_members_added/group_members_removed",
  "data": {
    "sender_id": "9994d052-3fc3-42de-be9c-0d692b6a0e39",
    "chat_id": "566bfca7-3ab0-4242-98b2-61d459acd879",
    "members": ["4bf2ac2a-1a4c-48fc-ac64-4e9418107c49"] // Only added/removed members
  }
}
```