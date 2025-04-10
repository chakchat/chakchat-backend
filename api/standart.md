**Here is a standard for API in this project.**
# REST API
`Idemotency-Key` header should be unique for every action but should be repeated when action is retried to ensure idempotency.
## Response body
Every Response can be deserialize as this object:
```json
{
  "error_type": null, // string
  "error_message": null, // string
  "error_details": [], // array of details
  "data": null, // response data. If success then it is object.
}
```
The `error_xxx` fields are `null` when response is successful.
Unused fields are omitted as shown in the examples below.
### Success response with object data
```json
{
  "data": {
    "field1": "value1",
    "field2": "value2"
  }
}
```
### Success response with array data
`data` field should be JSON object because array cabbot be extended by additional fields in the future.
```json
{
  "data": {
    "users": [
      {
        "name": "John Doe",
        "username": "johndoe",
        "phone": "+1234567890",
        "photo": "https://example.com/johndoe.jpg",
        "dateOfBirth": "1990-01-01"
      },
      {
        "name": "Jane Doe",
        "username": "janedoe",
        "phone": "+1234567891",
        "photo": "https://example.com/janedoe.jpg",
        "dateOfBirth": "1990-01-02"
      }
    ]
  }
}
```
For example, if you want to return paginated list of users, you should return something like this:
```json
{
  "data": {
    "entities": [
      {
        "id": "35bdbf25-7715-41d2-b77b-6f69b49ce0a9",
        "name": "John Doe",
        "username": "johndoe",
        "phone": "+1234567890",
        "photo": "https://example.com/johndoe.jpg",
        "dateOfBirth": "1990-01-01"
      }
    ],
    "pagination": {
      "page": 1,
      "size": 1,
      "total": 10
    }
  }
}
```
### Error response
```json
{
  "error_type": "invalid_input",
  "error_message": "Phone number is invalid",
  "error_details": [
    {
      "field": "phone",
      "message": "Phone number is invalid"
    }
  ]
}
```
`error_type` should be machine-readable error type.
### Used `error_type` values
- `internal` - Internal server error. No details provided.
- `invalid_json` - Invalid format of JSON. It can't be parsed correctly.
- `validation_failed` - Some fields are invalid. See `error_details` for details.
- `user_not_found` - Such User does not exist.
- `idempotency_key_missing` - No Idempotency-Key header is provided.
- `send_code_freq_exceeded` - Too many attempts to send code in a short period of time.
- `signin_key_not_found` - Such Sign-in key does not exist.
- `wrong_code` - Wrong verification code is provided.
- `refresh_token_expired` - Refresh JWT token is expired.
- `refresh_token_invalidated` - Refresh JWT token is invalidated.
- `invalid_token` - JWT token is invalid. It can't be parsed correctly or fails some validation not described in other error types.
- `invalid_token_type` - Invalid `typ` field in JWT token.
- `unauthorized` - Invalid authorization or no authorization provided.
- `access_token_expired` - Access token is expired.
- `not_found` - Not found. No such endpoint or resource.
- `user_already_exists` - User already exists
- `signup_key_not_found` - Sign-up key does not exist.
- `username_already_exists` - Username already exists.
- `forbidden` - Forbidden. User is not allowed to perform this action.
- `invalid_param` - Invalid query or route parameter.
