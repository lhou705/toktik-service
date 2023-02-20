package main

import (
	"context"
	"errors"
	"github.com/acmestack/gorm-plus/gplus"
	"github.com/cloudwego/kitex/pkg/klog"
	"gorm.io/gorm"
	"toktik/service/user/kitex_gen/user"
)

// UserImpl implements the last service interface defined in the IDL.
type UserImpl struct{}

var selects = []string{"users.id as id",
	"users.name as name",
	"users.follow_count as follow_count",
	"users.follower_count as follower_count",
	"users.avatar as avatar",
	"users.background_image as background_image",
	"users.signature as signature",
	"users.total_favorited as total_favorited",
	"users.work_count as work_count",
	"users.favorite_count as favorite_count",
	"follows.is_follow as is_follow"}

// CheckUser implements the UserImpl interface.
func (s *UserImpl) CheckUser(ctx context.Context, req *user.CheckUserReq) (resp *user.CheckUserResp, err error) {
	query, model := gplus.NewQuery[User]()
	query.Eq(&model.Name, req.GetUsername()).Eq(&model.Password, req.Password).Select(&model.Name, &model.ID)
	result, db := gplus.SelectOne(query)
	if db.Error != nil {
		klog.CtxErrorf(ctx, "查找用户失败，原因%v", db.Error)
		return nil, db.Error
	}
	resp = &user.CheckUserResp{Username: result.Name, UserId: result.ID}
	return resp, nil
}

// CreateUser implements the UserImpl interface.
func (s *UserImpl) CreateUser(ctx context.Context, req *user.RegisterUserReq) (resp *user.RegisterUserResp, err error) {
	// 检查用户是否存在
	query, model := gplus.NewQuery[User]()
	query.Eq(&model.Name, req.GetUsername())
	_, db := gplus.SelectOne(query)
	if db.Error != nil && !errors.Is(db.Error, gorm.ErrRecordNotFound) {
		return nil, errors.New("此用户已存在")
	}
	db = gplus.Insert(&User{Name: req.GetUsername(), Password: req.GetPassword()})
	if db.Error != nil {
		klog.CtxErrorf(ctx, "创建用户失败，原因%v", db.Error)
		return nil, db.Error
	}
	query.Eq(&model.Name, req.GetUsername()).Eq(&model.Password, req.Password).Select(&model.Name, &model.ID)
	result, db := gplus.SelectOne(query)
	if db.Error != nil {
		klog.CtxErrorf(ctx, "查找用户失败，原因%v", db.Error)
		return nil, db.Error
	}
	resp = &user.RegisterUserResp{Username: result.Name, UserId: result.ID}
	return resp, nil
}

// GetUserInfoByUserId implements the UserImpl interface.
func (s *UserImpl) GetUserInfoByUserId(ctx context.Context, req *user.GetUserInfoByUserIdReq) (resp *user.GetUserInfoByUserIdResp, err error) {
	var followList user.GetUserInfoByUserIdResp
	err = Db.Model(&User{}).Select(selects).
		Joins("left join follows on users.id = follows.follower_id and follows.follow_id = ?", req.GetUserId()).
		Where("users.id = ? ", req.GetId()).
		Find(&followList).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			resp = nil
			return resp, nil
		}
		klog.CtxErrorf(ctx, "获取粉丝列表失败，原因：%v", err)
		resp = nil
		return resp, err
	}
	return &followList, nil
}

// GetUserInfoByUsername implements the UserImpl interface.
func (s *UserImpl) GetUserInfoByUsername(ctx context.Context, req *user.GetUserInfoByUsernameReq) (resp *user.GetUserInfoByUsernameResp, err error) {
	query, model := gplus.NewQuery[User]()
	query.Eq(&model.Name, req.GetUsername()).Select(&model.ID)
	result, db := gplus.SelectOne(query)
	if db.Error != nil {
		klog.CtxErrorf(ctx, "查找用户失败，原因%v", db.Error)
		return nil, db.Error
	}
	resp = &user.GetUserInfoByUsernameResp{
		Id: result.ID,
	}
	return resp, nil
}

