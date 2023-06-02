package service

import (
	"context"
	"fmt"

	"github.com/gogf/gf/v2/os/glog"
	goconfluence "github.com/hailaz/confluence-go-api"
)

var basePath = "https://goframe.org"
var apiPath = basePath + "/rest/api"
var cAPI *goconfluence.API
var qaPath = basePath + "/pages/viewpage.action?pageId=7296348"

// NewSearchApi description
//
// createTime: 2022-12-17 16:15:33
//
// author: hailaz
func NewSearchApi(ctx context.Context, token string) {
	// goconfluence.SetDebug(true)
	// initialize a new api instance
	api, err := goconfluence.NewAPI(apiPath, "", token)
	if err != nil {
		glog.Fatal(ctx, err)
	}
	cAPI = api
	// get current user information
	currentUser, err := api.CurrentUser()
	if err != nil {
		glog.Fatal(ctx, err)
	}
	glog.Debugf(ctx, "%+v\n", currentUser)
}

// Search description
//
// createTime: 2022-12-09 16:55:41
//
// author: hailaz
func Search(ctx context.Context, key string) string {
	glog.Debug(ctx, "search key: ", key)
	resStr := ""
	cql := fmt.Sprintf("siteSearch ~ '%s' AND space in ('%s')", key, "gf")
	res, err := cAPI.Search(goconfluence.SearchQuery{
		CQL:   cql,
		Limit: 3,
	})
	if err != nil {
		glog.Error(ctx, err)
		return "哎呀，出错啦~"
	}
	// g.Dump(res)
	if len(res.Results) > 0 {
		resStr = "搜索结果：\n"
		for _, v := range res.Results {
			resStr += v.Content.Title + "\n"
			resStr += basePath + v.Content.Links.WebUI + "\n"
			glog.Debug(ctx, v.Content.Title)
			glog.Debug(ctx, basePath+v.Content.Links.WebUI)
		}
		resStr += "其它常见问题：\n" + qaPath
	} else {
		resStr = "没有找到呢~换个关键字试试^_^\n也可以访问最新文档\n" + basePath
	}

	return resStr
}
