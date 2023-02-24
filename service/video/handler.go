package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/cloudwego/kitex/pkg/klog"
	"gorm.io/gorm"
	"path"
	"time"
	"toktik/service/video/kitex_gen/cos"
	"toktik/service/video/kitex_gen/video"
)

// VideoImpl implements the last service interface defined in the IDL.
type VideoImpl struct{}
type videoItem struct {
	ID            int64
	PlayUrl       string
	CoverUrl      string
	FavoriteCount int64
	CommentCount  int64
	IsFavorite    bool
	Title         string
}

type commentItem struct {
	ID         int64
	Content    string
	CreateDate string
}

func addUserToList(ctx context.Context, source []*videoItem, followerId int64) []*video.VideoItem {
	var cache = make(map[int64]*video.User)
	var dst []*video.VideoItem
	for _, item := range source {
		var dstItem = &video.VideoItem{
			Id:            item.ID,
			PlayUrl:       item.PlayUrl,
			CoverUrl:      item.CoverUrl,
			FavoriteCount: item.FavoriteCount,
			CommentCount:  item.CommentCount,
			IsFavorite:    item.IsFavorite,
			Title:         item.Title,
		}
		if cache[item.ID] == nil {
			// 说明缓存中没有，查找并添加到缓存中
			var user *video.User
			err := Db.Model(&User{}).Select(
				"users.id as id",
				"users.name as name",
				"users.follow_count as follow_count",
				"users.follower_count as follower_count",
				"users.avatar as avatar",
				"users.background_image as background_image",
				"users.signature as signature",
				"users.total_favorited as total_favorited",
				"users.work_count as work_count",
				"users.favorite_count as favorite_count",
				"follows.is_follow as is_follow").
				Joins("left join videos on videos.author_id = users.id").
				Joins("left join follows on follows.follower_id = ? and follows.follow_id = videos.author_id", followerId).
				Where("videos.id = ?", item.ID).
				Scan(&user).Error
			if err == nil || errors.Is(err, gorm.ErrRecordNotFound) {
				cache[item.ID] = user
				dstItem.Author = user
				dst = append(dst, dstItem)
				continue
			} else {
				klog.CtxErrorf(ctx, "查找视频%d的作者出错，原因：%v", item.ID, err)
			}
		} else {
			dstItem.Author = cache[item.ID]
			dst = append(dst, dstItem)
		}
	}
	return dst
}

func addUserToCommentList(ctx context.Context, source []*commentItem, followerId int64) []*video.Comment {
	var cache = make(map[int64]*video.User)
	var dst []*video.Comment
	for _, item := range source {
		var dstItem = &video.Comment{
			Id:         item.ID,
			Content:    item.Content,
			CreateDate: item.CreateDate,
		}
		if cache[item.ID] == nil {
			// 说明缓存中没有，查找并添加到缓存中
			var user *video.User
			err := Db.Model(&User{}).Select(
				"users.id as id",
				"users.name as name",
				"users.follow_count as follow_count",
				"users.follower_count as follower_count",
				"users.avatar as avatar",
				"users.background_image as background_image",
				"users.signature as signature",
				"users.total_favorited as total_favorited",
				"users.work_count as work_count",
				"users.favorite_count as favorite_count",
				"follows.is_follow as is_follow").
				Joins("left join comments on comments.user_id = users.id").
				Joins("left join follows on follows.follower_id = ? and follows.follow_id = comments.user_id", followerId).
				Where("comments.id = ?", item.ID).
				Scan(&user).Error
			if err == nil || errors.Is(err, gorm.ErrRecordNotFound) {
				cache[item.ID] = user
				dstItem.User = user
				dst = append(dst, dstItem)
				continue
			} else {
				klog.CtxErrorf(ctx, "查找视频%d的作者出错，原因：%v", item.ID, err)
			}
		} else {
			dstItem.User = cache[item.ID]
			dst = append(dst, dstItem)
		}
	}
	return dst
}

var selects = []string{"videos.id as id",
	"videos.play_url as play_url",
	"videos.cover_url as cover_url",
	"videos.favorite_count as favorite_count",
	"videos.comment_count as comment_count",
	"videos.title as title",
	"favorites.is_favorite as is_favorite"}

