package service

import (
	"context"
	"runtime"
	"strings"
	"time"

	"github.com/eatmoreapple/openwechat"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/os/glog"
	"github.com/gogf/gf/v2/util/grand"
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
	reloadStorage := openwechat.NewFileHotReloadStorage("./log/storage.json")

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

var (
	times = 0
	msgs  = []string{"hi", "你好", "余额", "收入", "支出", "账单", "贷款", "理财", "投资", "股票", "基金", "保险", "房产", "车辆", "信用卡", "借记卡", "贷记卡", "借款"}
)

// KeepAlive 保活
//
// createTime: 2024-01-19 20:27:38
func (h *MsgHandler) KeepAlive(ctx context.Context) {
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

	now := time.Now()
	if times == 0 || (times == 24 && now.Hour() == 0) {
		times = now.Hour()
	}
	if now.Hour() == times {
		glog.Debugf(ctx, "今天保活第%d次", times)
		glog.Debugf(ctx, "now : %+v", now)
		if chat := mp.GetByNickName("微信支付"); chat != nil {
			chat.SendText(msgs[grand.N(0, len(msgs)-1)])
		} else {
			glog.Errorf(ctx, "注意，没有找到保活对象")
			fs, err := self.Friends(true)
			if err != nil {
				glog.Errorf(ctx, "get friends error : %v", err)
				return
			}
			// 实在不行就改一个好友备注为ping，然后发送
			if chat := fs.GetByRemarkName("ping"); chat != nil {
				chat.SendText("余额")
			}
		}
		times++
	}
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
		if resp.Selector == openwechat.SelectorNormal {
			// h.KeepAlive(ctx)
		}
	} else {
		glog.Debugf(ctx, "sync check error: %s", resp.Err())
	}
}

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
	glog.Debugf(ctx, "Received Group %v Text Msg : [%v]", group.NickName, msg.Content)

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

	cutAt := strings.Join(strings.Split(msg.Content, " ")[1:], " ")

	glog.Debugf(ctx, "cutAt:[%v]", cutAt)

	if strings.HasPrefix(cutAt, "日历") {
		// 日历侠
		img := NewMyImage("src/null.png", "src/simsun.ttc")
		glog.Debugf(ctx, "Received Text Msg [%v]", cutAt)
		if agrs := strings.Split(cutAt, " "); len(agrs) > 1 {
			img.Rili(GetTime(agrs[1]))
		} else {
			img.Rili(time.Now())
		}
		msg.ReplyImage(img.Reader())
	} else {
		return nil
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
	EmailDataSetting.SendEMail(GetQrcodeMsg("https://login.weixin.qq.com/l/"+uuid), "微信登录二维码", nil)
	if runtime.GOOS == "windows" {
		// 运行在Windows系统上
		openwechat.PrintlnQrcodeUrl(uuid)
	} else {
		glog.Debugf(context.Background(), "login in linux")
		q, _ := qrcode.New("https://login.weixin.qq.com/l/"+uuid, qrcode.Low)
		glog.Debugf(context.Background(), q.ToString(true))
	}
}
