package service

import (
	"fmt"
	"strings"
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

// TestDate description
//
// createTime: 2024-03-28 23:01:45
func TestDate(t *testing.T) {
	dateFormats := []string{
		"20240328",
		"2024-03-28",
		"2024年3月29日",
		"2024年03月29日",
		"2024年3月40日",
	}

	for _, dateFormat := range dateFormats {
		GetTime(dateFormat)
	}
}

// TestArgs description
//
// createTime: 2024-03-28 23:13:51
func TestArgs(t *testing.T) {
	if agrs := strings.Split("@哆啦A梦 日历", " "); len(agrs) > 1 {
		fmt.Println(agrs)
		GetTime(agrs[1])
	}
}
