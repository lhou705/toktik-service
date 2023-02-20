package handler

import (
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"io"
	"net/http"
	"strconv"
	"time"
	"toktik/service/web/client"
	"toktik/service/web/common"
	"toktik/service/web/kitex_gen/video"
)

type VideoHandler struct {
}

func (v *VideoHandler) GetFeedList(ctx context.Context, c *app.RequestContext) {
	lastTime, err := strconv.ParseInt(c.Query("last_time"), 10, 64)
	if err != nil {
		//
		lastTime = 0
	}
	if lastTime == 0 {
		lastTime = time.Now().Unix()
	}
	userId := c.GetInt64("id")
	resp, err := client.VideoClient.GetFeedList(ctx, &video.GetFeedListReq{LastTime: lastTime, UserId: userId})
	if err != nil {
		c.JSON(http.StatusOK, common.BaseResponse{
			StatusCode: common.ReqError,
			StatusMsg:  "请求视频错误",
		})
		hlog.CtxErrorf(ctx, "请求视频错误，原因：%v", err)
		return
	}
	c.JSON(http.StatusOK, common.GetFeedListResponse{
		StatusCode: common.Success,
		StatusMsg:  "获取成功",
		VideoList:  resp.GetVideoList(),
		NextTime:   resp.GetNextTime(),
	})
}

func (v *VideoHandler) GetPublishList(ctx context.Context, c *app.RequestContext) {
	userId, err := strconv.ParseInt(c.Query("user_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusOK, common.BaseResponse{
			StatusCode: common.InvalidFiledType,
			StatusMsg:  "user_id的字段类型错误，应该为int型",
		})
		hlog.CtxErrorf(ctx, "字段user_id验证错误，原因：%v", err)
		return
	}
	resp, err := client.VideoClient.GetPublishList(ctx, &video.GetPublishListReq{UserId: userId})
	if err != nil {
		c.JSON(http.StatusOK, common.BaseResponse{
			StatusCode: common.ReqError,
			StatusMsg:  "请求发布列表错误",
		})
		hlog.CtxErrorf(ctx, "请求发布列表错误，原因：%v", err)
		return
	}
	c.JSON(http.StatusOK, common.GetVideoList{
		StatusCode: common.Success,
		StatusMsg:  "请求成功",
		VideoList:  resp.GetVideoList(),
	})
}

func (v *VideoHandler) SetFavoriteStatus(ctx context.Context, c *app.RequestContext) {
	videoId, err := strconv.ParseInt(c.Query("video_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusOK, common.BaseResponse{
			StatusCode: common.InvalidFiledType,
			StatusMsg:  "video_id的字段类型错误，应该为int型",
		})
		hlog.CtxErrorf(ctx, "字段video_id验证错误，原因：%v", err)
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
	userId := c.GetInt64("id")
	if actionType == 1 {
		resp, err := client.VideoClient.FavoriteVideoStatus(ctx, &video.FavoriteVideoReq{
			UserId:  userId,
			VideoId: videoId,
		})
		if err != nil {
			c.JSON(http.StatusOK, common.BaseResponse{
				StatusCode: common.ReqError,
				StatusMsg:  "点赞时出现错误",
			})
			hlog.CtxErrorf(ctx, "点赞时出现错误，原因：%v", err)
			return
		}
		if resp.GetIsSuccess() {
			c.JSON(http.StatusOK, common.BaseResponse{
				StatusCode: common.Success,
				StatusMsg:  "点赞成功",
			})
			return
		}
	} else if actionType == 2 {
		resp, err := client.VideoClient.UnFavoriteVideoStatus(ctx, &video.FavoriteVideoReq{
			UserId:  userId,
			VideoId: videoId,
		})
		if err != nil {
			c.JSON(http.StatusOK, common.BaseResponse{
				StatusCode: common.ReqError,
				StatusMsg:  "取消点赞时出现错误",
			})
			hlog.CtxErrorf(ctx, "取消点赞时出现错误，原因：%v", err)
			return
		}
		if resp.GetIsSuccess() {
			c.JSON(http.StatusOK, common.BaseResponse{
				StatusCode: common.Success,
				StatusMsg:  "取消点赞成功",
			})
			return
		}
	} else {
		c.JSON(http.StatusOK, common.BaseResponse{
			StatusCode: common.InvalidFiledValue,
			StatusMsg:  "字段action_type的值只能为1或者2",
		})
		return
	}
}

func (v *VideoHandler) GetFavoriteList(ctx context.Context, c *app.RequestContext) {
	userId, err := strconv.ParseInt(c.Query("user_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusOK, common.BaseResponse{
			StatusCode: common.InvalidFiledType,
			StatusMsg:  "user_id的字段类型错误，应该为int型",
		})
		hlog.CtxErrorf(ctx, "字段user_id验证错误，原因：%v", err)
		return
	}
	resp, err := client.VideoClient.GetFavoriteList(ctx, &video.GetFavoriteListReq{UserId: userId})
	if err != nil {
		c.JSON(http.StatusOK, common.BaseResponse{
			StatusCode: common.ReqError,
			StatusMsg:  "请求点赞列表错误",
		})
		hlog.CtxErrorf(ctx, "请求点赞列表错误，原因：%v", err)
		return
	}
	c.JSON(http.StatusOK, common.GetVideoList{
		StatusCode: common.Success,
		StatusMsg:  "获取成功",
		VideoList:  resp.GetVideoList(),
	})
}