// GetFeedList implements the VideoImpl interface.
func (s *VideoImpl) GetFeedList(ctx context.Context, req *video.GetFeedListReq) (resp *video.GetFeedListResp, err error) {
	var result []*videoItem
	where := "videos.created_at < ? "
	err = Db.Model(&Video{}).Select(selects).
		Joins("left join favorites on favorites.video_id = videos.id and favorites.user_id = ?", req.GetUserId()).Where(
		where, req.GetLastTime()).Order("videos.created_at desc").Scan(&result).Error
	resp = &video.GetFeedListResp{}
	if err != nil {
		klog.CtxErrorf(ctx, "查找视频失败，原因：%v", err)
		return resp, err
	}
	res := addUserToList(ctx, result, req.GetUserId())

	if len(result) == 0 {
		resp.VideoList = nil
		resp.NextTime = 0
		return resp, nil
	}
	var v *Video
	err = Db.Model(&Video{ID: result[len(result)-1].ID}).Select("created_at").First(&v).Error
	if err != nil {
		klog.CtxErrorf(ctx, "查找视频失败，原因：%v", err)
		return resp, err
	}
	resp.VideoList = res
	resp.NextTime = v.CreatedAt
	return resp, nil
}

// PublishVideo implements the VideoImpl interface.
func (s *VideoImpl) PublishVideo(ctx context.Context, req *video.PublishVideoReq) (resp *video.PublishVideoResp, err error) {
	// 将数据发给cos

	playKey := fmt.Sprintf("play/%d%s", time.Now().Unix(), path.Ext(req.GetFileName()))
	coverKey := fmt.Sprintf("cover/%d.jpg", time.Now().Unix())
	go func() {
		res, err := cosClient.Upload(ctx, &cos.UploadReq{
			File: req.GetData(),
			Key:  playKey,
		})
		if err != nil || !res.GetIsSuccess() {
			klog.CtxErrorf(ctx, "上传用户%d视频%s失败，原因：%v", req.GetUserId(), req.GetTitle(), err)
		}
	}()
	klog.CtxInfof(ctx, "上传视频%s成功", playKey)
	// 创建记录
	err = Db.Create(&Video{
		AuthorId: req.GetUserId(),
		PlayUrl:  config.Cos.CdnAddr + playKey,
		CoverUrl: config.Cos.CdnAddr + coverKey,
		Title:    req.GetTitle(),
	}).Error
	if err != nil {
		resp.IsSuccess = false
		klog.CtxErrorf(ctx, "创建用户%d视频%s失败，原因：%v", req.GetUserId(), req.GetTitle(), err)
		return resp, err
	}
	err = Db.Model(&User{ID: req.GetUserId()}).UpdateColumn("work_count", gorm.Expr("work_count + 1")).Error
	if err != nil {
		resp.IsSuccess = false
		klog.CtxErrorf(ctx, "修改用户%d的作品数量错误，原因：%v", req.GetUserId(), err)
		return resp, err
	}
	resp.IsSuccess = true
	return resp, nil
}

// GetPublishList implements the VideoImpl interface.
func (s *VideoImpl) GetPublishList(ctx context.Context, req *video.GetPublishListReq) (resp *video.GetPublishListResp, err error) {
	var result []*videoItem
	// 先查视频
	err = Db.Model(&Video{}).Select(selects).
		Joins("left join favorites on favorites.video_id = videos.id and favorites.user_id= videos.author_id").Where(
		"videos.author_id = ?", req.GetUserId()).Scan(&result).Error
	resp = &video.GetPublishListResp{}
	if err != nil {
		klog.CtxErrorf(ctx, "查找视频失败，原因：%v", err)
		return resp, err
	}
	// 在查用户 -- 缓存加快速度
	res := addUserToList(ctx, result, req.GetUserId())
	// 组装
	resp.VideoList = res
	return resp, nil
}

