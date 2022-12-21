package service

import (
	"testing"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/os/glog"
)

// init description
//
// createTime: 2022-12-21 21:00:55
//
// author: hailaz
func init() {
	ctx := gctx.New()
	emailcode, err := g.Cfg().Get(ctx, "emailcode")
	if err != nil {
		glog.Fatal(ctx, err)
	}
	SetEmailCode(emailcode.String())
}

// Test_Mail description
//
// createTime: 2022-12-19 18:02:48
//
// author: hailaz
func Test_Mail(t *testing.T) {
	// gomail
	SendEMail(GetQrcodeMsg("http://www.hailaz.cn"), "test", []string{"hailaz@qq.com"})
}
