package service

import (
	"context"
	"runtime"
	"strings"
	"time"

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
	reloadStorage := openwechat.NewFileHotReloadStorage("storage.json")

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

var mymsg = "https://item.m.jd.com/product/10026691993401.html?utm_user=plusmember&gx=RnAowmJYaTbZypgWrIMYHXCWUFQ&gxd=RnAokjVfPGHanZ8d_YByWrqF9uKt6mw&ad_od=share&utm_source=androidapp&utm_medium=appshare&utm_campaign=t_335139774&utm_term=CopyURL"

// FuncName description
//
// createTime: 2022-12-19 00:06:47
//
// author: hailaz
func (h *MsgHandler) SyncCheckCallback(resp openwechat.SyncCheckResponse) {
	ctx := gctx.New()
	glog.Debugf(ctx, "RetCode:%s  Selector:%s", resp.RetCode, resp.Selector)
	if resp.Success() {
		if resp.Selector == openwechat.SelectorNormal {
			self, err := h.bot.GetCurrentUser()
			if err != nil {
				glog.Errorf(ctx, "get current user error : %v", err)
				return
			}
			glog.Debugf(ctx, "self : %+v", *self.User)

			mp, err := self.Mps(false)
			if err != nil {
				glog.Errorf(ctx, "get friends error : %v", err)
				return
			}
			for _, v := range mp {
				glog.Debug(ctx, v.ID(), v.NickName, v.UserName)
				// v.SendText("你好")
			}

			// mp.GetByNickName("微信支付").SendText("你好")

			fs, err := self.Friends(true)
			if err != nil {
				glog.Errorf(ctx, "get friends error : %v", err)
				return
			}
			for _, v := range fs {
				glog.Debug(ctx, v.ID(), v.NickName, v.RemarkName, v.UserName)
			}
			// glog.Debugf(ctx, "friends : %+v", fs)
			// fs.GetByNickName("哆啦A梦").SendText("你好")
			glog.Debugf(ctx, "-times : %d", times)
			now := time.Now()
			if times == 0 || now.Hour() == 0 {
				times = now.Hour()
			}
			glog.Debugf(ctx, "--times : %d", times)
			if now.Hour() == times {
				glog.Debugf(ctx, "now : %+v", now)
				// msg := now.Format("2006-01-02 15:04:05") + " " + grand.Letters(now.Second())
				fs.GetByRemarkName("ping").SendText(mymsg)
				times++
			}
			glog.Debugf(ctx, "---times : %d", times)
			// if now.Minute()%20 == 3 && now.Second() < 30 {
			// 	glog.Debugf(ctx, "now : %+v", now)
			// 	// msg := now.Format("2006-01-02 15:04:05") + " " + grand.Letters(now.Second())
			// 	fs.GetByNickName("AA39萌小宝~网购查券助手").SendText(mymsg)
			// }

		}

	} else {
		glog.Debugf(ctx, "sync check error: %s", resp.Err())
	}

}

var times = 0

// Handler 全局处理入口
func (h *MsgHandler) Handler(msg *openwechat.Message) {
	ctx := gctx.New()
	// glog.Debugf(ctx, "hadler Received msg : %+v", *msg)
	glog.Debugf(ctx, "hadler Received msg :%s  %v", msg.MsgType, msg.Content)

	// 处理群消息
	if msg.IsSendByGroup() {
		h.GroupMsg(ctx, msg)
		return
	}

	// 好友申请
	if msg.IsFriendAdd() {
		return
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
		return ReadMsg(ctx, msg)
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
	glog.Debugf(ctx, "Received User%s[%s] %v \nText Msg : %v", sender, sender.ID(), sender.NickName, msg.Content)

	if sender.NickName == "微信团队" {
		glog.Debugf(ctx, "Received Uin %v", sender.Uin)
		return nil
	}

	switch msg.MsgType {
	case openwechat.MsgTypeText:
		return ReadMsg(ctx, msg)
	}

	return err
}

// ReadMsg description
//
// createTime: 2023-06-02 22:11:33
//
// author: hailaz
func ReadMsg(ctx context.Context, msg *openwechat.Message) error {
	requestText := strings.TrimSpace(msg.Content)
	requestText = strings.Trim(requestText, "\n")
	if requestText != "" && strings.HasPrefix(requestText, "gf ") {
		requestText = strings.TrimPrefix(requestText, "gf ")
		reply := Search(ctx, requestText)
		_, err := msg.ReplyText(reply)
		if err != nil {
			glog.Debugf(ctx, "response user error: %v \n", err)
			return err
		}
	}
	return nil
}

// QrCodeCallBack 登录扫码回调，
func (h *MsgHandler) QrCodeCallBack(uuid string) {
	// SendEMail(GetQrcodeMsg("https://login.weixin.qq.com/l/"+uuid), "微信登录二维码", []string{"hailaz@qq.com"})
	if runtime.GOOS == "windows" {
		// 运行在Windows系统上
		openwechat.PrintlnQrcodeUrl(uuid)
	} else {
		glog.Debugf(context.Background(), "login in linux")
		q, _ := qrcode.New("https://login.weixin.qq.com/l/"+uuid, qrcode.Low)
		glog.Debugf(context.Background(), q.ToString(true))
	}
}
