
#  Communication Protocol

This document describes the message protocol for communication between the poker backend and client. The only purpose of this is to show the message structure and types over the proposed  WebSockets communication channel.
  
  If the communication channel changes from WebSockets this document should be disregarded.

##  General structure

Events are given the following values:
```
CHAT_MSG = 0x80
```
###  Server to Client ONLY
These messages are only sent from the server to the client .
```
TABLE_ACCEPT = 0x01
TABLE_DECLINE = 0x02
TABLE_TIMER = 0x03
PLAYER_ACTION = 0x04
PLAYER_SHOW = 0x05
HAND_RESULT = 0x06
```
###  Client to Server ONLY
These messages are only sent from the client to the server.
```
TABLE_SIT = 0x21
TABLE_STAND = 0x22
HAND_ACTION = 0x23
```



## General Message Structure

The general message structure looks like the following.  

```json
{
	"type": <EventType>
	"timestamp": <Integer>
	"payload": <Object>
}
```
  
 
## Client To Server

### Sit
Event `TABLE_SIT`, allows a player to sit in a table given a table ID.

#### Payload
```
{
	table_id: <String>,
	user_identity_token: <String>
}
```

#### Potential Responses
**Event: TABLE_DECLINE**
```
{
	table_id: <String>,
	reason: "User already sitting at table"
}
```
**Event: TABLE_ACCEPT**
Accepting will give the user their position in the table.
```
{
	table_id: <String>,
	position: <Integer>
}
```

### Stand

TBD
## Server to Client