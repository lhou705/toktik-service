package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/acmestack/gorm-plus/gplus"
	"github.com/cloudwego/kitex/pkg/klog"
	"gorm.io/gorm"
	"time"
	message "toktik/service/message/kitex_gen/message"
)

// MessageImpl implements the last service interface defined in the IDL.
type MessageImpl struct{}

// SendMessage implements the MessageImpl interface.
func (s *MessageImpl) SendMessage(ctx context.Context, req *message.SendMessageReq) (resp *message.SendMessageResp, err error) {
	db := gplus.Insert(&Message{
		FromUserId: req.GetFromUserId(),
		ToUserId:   req.GetToUserId(),
		Content:    req.GetContent(),
	})
	resp = &message.SendMessageResp{}
	if db.Error != nil {
		resp.IsSuccess = false
		klog.CtxErrorf(ctx, "用户%d发给%d的消息失败，原因：%v", req.GetFromUserId(), req.GetToUserId(), db.Error)
		return resp, db.Error
	}
	resp.IsSuccess = true
	return resp, nil
}

// GetMessageList implements the MessageImpl interface.
func (s *MessageImpl) GetMessageList(ctx context.Context, req *message.GetMessageListReq) (resp *message.GetMessageListResp, err error) {
	var result []*message.MessageItem
	err = Db.Model(&Message{}).
		Select("id", "to_user_id", "from_user_id", "content", "created_at as create_time").
		Where("((from_user_id = ? and to_user_id = ?) or (from_user_id = ? and to_user_id = ?)) and (created_at > ? and created_at <= ? )",
			req.GetFromUserId(), req.GetToUserId(), req.GetToUserId(), req.GetFromUserId(), req.GetPreMsgTime(), time.Now().Unix()*1000).
		Find(&result).
		Error
	resp = &message.GetMessageListResp{}
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		errStr := fmt.Sprintf("查找用户%d到%d的聊天记录失败。原因：%v", req.GetFromUserId(), req.GetToUserId(), err)
		klog.CtxErrorf(ctx, errStr)
		resp.MessageList = nil
		return resp, err
	}
	resp.MessageList = result
	return resp, nil
}