// GetFavoriteList implements the VideoImpl interface.
func (s *VideoImpl) GetFavoriteList(ctx context.Context, req *video.GetFavoriteListReq) (resp *video.GetFavoriteListResp, err error) {
	var result []*videoItem
	// 先查视频
	err = Db.Model(&Video{}).Select(selects).
		Joins("left join favorites on favorites.video_id = videos.id and favorites.is_favorite = true").Where(
		"favorites.user_id = ?", req.GetUserId()).Scan(&result).Error
	resp = &video.GetFavoriteListResp{}
	if err != nil {
		klog.CtxErrorf(ctx, "查找视频失败，原因：%v", err)
		return resp, err
	}
	// 在查用户 -- 缓存加快速度
	res := addUserToList(ctx, result, req.GetUserId())
	// 组装
	resp.VideoList = res
	return resp, nil
}

// SendComment implements the VideoImpl interface.
func (s *VideoImpl) SendComment(ctx context.Context, req *video.SendCommentReq) (resp *video.SendCommentResp, err error) {
	// 评论表创建记录
	err = Db.Create(&Comment{UserId: req.GetUserId(), VideoId: req.GetVideoId(), Content: req.GetContent()}).Error
	resp = &video.SendCommentResp{}
	if err != nil {
		klog.CtxErrorf(ctx, "创建用户%d给视频%d的评论失败，原因：%v", req.GetUserId(), req.GetVideoId(), req.GetContent(), err)
		return resp, err
	}
	// 视频表+1
	err = Db.Model(&Video{ID: req.GetVideoId()}).UpdateColumn("comment_count", gorm.Expr("comment_count + 1")).Error
	if err != nil {
		klog.CtxErrorf(ctx, "修改视频%d的评论数量错误，原因：%v", req.GetVideoId(), err)
		return
	}
	// 查找最新创建的值
	var comment []*commentItem
	err = Db.Model(&Comment{}).Select("id, content,from_unixtime(created_at,'%m-%d') as create_date ").
		Where(&Comment{UserId: req.GetUserId(), VideoId: req.GetVideoId(), Content: req.GetContent()}).Order("created_at desc").Limit(1).Scan(&comment).Error
	if err != nil {
		klog.CtxErrorf(ctx, "查找用户%d在视频%d的评论错误，原因：%v", req.GetUserId(), req.GetVideoId(), err)
		return
	}
	res := addUserToCommentList(ctx, comment, req.GetUserId())
	resp.Comment = res[0]
	return resp, nil
}

// DeleteComment implements the VideoImpl interface.
func (s *VideoImpl) DeleteComment(ctx context.Context, req *video.DeleteCommentReq) (resp *video.DeleteCommentResp, err error) {
	err = Db.Delete(&Comment{}, &req.CommentId).Error
	resp = &video.DeleteCommentResp{}
	if err != nil {
		klog.CtxErrorf(ctx, "删除用户%d给视频%d的评论失败，原因：%v", req.GetUserId(), req.GetVideoId(), err)
		resp.IsSuccess = false
		return resp, err
	}
	err = Db.Model(&Video{ID: req.GetVideoId()}).UpdateColumn("comment_count", gorm.Expr("comment_count - 1")).Error
	if err != nil {
		klog.CtxErrorf(ctx, "修改视频%d的评论数量错误，原因：%v", req.GetVideoId(), err)
		return
	}
	resp.IsSuccess = true
	return resp, nil
}

// GetCommentListComment implements the VideoImpl interface.
func (s *VideoImpl) GetCommentListComment(ctx context.Context, req *video.GetCommentListReq) (resp *video.GetCommentListResp, err error) {
	var commentList []*commentItem
	resp = &video.GetCommentListResp{}
	err = Db.Model(&Comment{}).Select("id,content,from_unixtime(created_at,'%m-%d') as create_date").
		Where("video_id = ?", req.GetVideoId()).Order("created_at desc").Scan(&commentList).Error
	if err != nil {
		klog.CtxErrorf(ctx, "查找视频%d的评论失败，原因：%v", req.GetVideoId(), err)
		return resp, err
	}
	res := addUserToCommentList(ctx, commentList, req.GetUserId())
	resp.CommentList = res
	return resp, nil
}

