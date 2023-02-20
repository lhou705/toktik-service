package common

import (
	"toktik/service/web/kitex_gen/message"
	"toktik/service/web/kitex_gen/user"
	"toktik/service/web/kitex_gen/video"
)

type BaseResponse struct {
	StatusCode int32  `json:"status_code,omitempty"`
	StatusMsg  string `json:"status_msg,omitempty"`
}

type LoginAndRegisterResponse struct {
	StatusCode int32  `json:"status_code,omitempty" `
	StatusMsg  string `json:"status_msg,omitempty" `
	Token      string `json:"token"`
	UserId     int64  `json:"user_id"`
}

type GetUserInfoResponse struct {
	StatusCode int32                         `json:"status_code,omitempty"`
	StatusMsg  string                        `json:"status_msg,omitempty"`
	User       *user.GetUserInfoByUserIdResp `json:"user"`
}

type GetUserListResponse struct {
	StatusCode int32                           `json:"status_code,omitempty"`
	StatusMsg  string                          `json:"status_msg,omitempty"`
	User       []*user.GetUserInfoByUserIdResp `json:"user_list"`
}

type GetFriendListResponse struct {
	StatusCode int32              `json:"status_code,omitempty"`
	StatusMsg  string             `json:"status_msg,omitempty"`
	User       []*user.FriendUser `json:"user_list"`
}

type GetMessageListResponse struct {
	StatusCode  int32                  `json:"status_code,omitempty"`
	StatusMsg   string                 `json:"status_msg,omitempty"`
	MessageList []*message.MessageItem `json:"message_list"`
}

type GetFeedListResponse struct {
	StatusCode int32              `json:"status_code,omitempty"`
	StatusMsg  string             `json:"status_msg,omitempty"`
	VideoList  []*video.VideoItem `json:"video_list"`
	NextTime   int64              `json:"next_time"`
}

type GetVideoList struct {
	StatusCode int32              `json:"status_code,omitempty"`
	StatusMsg  string             `json:"status_msg,omitempty"`
	VideoList  []*video.VideoItem `json:"video_list"`
}

type GetComment struct {
	StatusCode int32          `json:"status_code,omitempty"`
	StatusMsg  string         `json:"status_msg,omitempty"`
	Comment    *video.Comment `json:"comment"`
}

type GetCommentList struct {
	StatusCode  int32            `json:"status_code,omitempty"`
	StatusMsg   string           `json:"status_msg,omitempty"`
	CommentList []*video.Comment `json:"comment_list"`
}