func (v *VideoHandler) SendComment(ctx context.Context, c *app.RequestContext) {
	videoId, err := strconv.ParseInt(c.Query("video_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusOK, common.BaseResponse{
			StatusCode: common.InvalidFiledType,
			StatusMsg:  "video_id的字段类型错误，应该为int型",
		})
		hlog.CtxErrorf(ctx, "字段video_id验证错误，原因：%v", err)
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
	userId := c.GetInt64("id")
	if actionType == 1 {
		commentText := c.Query("comment_text")
		resp, err := client.VideoClient.SendComment(ctx, &video.SendCommentReq{
			UserId:  userId,
			VideoId: videoId,
			Content: commentText,
		})
		if err != nil {
			c.JSON(http.StatusOK, common.BaseResponse{
				StatusCode: common.ReqError,
				StatusMsg:  "发布评论失败",
			})
			hlog.CtxErrorf(ctx, "发布评论失败", err)
			return
		}
		c.JSON(http.StatusOK, common.GetComment{
			StatusCode: common.Success,
			StatusMsg:  "获取成功",
			Comment:    resp.GetComment(),
		})
	} else if actionType == 2 {
		commentId, err := strconv.ParseInt(c.Query("comment_id"), 10, 64)
		if err != nil {
			c.JSON(http.StatusOK, common.BaseResponse{
				StatusCode: common.InvalidFiledType,
				StatusMsg:  "comment_id的字段类型错误，应该为int型",
			})
			hlog.CtxErrorf(ctx, "字段comment_id验证错误，原因：%v", err)
			return
		}
		resp, err := client.VideoClient.DeleteComment(ctx, &video.DeleteCommentReq{
			UserId:    userId,
			VideoId:   videoId,
			CommentId: commentId,
		})
		if err != nil {
			c.JSON(http.StatusOK, common.BaseResponse{
				StatusCode: common.ReqError,
				StatusMsg:  "删除评论时出现错误",
			})
			hlog.CtxErrorf(ctx, "删除评论时出现错误", err)
			return
		}
		if resp.GetIsSuccess() {
			c.JSON(http.StatusOK, common.BaseResponse{
				StatusCode: common.Success,
				StatusMsg:  "删除评论成功",
			})
		}
		c.JSON(http.StatusOK, common.BaseResponse{
			StatusCode: common.ReqError,
			StatusMsg:  "删除评论失败",
		})
	}
}

func (v *VideoHandler) GetCommentList(ctx context.Context, c *app.RequestContext) {
	videoId, err := strconv.ParseInt(c.Query("video_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusOK, common.BaseResponse{
			StatusCode: common.InvalidFiledType,
			StatusMsg:  "video_id的字段类型错误，应该为int型",
		})
		hlog.CtxErrorf(ctx, "字段video_id验证错误，原因：%v", err)
		return
	}
	userId := c.GetInt64("id")
	resp, err := client.VideoClient.GetCommentListComment(ctx, &video.GetCommentListReq{
		UserId:  userId,
		VideoId: videoId,
	})
	if err != nil {
		c.JSON(http.StatusOK, common.BaseResponse{
			StatusCode: common.ReqError,
			StatusMsg:  "请求评论列表列表错误",
		})
		hlog.CtxErrorf(ctx, "请求评论列表错误，原因：%v", err)
		return
	}
	c.JSON(http.StatusOK, common.GetCommentList{
		StatusCode:  http.StatusOK,
		StatusMsg:   "获取成功",
		CommentList: resp.GetCommentList(),
	})
}

func (v *VideoHandler) Publish(ctx context.Context, c *app.RequestContext) {
	userId := c.GetInt64("id")
	title := c.PostForm("title")
	data, err := c.FormFile("data")
	if err != nil {
		c.JSON(http.StatusOK, common.BaseResponse{
			StatusCode: common.InvalidFiledType,
			StatusMsg:  "获取文件错误",
		})
		hlog.CtxErrorf(ctx, "获取用户%d投递的文件%s错误%v", userId, title, err)
		return
	}
	fileContent, err := data.Open()
	if err != nil {
		c.JSON(http.StatusOK, common.BaseResponse{
			StatusCode: common.InvalidFiledType,
			StatusMsg:  "读取文件错误",
		})
		hlog.CtxErrorf(ctx, "读取用户%d投递的文件%s错误%v", userId, title, err)
		return
	}
	byteContent, err := io.ReadAll(fileContent)
	if err != nil {
		c.JSON(http.StatusOK, common.BaseResponse{
			StatusCode: common.InvalidFiledType,
			StatusMsg:  "解析文件错误",
		})
		hlog.CtxErrorf(ctx, "解析用户%d投递的文件%s错误%v", userId, title, err)
		return
	}
	resp, err := client.VideoClient.PublishVideo(ctx, &video.PublishVideoReq{
		Data:     byteContent,
		Title:    title,
		UserId:   userId,
		FileName: data.Filename,
	})
	if err != nil || !resp.GetIsSuccess() {
		c.JSON(http.StatusOK, common.BaseResponse{
			StatusCode: common.InvalidFiledType,
			StatusMsg:  "上传文件错误",
		})
		hlog.CtxErrorf(ctx, "上传用户%d投递的文件%s错误%v", userId, title, err)
		return
	}
	c.JSON(http.StatusOK, common.BaseResponse{
		StatusCode: common.Success,
		StatusMsg:  "上传成功",
	})
}
