package mw

import (
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"time"
)

func AccessLog(c context.Context, ctx *app.RequestContext) {
	start := time.Now()
	ctx.Next(c)
	end := time.Now()
	latency := end.Sub(start).Microseconds
	hlog.CtxTracef(c, "status=%d cost=%d method=%s full_path=%s client_ip=%s host=%s query=%s",
		ctx.Response.StatusCode(), latency,
		ctx.Request.Header.Method(), ctx.Request.URI().PathOriginal(), ctx.ClientIP(), ctx.Request.Host(),
		ctx.Request.QueryString())
}
