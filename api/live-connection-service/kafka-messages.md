# Kafka message

## Message from messaging service
Information about new message should have format:
```json
{
     "receivers": "usersIds",
     "data": {
        "type": "send/edit/delete",
        "update_id": "id",
        "message_id": "messageId",
        "text": "message_text", //null if it is file
        "file": {
          "file_name": "string",
          "file_size": 0,
          "mime_type": "string",
          "file_id": "3fa85f64-5717-4562-b3fc-2c963f66afa6",
          "file_url": "string",
          "created_at": "2025-03-26T10:14:33.367Z"
        }, // nul if it is text message
        "sender_id": "userId",
        "created_at": "dateTime"
    }   
}
```

Every chat's update should have format

```json
{
  "receivers": "usersIds",
  "data":{
    "chat_id": "chatId",
    "update_id": "id",
    "update_type": "created/deleted/blocked/unblocked/change",
    "chat_type": "type",
    "description": "description",
    "group_photo": "group photo url",
    "created_at": "dateTime"
  }
}
```

Information about group members should have format:

```json
{
  "receivers": "usersIds",
  "data": {
    "chat_id": "chatId",
    "update_type": "add/remove",
    "member_id": "userId"

  }
}

Information about reactions should have format:

```json
{
  "receivers": "usersIds",
  "data": {
    "chat_id": "id",
    "update_id": "id",
    "created_at": "dateTime"
  }
}




