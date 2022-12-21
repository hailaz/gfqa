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

var emailCode = ""

// SetEmailCode description
//
// createTime: 2022-12-21 20:55:56
//
// author: hailaz
func SetEmailCode(code string) {
	emailCode = code
}

// SendEMail description
//
// createTime: 2022-12-19 18:03:22
//
// author: hailaz
func SendEMail(m *gomail.Message, subject string, to []string) {
	if m == nil {
		return
	}
	sender := "2464629800@qq.com"

	// send email
	m.SetHeader("From", sender)
	m.SetHeader("To", to...)
	m.SetHeader("Subject", subject)

	d := gomail.NewDialer("smtp.qq.com", 465, sender, emailCode)
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
