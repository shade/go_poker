# Protocol Documentation
<a name="top"></a>

## Table of Contents

- [internal/proto/messages.proto](#internal/proto/messages.proto)
    - [ActionAck](#.ActionAck)
    - [ActionMsg](#.ActionMsg)
    - [Card](#.Card)
    - [CardSet](#.CardSet)
    - [ChatMsgRecv](#.ChatMsgRecv)
    - [ChatMsgSend](#.ChatMsgSend)
    - [Packet](#.Packet)
    - [PlayerSeat](#.PlayerSeat)
    - [SitAck](#.SitAck)
    - [SitEvent](#.SitEvent)
    - [StandAck](#.StandAck)
    - [StandEvent](#.StandEvent)
    - [TableState](#.TableState)
  
    - [ActionType](#.ActionType)
    - [EventType](#.EventType)
  
- [Scalar Value Types](#scalar-value-types)



<a name="internal/proto/messages.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## internal/proto/messages.proto



<a name=".ActionAck"></a>

### ActionAck



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| ok | [bool](#bool) |  |  |
| error | [string](#string) |  |  |






<a name=".ActionMsg"></a>

### ActionMsg



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| type | [ActionType](#ActionType) |  |  |
| chips | [int32](#int32) |  |  |






<a name=".Card"></a>

### Card



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| suit | [int32](#int32) |  |  |
| rank | [int32](#int32) |  |  |
| display | [string](#string) |  |  |






<a name=".CardSet"></a>

### CardSet



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| cards | [Card](#Card) | repeated |  |






<a name=".ChatMsgRecv"></a>

### ChatMsgRecv



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| message_id | [int32](#int32) |  |  |
| player_id | [string](#string) |  |  |
| data | [string](#string) |  |  |
| timestamp | [int32](#int32) |  |  |






<a name=".ChatMsgSend"></a>

### ChatMsgSend



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| data | [string](#string) |  |  |






<a name=".Packet"></a>

### Packet



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| event | [EventType](#EventType) |  |  |
| sit | [SitEvent](#SitEvent) |  |  |
| stand | [StandEvent](#StandEvent) |  |  |
| sit_ack | [SitAck](#SitAck) |  |  |
| stand_ack | [StandAck](#StandAck) |  |  |
| join_state | [TableState](#TableState) |  |  |
| hand | [CardSet](#CardSet) |  |  |
| flop | [CardSet](#CardSet) |  |  |
| turn | [CardSet](#CardSet) |  |  |
| river | [CardSet](#CardSet) |  |  |
| action | [ActionMsg](#ActionMsg) |  |  |
| action_ack | [ActionAck](#ActionAck) |  |  |
| msg_send | [ChatMsgSend](#ChatMsgSend) |  |  |
| msg_recv | [ChatMsgRecv](#ChatMsgRecv) |  |  |






<a name=".PlayerSeat"></a>

### PlayerSeat



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| player | [string](#string) |  |  |
| balance | [int32](#int32) |  |  |
| seat_num | [int32](#int32) |  |  |






<a name=".SitAck"></a>

### SitAck



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| table_id | [string](#string) |  |  |
| sat_down | [bool](#bool) |  |  |
| seat_num | [int32](#int32) |  |  |
| reason | [string](#string) |  |  |






<a name=".SitEvent"></a>

### SitEvent



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| table_id | [string](#string) |  |  |
| user_identity_token | [string](#string) |  |  |
| chips | [int32](#int32) |  |  |






<a name=".StandAck"></a>

### StandAck



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| table_id | [string](#string) |  |  |
| stood_up | [bool](#bool) |  |  |
| balance | [int32](#int32) |  |  |
| reason | [string](#string) |  |  |






<a name=".StandEvent"></a>

### StandEvent



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| table_id | [string](#string) |  |  |






<a name=".TableState"></a>

### TableState



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| table_id | [string](#string) |  |  |
| min_buy | [int32](#int32) |  |  |
| max_seats | [int32](#int32) |  |  |
| big_blind | [int32](#int32) |  |  |
| seats | [PlayerSeat](#PlayerSeat) | repeated |  |





 


<a name=".ActionType"></a>

### ActionType


| Name | Number | Description |
| ---- | ------ | ----------- |
| UNKNOWN | 0 |  |
| CHECK | 1 |  |
| CALL | 2 |  |
| BET | 3 |  |
| RAISE | 4 |  |
| ALL_IN | 5 |  |



<a name=".EventType"></a>

### EventType


| Name | Number | Description |
| ---- | ------ | ----------- |
| UNDEFINED | 0 |  |
| JOIN_TABLE | 32 |  |
| TABLE_STATE | 33 |  |
| TABLE_SIT | 34 |  |
| TABLE_STAND | 35 |  |
| TABLE_SIT_ACK | 36 |  |
| TABLE_STAND_ACK | 37 |  |
| HAND | 51 |  |
| FLOP | 52 |  |
| TURN | 53 |  |
| RIVER | 54 |  |
| TIMER | 64 |  |
| ACTION | 65 |  |
| ACTION_ACK | 66 |  |
| CHAT_MSG_SEND | 128 |  |
| CHAT_MSG_RECV | 129 |  |


 

 

 



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

