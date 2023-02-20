package handler

import (
	"context"
	"errors"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"net/http"
	"strconv"
	"toktik/service/web/client"
	"toktik/service/web/common"
	"toktik/service/web/kitex_gen/user"
	"toktik/service/web/utils"
)

type UserHandler struct{}

func (u *UserHandler) GetUserInfo(ctx context.Context, c *app.RequestContext) {
	userId, err := strconv.ParseInt(c.Query("user_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusOK, common.BaseResponse{
			StatusCode: common.InvalidFiledType,
			StatusMsg:  "user_id的字段类型错误，应该为int型",
		})
		hlog.CtxErrorf(ctx, "字段user_id验证错误，原因：%v", err)
		return
	}
	id := c.GetInt64("id")
	resp, err := client.UserClient.GetUserInfoByUserId(ctx, &user.GetUserInfoByUserIdReq{UserId: userId, Id: id})
	if err != nil {
		c.JSON(http.StatusOK, common.BaseResponse{
			StatusCode: common.ReqError,
			StatusMsg:  "请求用户信息失败",
		})
		hlog.CtxErrorf(ctx, "请求用户信息错误错误，原因：%v", err)
		return
	}
	c.JSON(http.StatusOK, common.GetUserInfoResponse{
		StatusCode: common.Success,
		StatusMsg:  "获取用户信息成功",
		User:       resp,
	})
}

func (u *UserHandler) Login(ctx context.Context, c *app.RequestContext) {
	username := c.Query("username")
	if len(username) <= 0 {
		c.JSON(http.StatusOK, common.BaseResponse{
			StatusCode: common.FiledIsNone,
			StatusMsg:  "username不能为空",
		})
		hlog.CtxErrorf(ctx, "字段username为空")
		return
	}
	password := c.Query("password")
	if len(password) <= 0 {
		c.JSON(http.StatusOK, common.BaseResponse{
			StatusCode: common.FiledLengthInvalid,
			StatusMsg:  "password长度不够，至少六位",
		})
		hlog.CtxErrorf(ctx, "字段password长度不够")
		return
	}
	// 验证用户
	res, err := client.UserClient.CheckUser(ctx, &user.CheckUserReq{
		Username: username,
		Password: password,
	})
	if err != nil {
		c.JSON(http.StatusOK, common.BaseResponse{
			StatusCode: common.ReqError,
			StatusMsg:  "登录失败，请检查您的输入",
		})
		hlog.CtxErrorf(ctx, "登陆验证失败，原因%v", err)
		return
	}
	if res.GetUsername() == "" || res.GetUserId() == 0 {
		c.JSON(http.StatusOK, common.BaseResponse{
			StatusCode: common.ReqError,
			StatusMsg:  "找不到此用户",
		})
		hlog.CtxInfof(ctx, "找不到用户")
		return
	}
	token, err := utils.GenerateToken(res.GetUsername())
	if err != nil {
		c.JSON(http.StatusOK, common.BaseResponse{
			StatusCode: common.ReqError,
			StatusMsg:  "登陆失败，请检查您的输入",
		})
		hlog.CtxErrorf(ctx, "生成token失败，原因%v", err)
		return
	}
	c.JSON(http.StatusOK, common.LoginAndRegisterResponse{
		StatusCode: common.Success,
		StatusMsg:  "登陆成功",
		Token:      token,
		UserId:     res.GetUserId(),
	})
}

func (u *UserHandler) Register(ctx context.Context, c *app.RequestContext) {
	username := c.Query("username")
	if len(username) <= 0 {
		c.JSON(http.StatusOK, common.BaseResponse{
			StatusCode: common.FiledIsNone,
			StatusMsg:  "username不能为空",
		})
		hlog.CtxErrorf(ctx, "字段username为空")
		return
	}
	password := c.Query("password")
	if len(password) <= 0 {
		c.JSON(http.StatusOK, common.BaseResponse{
			StatusCode: common.FiledLengthInvalid,
			StatusMsg:  "password长度不够，至少六位",
		})
		hlog.CtxErrorf(ctx, "字段password长度不够")
		return
	}
	res, err := client.UserClient.CreateUser(ctx, &user.RegisterUserReq{
		Username: username,
		Password: password,
	})
	if err != nil {
		if errors.Is(err, errors.New("此用户已存在")) {
			c.JSON(http.StatusOK, common.BaseResponse{
				StatusCode: common.ReqError,
				StatusMsg:  "注册失败，此用户已存在",
			})
			hlog.CtxErrorf(ctx, "注册验证失败，原因%v", err)
			return
		}
		c.JSON(http.StatusOK, common.BaseResponse{
			StatusCode: common.ReqError,
			StatusMsg:  "注册失败，请检查您的输入",
		})
		hlog.CtxErrorf(ctx, "注册验证失败，原因%v", err)
		return
	}
	if res.GetUsername() == "" || res.GetUserId() == 0 {
		c.JSON(http.StatusOK, common.BaseResponse{
			StatusCode: common.ReqError,
			StatusMsg:  "找不到此用户",
		})
		hlog.CtxInfof(ctx, "找不到用户")
		return
	}
	token, err := utils.GenerateToken(res.GetUsername())
	if err != nil {
		c.JSON(http.StatusOK, common.BaseResponse{
			StatusCode: common.ReqError,
			StatusMsg:  "注册失败，请检查您的输入",
		})
		hlog.CtxErrorf(ctx, "生成token失败，原因%v", err)
		return
	}
	c.JSON(http.StatusOK, common.LoginAndRegisterResponse{
		StatusCode: common.Success,
		StatusMsg:  "注册成功",
		Token:      token,
		UserId:     res.GetUserId(),
	})
}
