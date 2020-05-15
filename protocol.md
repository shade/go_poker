
#  Communication Protocol

This document describes the message protocol for communication between the poker backend and client. The only purpose of this is to show the message structure and types over the proposed  WebSockets communication channel.
  
  If the communication channel changes from WebSockets this document should be disregarded.

##  General structure

Events are given the following values:
```
CHAT_MSG = 0x80
TABLE_SIT = 0x21
TABLE_STAND = 0x22
```
###  Server to Client ONLY
These messages are only sent from the server to the client .
```
TABLE_TIMER = 0x04
PLAYER_ACTION = 0x05
PLAYER_SHOW = 0x06
HAND_RESULT = 0x07
```
###  Client to Server ONLY
These messages are only sent from the client to the server.
```
HAND_ACTION = 0x23
```

## General Message Structure

The general message structure looks like the following.  

```
{
	type: <EventType>
	timestamp: <Integer> -- UNIX Timestamp
	payload: <Object> -- The JSON payloads as defined below
}
```
  
 
## Sit
Event `TABLE_SIT`, allows a player to sit in a table given a table ID.

#### Payload
```
{
	table_id: <String>,
	user_identity_token: <String>
	seat_num: <Integer>
	chips: <Integer>
}
```

#### Potential Responses
`TABLE_SIT`
```
{
	table_id: <String>,
	sat_down: <Boolean>,
	seat_num: <Integer>
	reason: <String>
}
sat_down will hold boolean value of whether player successfully sat down at seat_num
```

## Stand
Event `TABLE_STAND`, allows a player to stand when they want to leave a table, irrespective of seat number

#### Payload
```
{
	table_id: <String>
}
```

#### Potential Responses
`TABLE_STAND`
```
{
	table_id: <String>
	balance: <Integer> - The amount of money returned to the user in chips ??????????
	left_seat: <Boolean>
	reason: <String>
}
```
