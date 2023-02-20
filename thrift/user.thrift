namespace go user

service User {
CheckUserResp CheckUser(1:CheckUserReq req)
RegisterUserResp CreateUser(1:RegisterUserReq req)
GetUserInfoByUsernameResp GetUserInfoByUsername(1:GetUserInfoByUsernameReq req)
GetUserInfoByUserIdResp GetUserInfoByUserId(1: GetUserInfoByUserIdReq req)
FollowResp FollowStatus(1:FollowReq req)
FollowResp UnFollowStatus(1:FollowReq req)
GetFollowListResp GetFollowList(1:GetFollowListReq req)
GetFollowerListResp GetFollowerList(1:GetFollowerListReq req)
GetFriendListResp GetFriendList(1:GetFriendListReq req)
}

struct CheckUserResp {
1:string username
2:i64 user_id
}

struct CheckUserReq {
1:string username
2:string password
}

struct RegisterUserResp {
1:string username
2:i64 user_id
}

struct RegisterUserReq {
1:string username
2:string password
}

struct GetUserInfoByUsernameResp {
1:i64 id
}

struct GetUserInfoByUserIdResp {
1:i64 id
2:string name
3:i64 follow_count
4:i64 follower_count
5:bool is_follow
6:string avatar
7:string background_image
8:string signature
9:i64 total_favorited
10:i64 work_count
11:i64 favorite_count
}

struct GetUserInfoByUsernameReq {
1:string username
}

struct GetUserInfoByUserIdReq {
1:i64 user_id
2:i64 id
}

struct FollowReq {
1:i64 follow_id
3:i64 follower_id
}

struct FollowResp {
1:bool is_success
}

struct GetFollowListResp {
1:list<GetUserInfoByUserIdResp> user_list
}

struct GetFollowListReq {
1:i64 user_id
}

struct GetFollowerListResp {
1:list<GetUserInfoByUserIdResp> user_list
}

struct GetFollowerListReq {
1:i64 user_id
}

struct GetFriendListReq {
1:i64 user_id
}

struct GetFriendListResp {
1:list<FriendUser> user_list
}

struct FriendUser {
1:i64 id
2:string name
3:i64 follow_count
4:i64 follower_count
5:bool is_follow
6:string avatar
7:string background_image
8:string signature
9:i64 total_favorited
10:i64 work_count
11:i64 favorite_count
12:string message
13:i64 msgType
}