// FavoriteVideoStatus implements the VideoImpl interface.
func (s *VideoImpl) FavoriteVideoStatus(ctx context.Context, req *video.FavoriteVideoReq) (resp *video.FavoriteVideoResp, err error) {
	err = Db.Where(&Favorite{VideoId: req.GetVideoId(), UserId: req.GetUserId()}).Assign(map[string]any{
		"is_favorite": true,
	}).FirstOrCreate(&Favorite{VideoId: req.GetVideoId(), UserId: req.GetUserId(), IsFavorite: true}).Error
	resp = &video.FavoriteVideoResp{}
	if err != nil {
		klog.CtxErrorf(ctx, "修改或创建用户%d->视频%d的点赞错误，原因：%v", req.GetUserId(), req.GetVideoId(), err)
		return resp, err
	}
	err = Db.Model(&Video{ID: req.GetVideoId()}).UpdateColumn("favorite_count", gorm.Expr("favorite_count + 1")).Error
	if err != nil {
		klog.CtxErrorf(ctx, "修改视频%d的点赞数错误，原因：%v", req.GetVideoId(), err)
		return resp, err
	}
	// 修改作者的总点赞数
	var v *Video
	err = Db.Model(&Video{ID: req.GetVideoId()}).Select("author_id").First(&v).Error
	if err != nil {
		klog.CtxErrorf(ctx, "查找视频%d的作者错误，原因：%v", req.GetVideoId(), err)
		return resp, err
	}
	err = Db.Model(&User{ID: v.AuthorId}).UpdateColumn("total_favorited", gorm.Expr("total_favorited + 1")).Error
	if err != nil {
		klog.CtxErrorf(ctx, "修改用户%d的作品总赞数错误，原因：%v", v.AuthorId, err)
		return resp, err
	}
	err = Db.Model(&User{ID: req.GetUserId()}).UpdateColumn("favorite_count", gorm.Expr("favorite_count + 1")).Error
	if err != nil {
		klog.CtxErrorf(ctx, "修改用户%d总赞数错误，原因：%v", v.AuthorId, err)
		return resp, err
	}
	resp.IsSuccess = true
	return resp, nil
}

// UnFavoriteVideoStatus implements the VideoImpl interface.
func (s *VideoImpl) UnFavoriteVideoStatus(ctx context.Context, req *video.FavoriteVideoReq) (resp *video.FavoriteVideoResp, err error) {
	err = Db.Where(&Favorite{VideoId: req.GetVideoId(), UserId: req.GetUserId()}).Assign(map[string]any{
		"is_favorite": false,
	}).FirstOrCreate(&Favorite{VideoId: req.GetVideoId(), UserId: req.GetUserId(), IsFavorite: false}).Error
	resp = &video.FavoriteVideoResp{}
	if err != nil {
		klog.CtxErrorf(ctx, "修改或创建用户%d->视频%d的点赞错误，原因：%v", req.GetUserId(), req.GetVideoId(), err)
		return resp, err
	}
	err = Db.Model(&Video{ID: req.GetVideoId()}).UpdateColumn("favorite_count", gorm.Expr("favorite_count - 1")).Error
	if err != nil {
		klog.CtxErrorf(ctx, "修改视频%d的点赞数错误，原因：%v", req.GetVideoId(), err)
		return resp, err
	}
	// 修改作者的总点赞数
	var v *Video
	err = Db.Model(&Video{ID: req.GetVideoId()}).Select("author_id").First(&v).Error
	if err != nil {
		klog.CtxErrorf(ctx, "查找视频%d的作者错误，原因：%v", req.GetVideoId(), err)
		return resp, err
	}
	err = Db.Model(&User{ID: v.AuthorId}).UpdateColumn("total_favorited", gorm.Expr("total_favorited - 1")).Error
	if err != nil {
		klog.CtxErrorf(ctx, "修改用户%d的作品总赞数错误，原因：%v", v.AuthorId, err)
		return resp, err
	}
	err = Db.Model(&User{ID: req.GetUserId()}).UpdateColumn("favorite_count", gorm.Expr("favorite_count - 1")).Error
	if err != nil {
		klog.CtxErrorf(ctx, "修改用户%d的作品总赞数错误，原因：%v", v.AuthorId, err)
		return resp, err
	}
	resp.IsSuccess = true
	return resp, nil
}