// GetFollowList implements the UserImpl interface.
func (s *UserImpl) GetFollowList(ctx context.Context, req *user.GetFollowListReq) (resp *user.GetFollowListResp, err error) {
	var followList []*user.GetUserInfoByUserIdResp
	err = Db.Model(&Follow{}).Select(
		selects).
		Joins("inner join users on users.id = follows.follow_id").
		Where("follows.follower_id = ?", req.GetUserId()).
		Scan(&followList).Error
	resp = &user.GetFollowListResp{}
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			resp.UserList = nil
			return resp, nil
		}
		klog.CtxErrorf(ctx, "获取关注列表失败，原因：%v", err)
		resp.UserList = nil
		return resp, err
	}
	resp.UserList = followList
	return resp, nil
}

// GetFollowerList implements the UserImpl interface.
func (s *UserImpl) GetFollowerList(ctx context.Context, req *user.GetFollowerListReq) (resp *user.GetFollowerListResp, err error) {
	var followList []*user.GetUserInfoByUserIdResp
	err = Db.Model(&Follow{}).Select(
		selects).
		Joins("inner join users on users.id = follows.follower_id").
		Where("follows.follow_id = ?", req.GetUserId()).
		Scan(&followList).Error
	resp = &user.GetFollowerListResp{}
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			resp.UserList = nil
			return resp, nil
		}
		klog.CtxErrorf(ctx, "获取粉丝列表失败，原因：%v", err)
		resp.UserList = nil
		return resp, err
	}
	resp.UserList = followList
	return resp, nil
}

// GetFriendList implements the UserImpl interface.
func (s *UserImpl) GetFriendList(ctx context.Context, req *user.GetFriendListReq) (resp *user.GetFriendListResp, err error) {
	var followList []*user.FriendUser
	err = Db.Table("(?) as m, follows",
		Db.Model(&Message{}).
			Select("content as message,if(to_user_id = ?,0,1) as msgType", req.GetUserId()).
			Where("from_user_id = ? or to_user_id = ?", req.GetUserId(), req.GetUserId()).
			Order("created_at desc").
			Limit(1)).Select(
		selects,
		"m.msgType",
		"m.message").
		Joins("inner join users on users.id = follows.follow_id").
		Where("follows.follower_id = ? and follows.is_mutual = true", req.GetUserId()).
		Scan(&followList).
		Error
	resp = &user.GetFriendListResp{}
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			resp.UserList = nil
			return resp, nil
		}
		klog.CtxErrorf(ctx, "获取好友列表失败，原因：%v", err)
		resp.UserList = nil
		return resp, err
	}
	resp.UserList = followList
	return
}

// FollowStatus implements the UserImpl interface.
func (s *UserImpl) FollowStatus(ctx context.Context, req *user.FollowReq) (resp *user.FollowResp, err error) {
	err = Db.Where(&Follow{FollowId: req.GetFollowId(), FollowerId: req.GetFollowerId()}).Assign(map[string]any{
		"is_follow": true,
	}).FirstOrCreate(&Follow{FollowId: req.GetFollowId(), FollowerId: req.GetFollowerId(), IsFollow: true}).Error
	resp = &user.FollowResp{}
	if err != nil {
		klog.CtxErrorf(ctx, "修改或创建用户%d->用户%d的关注关系错误，原因：%v", req.GetFollowerId(), req.GetFollowId(), err)
		resp.IsSuccess = false
		return
	}
	err = Db.Model(&User{ID: req.GetFollowId()}).UpdateColumn("follower_count", gorm.Expr("follower_count + 1")).Error
	if err != nil {
		klog.CtxErrorf(ctx, "修改用户%d的粉丝数量错误，原因：%v", req.GetFollowId(), err)
		resp.IsSuccess = false
		return
	}
	err = Db.Model(&User{ID: req.GetFollowerId()}).UpdateColumn("follow_count", gorm.Expr("follow_count + 1")).Error
	if err != nil {
		klog.CtxErrorf(ctx, "修改用户%d的关注数量错误，原因：%v", req.GetFollowerId(), err)
		resp.IsSuccess = false
		return
	}
	var follow *Follow
	err = Db.Where(&Follow{FollowId: req.GetFollowerId(), FollowerId: req.GetFollowId()}).First(&follow).Error
	if err != nil {
		klog.CtxErrorf(ctx, "获取用户%d->用户%d的关注关系错误，原因：%v", req.GetFollowId(), req.GetFollowerId(), err)
		resp.IsSuccess = false
		return resp, err
	}
	if follow.IsFollow {
		err = Db.Where(&Follow{FollowId: req.GetFollowId(), FollowerId: req.GetFollowerId()}).
			Assign(map[string]any{
				"is_mutual": true,
			}).
			FirstOrCreate(&Follow{FollowId: req.GetFollowId(), FollowerId: req.GetFollowerId(), IsMutual: false}).Error
		if err != nil {
			klog.CtxErrorf(ctx, "修改或创建用户%d->用户%d的互关关系错误，原因：%v", req.GetFollowerId(), req.GetFollowId(), err)
			resp.IsSuccess = false
			return resp, err
		}
		err = Db.Where(&Follow{FollowId: req.GetFollowerId(), FollowerId: req.GetFollowId()}).Assign(map[string]any{
			"is_mutual": true,
		}).
			FirstOrCreate(&Follow{FollowId: req.GetFollowerId(), FollowerId: req.GetFollowId(), IsMutual: false}).Error
		if err != nil {
			klog.CtxErrorf(ctx, "修改或创建用户%d->用户%d的互关关系错误，原因：%v", req.GetFollowerId(), req.GetFollowId(), err)
			resp.IsSuccess = false
			return resp, err
		}
	}
	resp.IsSuccess = true
	return resp, nil
}

