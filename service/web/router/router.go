package router

import (
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/route"
	"toktik/service/web/handler"
	"toktik/service/web/mw"
)

func RegisterRouter(h *server.Hertz) {
	mainRouter := h.Group("/douyin", mw.AccessLog)
	// 视频流接口
	var v handler.VideoHandler
	mainRouter.GET("/feed/", mw.JWT, v.GetFeedList)
	// 挂载用户路由
	userRouter := mainRouter.Group("/user")
	registerUserRouter(userRouter)
	// 挂载发布路由
	publishRouter := mainRouter.Group("/publish", mw.JWT)
	registerPublishRouter(publishRouter)
	// 挂载点赞路由
	favoriteRouter := mainRouter.Group("/favorite", mw.JWT)
	registerFavoriteRouter(favoriteRouter)
	// 挂载评论路由
	commentRouter := mainRouter.Group("/comment", mw.JWT)
	registerCommentRouter(commentRouter)
	// 挂载社交路由
	relationRouter := mainRouter.Group("/relation", mw.JWT)
	registerRelationRouter(relationRouter)
	// 挂载消息路由
	messageRouter := mainRouter.Group("/message", mw.JWT)
	registerMessageRouter(messageRouter)
}

func registerUserRouter(userRouter *route.RouterGroup) {
	var u handler.UserHandler
	userRouter.POST("/login/", u.Login)
	userRouter.POST("/register/", u.Register)
	userRouter.GET("/", mw.JWT, u.GetUserInfo)
}

func registerPublishRouter(publishRouter *route.RouterGroup) {
	var v handler.VideoHandler
	publishRouter.POST("/action/", v.Publish)
	publishRouter.GET("/list/", v.GetPublishList)
}

func registerFavoriteRouter(favoriteRouter *route.RouterGroup) {
	var v handler.VideoHandler
	favoriteRouter.POST("/action/", v.SetFavoriteStatus)
	favoriteRouter.GET("/list/", v.GetFavoriteList)
}

func registerCommentRouter(commentRouter *route.RouterGroup) {
	var v handler.VideoHandler
	commentRouter.POST("/action/", v.SendComment)
	commentRouter.GET("/list/", v.GetCommentList)
}

func registerRelationRouter(relationRouter *route.RouterGroup) {
	var f handler.FollowHandler
	relationRouter.POST("/action/", f.SetFollowStatus)
	relationRouter.GET("/follow/list/", f.GetFollowList)
	relationRouter.GET("/follower/list/", f.GetFollowerList)
	relationRouter.GET("/friend/list/", f.GetFriendList)
}

func registerMessageRouter(messageRouter *route.RouterGroup) {
	var m handler.MessageHandler
	messageRouter.GET("/chat/", m.GetMessageList)
	messageRouter.POST("/action/", m.SendMessage)
}
