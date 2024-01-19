package cmd

import (
	"context"
	"log"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gcmd"
	"github.com/gogf/gf/v2/os/glog"

	"gfqa/internal/controller"
	"gfqa/internal/service"
)

var (
	Main = gcmd.Command{
		Name:  "main",
		Usage: "main",
		Brief: "start http server",
		Func: func(ctx context.Context, parser *gcmd.Parser) (err error) {
			boot(ctx)
			s := g.Server()
			s.SetLogger(glog.DefaultLogger())
			s.Group("/", func(group *ghttp.RouterGroup) {
				group.Middleware(ghttp.MiddlewareHandlerResponse)
				group.Bind(
					controller.Hello,
				)
			})
			s.SetAccessLogEnabled(true)
			s.Run()
			return nil
		},
	}
)

// boot description
//
// createTime: 2022-12-17 16:18:50
//
// author: hailaz
func boot(ctx context.Context) {
	// init log
	glog.SetDefaultLogger(g.Log())
	glog.SetFlags(glog.F_TIME_STD | glog.F_FILE_SHORT)
	// hello
	glog.Debug(ctx, "hello")

	log.SetOutput(&MyWrite{})

	glog.Debug(ctx, glog.DefaultLogger().GetWriter())

	log.Println("hello")

	// // email init
	emailcode, err := g.Cfg().Get(ctx, "emailSetting")
	if err == nil {
		err = emailcode.Scan(&service.EmailDataSetting)
		if err != nil {
			glog.Fatal(ctx, err)
		}
	}

	// gf doc init
	token, err := g.Cfg().Get(ctx, "doctoken")
	if err != nil {
		glog.Fatal(ctx, err)
	}
	glog.Debug(ctx, token.String())
	if token.String() == "" {
		glog.Fatal(ctx, "config.yaml doctoken is nil")
	}
	service.NewSearchApi(ctx, token.String())

	// wechat
	go service.RunWechat(ctx)
}

// MyWrite description
type MyWrite struct {
}

// Write description
//
// createTime: 2022-12-18 14:29:25
//
// author: hailaz
func (w *MyWrite) Write(p []byte) (n int, err error) {
	glog.Skip(1).Debug(context.Background(), string(p)[20:])
	return len(p), nil
}