// UnFollowStatus implements the UserImpl interface.
func (s *UserImpl) UnFollowStatus(ctx context.Context, req *user.FollowReq) (resp *user.FollowResp, err error) {
	err = Db.Where(&Follow{FollowId: req.GetFollowId(), FollowerId: req.GetFollowerId()}).Assign(map[string]any{
		"is_follow": false,
	}).FirstOrCreate(&Follow{FollowId: req.GetFollowId(), FollowerId: req.GetFollowerId(), IsFollow: false}).Error
	resp = &user.FollowResp{}
	if err != nil {
		klog.CtxErrorf(ctx, "修改或创建用户%d->用户%d的关注关系错误，原因：%v", req.GetFollowerId(), req.GetFollowId(), err)
		resp.IsSuccess = false
		return
	}
	err = Db.Model(&User{ID: req.GetFollowId()}).UpdateColumn("follower_count", gorm.Expr("follower_count - 1")).Error
	if err != nil {
		klog.CtxErrorf(ctx, "修改用户%d的粉丝数量错误，原因：%v", req.GetFollowId(), err)
		resp.IsSuccess = false
		return
	}
	err = Db.Model(&User{ID: req.GetFollowerId()}).UpdateColumn("follow_count", gorm.Expr("follower_count - 1")).Error
	if err != nil {
		klog.CtxErrorf(ctx, "修改用户%d的关注数量错误，原因：%v", req.GetFollowerId(), err)
		resp.IsSuccess = false
		return
	}
	err = Db.Where(&Follow{FollowId: req.GetFollowId(), FollowerId: req.GetFollowerId()}).Assign(map[string]any{
		"is_mutual": false,
	}).FirstOrCreate(&Follow{FollowId: req.GetFollowId(), FollowerId: req.GetFollowerId(), IsMutual: false}).Error
	if err != nil {
		klog.CtxErrorf(ctx, "修改或创建用户%d->用户%d的互关关系错误，原因：%v", req.GetFollowerId(), req.GetFollowId(), err)
		resp.IsSuccess = false
		return resp, err
	}
	err = Db.Where(&Follow{FollowId: req.GetFollowerId(), FollowerId: req.GetFollowId()}).Assign(map[string]any{
		"is_mutual": false,
	}).FirstOrCreate(&Follow{FollowId: req.GetFollowerId(), FollowerId: req.GetFollowId(), IsMutual: false}).Error
	if err != nil {
		klog.CtxErrorf(ctx, "修改或创建用户%d->用户%d的互关关系错误，原因：%v", req.GetFollowerId(), req.GetFollowId(), err)
		resp.IsSuccess = false
		return resp, err
	}
	resp.IsSuccess = true
	return resp, nil
}
