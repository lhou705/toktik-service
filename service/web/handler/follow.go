package handler

import (
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"net/http"
	"strconv"
	"toktik/service/web/client"
	"toktik/service/web/common"
	"toktik/service/web/kitex_gen/user"
)

type FollowHandler struct {
}

func (f *FollowHandler) SetFollowStatus(ctx context.Context, c *app.RequestContext) {
	toUserId, err := strconv.ParseInt(c.Query("to_user_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusOK, common.BaseResponse{
			StatusCode: common.InvalidFiledType,
			StatusMsg:  "to_user_id的字段类型错误，应该为int型",
		})
		hlog.CtxErrorf(ctx, "字段to_user_id验证错误，原因：%v", err)
		return
	}
	actionType, err := strconv.ParseInt(c.Query("action_type"), 10, 64)
	if err != nil {
		c.JSON(http.StatusOK, common.BaseResponse{
			StatusCode: common.InvalidFiledType,
			StatusMsg:  "action_type的字段类型错误，应该为int型",
		})
		hlog.CtxErrorf(ctx, "字段action_type验证错误，原因：%v", err)
		return
	}
	var resp *user.FollowResp
	if actionType == 1 {
		resp, err = client.UserClient.FollowStatus(ctx, &user.FollowReq{
			FollowId:   toUserId,
			FollowerId: c.GetInt64("id"),
		})
	} else if actionType == 2 {
		resp, err = client.UserClient.UnFollowStatus(ctx, &user.FollowReq{
			FollowId:   toUserId,
			FollowerId: c.GetInt64("id"),
		})
	} else {
		c.JSON(http.StatusOK, common.BaseResponse{
			StatusCode: common.InvalidFiledValue,
			StatusMsg:  "字段action_type的值只能为1或者2",
		})
		return
	}
	if resp.GetIsSuccess() && actionType == 1 {
		c.JSON(http.StatusOK, common.BaseResponse{
			StatusCode: common.Success,
			StatusMsg:  "关注成功",
		})
		return
	}
	if resp.GetIsSuccess() && actionType == 2 {
		c.JSON(http.StatusOK, common.BaseResponse{
			StatusCode: common.Success,
			StatusMsg:  "取消成功",
		})
		return
	}
	c.JSON(http.StatusOK, common.BaseResponse{
		StatusCode: common.ReqError,
		StatusMsg:  "请求失败",
	})
}

func (f *FollowHandler) GetFollowList(ctx context.Context, c *app.RequestContext) {
	userId, err := strconv.ParseInt(c.Query("user_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusOK, common.BaseResponse{
			StatusCode: common.InvalidFiledType,
			StatusMsg:  "user_id的字段类型错误，应该为int型",
		})
		hlog.CtxErrorf(ctx, "字段user_id验证错误，原因：%v", err)
		return
	}
	resp, err := client.UserClient.GetFollowList(ctx, &user.GetFollowListReq{UserId: userId})
	if err != nil {
		c.JSON(http.StatusOK, common.BaseResponse{
			StatusCode: common.ReqError,
			StatusMsg:  "请求用户关注信息失败",
		})
		hlog.CtxErrorf(ctx, "请求用户关注信息错误，原因：%v", err)
		return
	}
	c.JSON(http.StatusOK, common.GetUserListResponse{
		StatusCode: common.Success,
		StatusMsg:  "获取用户信息成功",
		User:       resp.GetUserList(),
	})
}

func (f *FollowHandler) GetFollowerList(ctx context.Context, c *app.RequestContext) {
	userId, err := strconv.ParseInt(c.Query("user_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusOK, common.BaseResponse{
			StatusCode: common.InvalidFiledType,
			StatusMsg:  "user_id的字段类型错误，应该为int型",
		})
		hlog.CtxErrorf(ctx, "字段user_id验证错误，原因：%v", err)
		return
	}
	resp, err := client.UserClient.GetFollowerList(ctx, &user.GetFollowerListReq{UserId: userId})
	if err != nil {
		c.JSON(http.StatusOK, common.BaseResponse{
			StatusCode: common.ReqError,
			StatusMsg:  "请求用户粉丝信息失败",
		})
		hlog.CtxErrorf(ctx, "请求用户粉丝信息错误，原因：%v", err)
		return
	}
	c.JSON(http.StatusOK, common.GetUserListResponse{
		StatusCode: common.Success,
		StatusMsg:  "获取成功",
		User:       resp.GetUserList(),
	})
}

func (f *FollowHandler) GetFriendList(ctx context.Context, c *app.RequestContext) {
	userId, err := strconv.ParseInt(c.Query("user_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusOK, common.BaseResponse{
			StatusCode: common.InvalidFiledType,
			StatusMsg:  "user_id的字段类型错误，应该为int型",
		})
		hlog.CtxErrorf(ctx, "字段user_id验证错误，原因：%v", err)
		return
	}
	resp, err := client.UserClient.GetFriendList(ctx, &user.GetFriendListReq{UserId: userId})
	if err != nil {
		c.JSON(http.StatusOK, common.BaseResponse{
			StatusCode: common.ReqError,
			StatusMsg:  "请求用户好友信息失败",
		})
		hlog.CtxErrorf(ctx, "请求用户好友信息错误，原因：%v", err)
		return
	}
	c.JSON(http.StatusOK, common.GetFriendListResponse{
		StatusCode: common.Success,
		StatusMsg:  "获取成功",
		User:       resp.GetUserList(),
	})
}
