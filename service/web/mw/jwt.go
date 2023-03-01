package mw

import (
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"net/http"
	"strings"
	"toktik/service/web/client"
	"toktik/service/web/common"
	"toktik/service/web/kitex_gen/user"
	"toktik/service/web/utils"
)

func JWT(ctx context.Context, c *app.RequestContext) {
	token := c.Query("token")
	if strings.Contains(c.Request.URI().String(), "/publish/action/") {
		token = c.PostForm("token")
	}

	claims, err := utils.ParseToken(token)
	if err != nil || claims == nil {
		if strings.Contains(c.Request.URI().String(), "/feed/") ||
			strings.Contains(c.Request.URI().String(), "/comment/list/") {
			c.Set("id", 0)
			c.Next(ctx)
			return
		}
		c.Abort()
		c.JSON(http.StatusOK, common.BaseResponse{
			StatusCode: common.LoginStatusExpired,
			StatusMsg:  "登陆状态已失效",
		})
		hlog.CtxErrorf(ctx, "验证出现错误，原因：%v", err)
		return
	}
	username := claims.Audience
	res, err := client.UserClient.GetUserInfoByUsername(ctx, &user.GetUserInfoByUsernameReq{Username: username})
	if err != nil {
		c.Abort()
		c.JSON(http.StatusOK, common.BaseResponse{
			StatusCode: common.ReqError,
			StatusMsg:  "您的登陆状态已失效",
		})
		return
	}
	c.Set("id", res.GetId())
	c.Next(ctx)
}
