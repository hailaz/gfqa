package service

import (
	"testing"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/os/glog"
)

// TestSearch description
//
// createTime: 2023-10-07 21:39:39
//
// author: hailaz
func TestSearch(t *testing.T) {
	ctx := gctx.New()
	// gf doc init
	token, err := g.Cfg().Get(ctx, "doctoken")
	if err != nil {
		glog.Fatal(ctx, err)
	}
	glog.Debug(ctx, token.String())

	NewSearchApi(ctx, token.String())
	Search(ctx, "开始")
}
