
# GO Poker API

The Go Poker API consists of the following 2 sections

 1. Identity REST Server
 2. Room API & WS Server

All endpoints including WebSockets must be prefixed with `/api/v1/`
 
## Identity

Identity generation only requires that users pass in a unique id. For each unique ID an access token can and will be created. As such, Identity can be derived by an adapted OAuth workflow, depending on the individual's needs. The following endpoints exist at the moment to support identity creation


| method | endpoint | description
|--|--|--|
| POST |`/identity/gen`| Generates a new ID for a given user, reserves a username
| POST |`/identity/token`| Fetches a new token for a user provided a secret


## Room


| method | endpoint |
|--|--|--|--|
| GET |`/rooms`|
| POST |`/rooms` 

For an individual attempting to join a room, a WebSocket must be opened at the endpoint located at

`/room/{room_id}/subscribe?token={token}&format={format}`
* `room_id` The room id fetched from the room API
* `token` The user identity token, fetched from the identity server
* `format` The data communication format used in communication, currently supported formats are `JSON`, and `Protobuf`.

After joining a room, a user must follow the protocol outlined in the protobuf documentation `[]`.
