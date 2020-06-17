# Protocol Documentation
<a name="top"></a>

## Table of Contents

- [messages.proto](#messages.proto)
    - [ActionAck](#messages.ActionAck)
    - [ActionMsg](#messages.ActionMsg)
    - [ChatMsgRecv](#messages.ChatMsgRecv)
    - [ChatMsgSend](#messages.ChatMsgSend)
    - [Packet](#messages.Packet)
    - [PlayerSeat](#messages.PlayerSeat)
    - [SitAck](#messages.SitAck)
    - [SitEvent](#messages.SitEvent)
    - [StandAck](#messages.StandAck)
    - [StandEvent](#messages.StandEvent)
  
    - [ActionType](#messages.ActionType)
    - [EventType](#messages.EventType)
  
- [table.proto](#table.proto)
    - [Card](#table.Card)
    - [CardSet](#table.CardSet)
    - [Player](#table.Player)
    - [TableOptions](#table.TableOptions)
    - [TableState](#table.TableState)
  
- [Scalar Value Types](#scalar-value-types)



<a name="messages.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## messages.proto



<a name="messages.ActionAck"></a>

### ActionAck



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| ok | [bool](#bool) |  |  |
| error | [string](#string) |  |  |
| nonce | [int32](#int32) |  |  |






<a name="messages.ActionMsg"></a>

### ActionMsg



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| type | [ActionType](#messages.ActionType) |  |  |
| chips | [int32](#int32) |  |  |
| nonce | [int32](#int32) |  |  |






<a name="messages.ChatMsgRecv"></a>

### ChatMsgRecv



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| message_id | [int32](#int32) |  |  |
| player_id | [string](#string) |  |  |
| data | [string](#string) |  |  |
| timestamp | [int32](#int32) |  |  |






<a name="messages.ChatMsgSend"></a>

### ChatMsgSend



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| data | [string](#string) |  |  |






<a name="messages.Packet"></a>

### Packet



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| event | [EventType](#messages.EventType) |  |  |
| sit | [SitEvent](#messages.SitEvent) |  |  |
| stand | [StandEvent](#messages.StandEvent) |  |  |
| sit_ack | [SitAck](#messages.SitAck) |  |  |
| stand_ack | [StandAck](#messages.StandAck) |  |  |
| join_state | [table.TableState](#table.TableState) |  |  |
| hand | [table.CardSet](#table.CardSet) |  |  |
| flop | [table.CardSet](#table.CardSet) |  |  |
| turn | [table.CardSet](#table.CardSet) |  |  |
| river | [table.CardSet](#table.CardSet) |  |  |
| action | [ActionMsg](#messages.ActionMsg) |  |  |
| action_ack | [ActionAck](#messages.ActionAck) |  |  |
| msg_send | [ChatMsgSend](#messages.ChatMsgSend) |  |  |
| msg_recv | [ChatMsgRecv](#messages.ChatMsgRecv) |  |  |






<a name="messages.PlayerSeat"></a>

### PlayerSeat



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| player | [string](#string) |  |  |
| balance | [int32](#int32) |  |  |
| seat_num | [int32](#int32) |  |  |






<a name="messages.SitAck"></a>

### SitAck



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| table_id | [string](#string) |  |  |
| sat_down | [bool](#bool) |  |  |
| seat_num | [int32](#int32) |  |  |
| reason | [string](#string) |  |  |






<a name="messages.SitEvent"></a>

### SitEvent



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| table_id | [string](#string) |  |  |
| user_identity_token | [string](#string) |  |  |
| chips | [int32](#int32) |  |  |






<a name="messages.StandAck"></a>

### StandAck



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| table_id | [string](#string) |  |  |
| stood_up | [bool](#bool) |  |  |
| balance | [int32](#int32) |  |  |
| reason | [string](#string) |  |  |






<a name="messages.StandEvent"></a>

### StandEvent



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| table_id | [string](#string) |  |  |





 


<a name="messages.ActionType"></a>

### ActionType


| Name | Number | Description |
| ---- | ------ | ----------- |
| UNKNOWN | 0 |  |
| CHECK | 1 |  |
| CALL | 2 |  |
| BET | 3 |  |
| RAISE | 4 |  |
| ALL_IN | 5 |  |



<a name="messages.EventType"></a>

### EventType


| Name | Number | Description |
| ---- | ------ | ----------- |
| UNDEFINED | 0 |  |
| JOIN_ROOM | 32 | Event used when we&#39;re join the room, for chat messages |
| TABLE_STATE | 33 |  |
| TABLE_SIT | 34 |  |
| TABLE_STAND | 35 |  |
| TABLE_SIT_ACK | 36 |  |
| TABLE_STAND_ACK | 37 |  |
| START_GAME | 38 |  |
| HAND | 51 |  |
| FLOP | 52 |  |
| TURN | 53 |  |
| RIVER | 54 |  |
| TIMER | 64 |  |
| ACTION | 65 |  |
| ACTION_ACK | 66 |  |
| CHAT_MSG_SEND | 128 |  |
| CHAT_MSG_RECV | 129 |  |


 

 

 



<a name="table.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## table.proto



<a name="table.Card"></a>

### Card



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| suit | [string](#string) |  | The suit of the card oneof, &#39;s&#39;, &#39;h&#39;, &#39;d&#39;, &#39;c&#39; |
| rank | [int32](#int32) |  | The rank of the card [0, 12) |






<a name="table.CardSet"></a>

### CardSet
Created due to the inability to have repetitions in one ofs.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| cards | [Card](#table.Card) | repeated |  |






<a name="table.Player"></a>

### Player



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| name | [string](#string) |  | The public id of the player. |
| chips | [int32](#int32) |  | The number of chips this player has currently has. |
| seat_num | [int32](#int32) |  | The position of seat this player is currently occupying, zero indexed. |






<a name="table.TableOptions"></a>

### TableOptions



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| name | [string](#string) |  | The name given to the table, must be unique. |
| owner | [string](#string) |  | The owner of this table, other than the admin this is the only player allowed to destroy this table. |
| min_buy | [int32](#int32) |  | The minimum number of chips required to join this table. |
| max_seats | [int32](#int32) |  | The maximum seats allowed at this table. |
| big_blind | [int32](#int32) |  | The big blind value, the small blind will be floor(big_blind/2). |
| seat_shuffle | [bool](#bool) |  | If this is true the seats will be shuffled after the dealer is determined at the start of the game. |
| seat_shuffle_rounds | [int32](#int32) |  | If this is non-zero the seats will be shuffled after the specified number of rounds have passed. |
| rate | [int32](#int32) |  | Conversion to CAD; (rate chips) = $0.01 CAD. |






<a name="table.TableState"></a>

### TableState



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| options | [TableOptions](#table.TableOptions) |  |  |
| players | [Player](#table.Player) | repeated |  |
| board_cards | [CardSet](#table.CardSet) |  |  |
| pot_size | [int32](#int32) |  |  |





 

 

 

 



## Scalar Value Types

| .proto Type | Notes | C++ | Java | Python | Go | C# | PHP | Ruby |
| ----------- | ----- | --- | ---- | ------ | -- | -- | --- | ---- |
| <a name="double" /> double |  | double | double | float | float64 | double | float | Float |
| <a name="float" /> float |  | float | float | float | float32 | float | float | Float |
| <a name="int32" /> int32 | Uses variable-length encoding. Inefficient for encoding negative numbers – if your field is likely to have negative values, use sint32 instead. | int32 | int | int | int32 | int | integer | Bignum or Fixnum (as required) |
| <a name="int64" /> int64 | Uses variable-length encoding. Inefficient for encoding negative numbers – if your field is likely to have negative values, use sint64 instead. | int64 | long | int/long | int64 | long | integer/string | Bignum |
| <a name="uint32" /> uint32 | Uses variable-length encoding. | uint32 | int | int/long | uint32 | uint | integer | Bignum or Fixnum (as required) |
| <a name="uint64" /> uint64 | Uses variable-length encoding. | uint64 | long | int/long | uint64 | ulong | integer/string | Bignum or Fixnum (as required) |
| <a name="sint32" /> sint32 | Uses variable-length encoding. Signed int value. These more efficiently encode negative numbers than regular int32s. | int32 | int | int | int32 | int | integer | Bignum or Fixnum (as required) |
| <a name="sint64" /> sint64 | Uses variable-length encoding. Signed int value. These more efficiently encode negative numbers than regular int64s. | int64 | long | int/long | int64 | long | integer/string | Bignum |
| <a name="fixed32" /> fixed32 | Always four bytes. More efficient than uint32 if values are often greater than 2^28. | uint32 | int | int | uint32 | uint | integer | Bignum or Fixnum (as required) |
| <a name="fixed64" /> fixed64 | Always eight bytes. More efficient than uint64 if values are often greater than 2^56. | uint64 | long | int/long | uint64 | ulong | integer/string | Bignum |
| <a name="sfixed32" /> sfixed32 | Always four bytes. | int32 | int | int | int32 | int | integer | Bignum or Fixnum (as required) |
| <a name="sfixed64" /> sfixed64 | Always eight bytes. | int64 | long | int/long | int64 | long | integer/string | Bignum |
| <a name="bool" /> bool |  | bool | boolean | boolean | bool | bool | boolean | TrueClass/FalseClass |
| <a name="string" /> string | A string must always contain UTF-8 encoded or 7-bit ASCII text. | string | String | str/unicode | string | string | string | String (UTF-8) |
| <a name="bytes" /> bytes | May contain any arbitrary sequence of bytes. | string | ByteString | str | []byte | ByteString | string | String (ASCII-8BIT) |

