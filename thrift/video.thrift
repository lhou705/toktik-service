namespace go video

service Video {
GetFeedListResp GetFeedList(1:GetFeedListReq req)
PublishVideoResp PublishVideo(1:PublishVideoReq req)
GetPublishListResp GetPublishList(1:GetPublishListReq req)
FavoriteVideoResp FavoriteVideoStatus(1:FavoriteVideoReq req)
FavoriteVideoResp UnFavoriteVideoStatus(1:FavoriteVideoReq req)
GetFavoriteListResp GetFavoriteList(1:GetFavoriteListReq req)
SendCommentResp SendComment(1:SendCommentReq req)
DeleteCommentResp DeleteComment(1:DeleteCommentReq req)
GetCommentListResp GetCommentListComment(1:GetCommentListReq req)
}

struct User {
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

struct VideoItem {
1:i64 id
2:User author
3:string play_url
4:string cover_url
5:i64 favorite_count
6:i64 comment_count
7:bool is_favorite
8:string title
}

struct GetFeedListResp {
1:list<VideoItem> video_list
2:i64 next_time
}

struct GetFeedListReq {
1:i64 last_time
2:i64 user_id
}

struct PublishVideoResp {
1:bool isSuccess
}

struct PublishVideoReq {
1:binary data
2:string title
3:i64 user_id
4:string file_name
}

struct GetPublishListResp {
1:list<VideoItem>video_list
}

struct GetPublishListReq {
1:i64 user_id
}

struct FavoriteVideoResp {
1:bool isSuccess
}

struct FavoriteVideoReq {
1:i64 user_id
2:i64 video_id
}

struct GetFavoriteListResp {
1:list<VideoItem>video_list
}

struct GetFavoriteListReq {
1:i64 user_id
}
struct Comment {
1:i64 id
2:User user
3:string content
4:string create_date
}
struct SendCommentResp {
1:Comment comment
}

struct SendCommentReq {
1:i64 user_id
2:i64 video_id
3:string content
}

struct DeleteCommentResp {
1:bool isSuccess
}

struct DeleteCommentReq {
1:i64 user_id
2:i64 video_id
3:i64 comment_id
}

struct GetCommentListResp {
1:list<Comment> comment_list
}

struct GetCommentListReq {
1:i64 user_id
2:i64 video_id
}