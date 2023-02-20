namespace go message

service Message {
SendMessageResp SendMessage(1:SendMessageReq req)
GetMessageListResp GetMessageList(1:GetMessageListReq req)
}

struct SendMessageReq {
1:i64 from_user_id
2:i64 to_user_id
3:string content
}

struct SendMessageResp {
1:bool isSuccess
}

struct GetMessageListReq {
1:i64 from_user_id
2:i64 to_user_id
3:i64 pre_msg_time
}
struct GetMessageListResp{
1:list<MessageItem> message_list
}

struct MessageItem {
1:i64 id
2:i64 from_user_id
3:i64 to_user_id
4:string content
5:i64 create_time
}