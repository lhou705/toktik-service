package handler

import (
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"net/http"
	"strconv"
	"toktik/service/web/client"
	"toktik/service/web/common"
	"toktik/service/web/kitex_gen/message"
)

type MessageHandler struct {
}

func (m *MessageHandler) SendMessage(ctx context.Context, c *app.RequestContext) {
	actionType, err := strconv.ParseInt(c.Query("action_type"), 10, 64)
	if err != nil {
		c.JSON(http.StatusOK, common.BaseResponse{
			StatusCode: common.InvalidFiledType,
			StatusMsg:  "action_type的字段类型错误，应该为int型",
		})
		hlog.CtxErrorf(ctx, "字段action_type验证错误，原因：%v", err)
		return
	}
	if actionType != 1 {
		c.JSON(http.StatusOK, common.BaseResponse{
			StatusCode: common.InvalidFiledValue,
			StatusMsg:  "action_type的值只能为1",
		})
		hlog.CtxErrorf(ctx, "字段action_typ的值只能为1")
		return
	}
	toUserId, err := strconv.ParseInt(c.Query("to_user_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusOK, common.BaseResponse{
			StatusCode: common.InvalidFiledType,
			StatusMsg:  "to_user_id的字段类型错误，应该为int型",
		})
		hlog.CtxErrorf(ctx, "字段to_user_id验证错误，原因：%v", err)
		return
	}
	content := c.Query("content")
	userId := c.GetInt64("id")
	resp, err := client.MessageClient.SendMessage(ctx, &message.SendMessageReq{
		FromUserId: userId,
		ToUserId:   toUserId,
		Content:    content,
	})
	if err != nil {
		c.JSON(http.StatusOK, common.BaseResponse{
			StatusCode: common.ReqError,
			StatusMsg:  "发送消息失败",
		})
		hlog.CtxErrorf(ctx, "发送消息失败，原因：%v", err)
		return
	}
	if resp.GetIsSuccess() {
		c.JSON(http.StatusOK, common.BaseResponse{
			StatusCode: common.Success,
			StatusMsg:  "发送消息成功",
		})
		return
	}
	c.JSON(http.StatusOK, common.BaseResponse{
		StatusCode: common.ReqError,
		StatusMsg:  "发送消息失败",
	})
	return
}

func (m *MessageHandler) GetMessageList(ctx context.Context, c *app.RequestContext) {
	toUserId, err := strconv.ParseInt(c.Query("to_user_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusOK, common.BaseResponse{
			StatusCode: common.InvalidFiledType,
			StatusMsg:  "to_user_id的字段类型错误，应该为int型",
		})
		hlog.CtxErrorf(ctx, "字段to_user_id验证错误，原因：%v", err)
		return
	}
	preMsgTime, err := strconv.ParseInt(c.Query("pre_msg_time"), 10, 64)
	if err != nil {
		c.JSON(http.StatusOK, common.BaseResponse{
			StatusCode: common.InvalidFiledType,
			StatusMsg:  "pre_msg_time的字段类型错误，应该为int型",
		})
		hlog.CtxErrorf(ctx, "字段pre_msg_time验证错误，原因：%v", err)
		return
	}
	userId := c.GetInt64("id")
	resp, err := client.MessageClient.GetMessageList(ctx, &message.GetMessageListReq{
		FromUserId: userId,
		ToUserId:   toUserId,
		PreMsgTime: preMsgTime,
	})
	if err != nil {
		c.JSON(http.StatusOK, common.BaseResponse{
			StatusCode: common.ReqError,
			StatusMsg:  "获取聊天记录失败",
		})
		hlog.CtxErrorf(ctx, "获取聊天记录失败，原因：%v", err)
		return
	}
	c.JSON(http.StatusOK, common.GetMessageListResponse{
		StatusCode:  common.Success,
		StatusMsg:   "获取聊天记录成功",
		MessageList: resp.GetMessageList(),
	})
}
