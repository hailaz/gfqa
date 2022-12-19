package service

import (
	"context"
	"runtime"
	"strings"

	"github.com/eatmoreapple/openwechat"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/os/glog"
	"github.com/skip2/go-qrcode"
)

// MsgHandler description
type MsgHandler struct {
	bot *openwechat.Bot
}

// NewHandler description
//
// createTime: 2022-12-19 00:09:54
//
// author: hailaz
func NewHandler(bot *openwechat.Bot) *MsgHandler {
	return &MsgHandler{
		bot: bot,
	}
}

// RunWechat description
//
// createTime: 2022-12-17 15:35:03
//
// author: hailaz
func RunWechat(ctx context.Context) {
	//bot := openwechat.DefaultBot()
	bot := openwechat.DefaultBot(openwechat.Desktop) // 桌面模式，上面登录不上的可以尝试切换这种模式

	handler := NewHandler(bot)

	// 注册消息处理函数
	bot.MessageHandler = handler.Handler

	// 注册登陆二维码回调
	bot.UUIDCallback = handler.QrCodeCallBack

	bot.SyncCheckCallback = handler.SyncCheckCallback

	// 创建热存储容器对象
	reloadStorage := openwechat.NewJsonFileHotReloadStorage("storage.json")

	// 执行热登录
	err := bot.HotLogin(reloadStorage)
	if err != nil {
		if err = bot.Login(); err != nil {
			glog.Errorf(ctx, "login error: %v \n", err)
			return
		}
	}
	// 阻塞主goroutine, 直到发生异常或者用户主动退出
	bot.Block()
}

// FuncName description
//
// createTime: 2022-12-19 00:06:47
//
// author: hailaz
func (h *MsgHandler) SyncCheckCallback(resp openwechat.SyncCheckResponse) {
	ctx := gctx.New()
	glog.Debugf(ctx, "RetCode:%s  Selector:%s", resp.RetCode, resp.Selector)
	if resp.Success() {
		// self, err := h.bot.GetCurrentUser()
		// if err != nil {
		// 	glog.Errorf(ctx, "get current user error : %v", err)
		// 	return
		// }
		// glog.Debugf(ctx, "self : %+v", *self.User)

		// fs, err := self.Friends(false)
		// if err != nil {
		// 	glog.Errorf(ctx, "get friends error : %v", err)
		// 	return
		// }
		// glog.Debugf(ctx, "friends : %+v", fs)
		// fs.GetByNickName("哆啦A梦").SendText("你好")

	} else {
		glog.Debugf(ctx, "sync check error: %s", resp.Err())
	}

}

// Handler 全局处理入口
func (h *MsgHandler) Handler(msg *openwechat.Message) {
	ctx := gctx.New()
	glog.Debugf(ctx, "hadler Received msg : %v", msg.Content)
	// 处理群消息
	if msg.IsSendByGroup() {
		h.GroupMsg(ctx, msg)
		return
	}

	// 好友申请
	if msg.IsFriendAdd() {
		_, err := msg.Agree("你好我是基于chatGPT引擎开发的微信机器人，你可以向我提问任何问题。")
		if err != nil {
			glog.Errorf(ctx, "add friend agree error : %v", err)
			return
		}
	}

	// 私聊
	h.UserMsg(ctx, msg)
}

// GroupMsg description
//
// createTime: 2022-12-17 16:27:05
//
// author: hailaz
func (h *MsgHandler) GroupMsg(ctx context.Context, msg *openwechat.Message) error {
	// 接收群消息
	sender, err := msg.Sender()
	if err != nil {
		glog.Error(ctx, err)
		return err
	}
	group := openwechat.Group{User: sender}
	glog.Debugf(ctx, "Received Group %v Text Msg : %v", group.NickName, msg.Content)

	// 不是@的不处理
	if !msg.IsAt() {
		return nil
	}

	// 获取@我的用户
	groupSender, err := msg.SenderInGroup()
	if err != nil {
		glog.Debugf(ctx, "get sender in group error :%v \n", err)
		return err
	}
	atText := "@" + groupSender.NickName + " \n"

	requestText := strings.TrimSpace(msg.Content)
	requestText = strings.Trim(requestText, "\n")
	if requestText != "" {
		reply := Search(gctx.New(), requestText)
		replyText := atText + reply
		_, err = msg.ReplyText(replyText)
		if err != nil {
			glog.Debugf(ctx, "response group error: %v \n", err)
			return err
		}
	}
	return err
}

// UserMsg description
//
// createTime: 2022-12-17 16:27:05
//
// author: hailaz
func (h *MsgHandler) UserMsg(ctx context.Context, msg *openwechat.Message) error {
	// 接收私聊消息
	sender, err := msg.Sender()
	if err != nil {
		glog.Error(ctx, err)
		return err
	}
	glog.Debugf(ctx, "Received User %v Text Msg : %v", sender.NickName, msg.Content)

	if sender.NickName == "微信团队" {
		glog.Debugf(ctx, "Received Uin %v", sender.Uin)
		return nil
	}

	requestText := strings.TrimSpace(msg.Content)
	requestText = strings.Trim(requestText, "\n")
	if requestText != "" {
		reply := Search(ctx, requestText)
		_, err = msg.ReplyText(reply)
		if err != nil {
			glog.Debugf(ctx, "response user error: %v \n", err)
			return err
		}
	}
	return err
}

// QrCodeCallBack 登录扫码回调，
func (h *MsgHandler) QrCodeCallBack(uuid string) {
	if runtime.GOOS == "windows" {
		// 运行在Windows系统上
		openwechat.PrintlnQrcodeUrl(uuid)
	} else {
		glog.Debugf(context.Background(), "login in linux")
		q, _ := qrcode.New("https://login.weixin.qq.com/l/"+uuid, qrcode.Low)
		glog.Debugf(context.Background(), q.ToString(true))
	}
}
