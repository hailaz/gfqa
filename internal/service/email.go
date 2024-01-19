package service

import (
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"os"

	"github.com/gogf/gf/v2/os/glog"
	"github.com/skip2/go-qrcode"
	"gopkg.in/gomail.v2"
)

// EmailData description
type EmailData struct {
	From     string   // 发送者
	FromCode string   // 发送者密码
	To       []string // 接收者
	IsOpen   bool     // 是否开启
}

var (
	EmailDataSetting = &EmailData{}
)

// SetFromCode description
//
// createTime: 2022-12-21 20:55:56
//
// author: hailaz
func (e *EmailData) SetFromCode(code string) {
	e.FromCode = code
}

// SetEmailFrom description
//
// createTime: 2022-12-21 20:55:56
//
// author: hailaz
func (e *EmailData) SetEmailFrom(from string) {
	e.From = from
}

// SetEmailFrom description
//
// createTime: 2022-12-21 20:55:56
//
// author: hailaz
func (e *EmailData) SetEmailTo(to string) {
	if e.To == nil {
		e.To = []string{}
	}
	e.To = append(e.To, to)
}

// SendEMail description
//
// createTime: 2022-12-19 18:03:22
//
// author: hailaz
func (e *EmailData) SendEMail(m *gomail.Message, subject string, to []string) {
	if !e.IsOpen {
		return
	}
	if m == nil {
		return
	}

	if len(e.From) == 0 || len(e.FromCode) == 0 {
		glog.Error(context.Background(), "emailAddress or emailCode is nil")
		return
	}

	if len(to) == 0 {
		to = e.To
	}

	// send email
	m.SetHeader("From", e.From)
	m.SetHeader("To", to...)
	m.SetHeader("Subject", subject)

	d := gomail.NewDialer("smtp.qq.com", 465, e.From, e.FromCode)
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	if err := d.DialAndSend(m); err != nil {
		panic(err)
	}

}

// GetMsg description
//
// createTime: 2022-12-19 18:37:09
//
// author: hailaz
func GetQrcodeMsg(url string) *gomail.Message {
	filename := "temp.jpg"
	err := qrcode.WriteFile(url, qrcode.Medium, 256, filename)
	if err != nil {
		glog.Error(context.Background(), err)
		return nil
	}

	var html = url + `<img src="cid:%s" alt="My image" />`
	m := gomail.NewMessage()
	m.Embed(filename, gomail.SetCopyFunc(
		func(w io.Writer) error {
			defer os.Remove(filename)
			h, err := os.Open(filename)
			if err != nil {
				return err
			}
			if _, err := io.Copy(w, h); err != nil {
				h.Close()
				return err
			}
			return h.Close()
		},
	))
	m.SetBody("text/html", fmt.Sprintf(html, filename))

	return m
}
