/***********************************************************************
* @ 使用第三方服务转发email
* @ brief
	、应交由固定地区的节点来转发，如center，否则容易被第三方当成异地登录，临时封禁~囧

	、国际服gmail，国内qq，天朝用gmail转发超时严重~真是蛋疼啊

	gmail：
		1、用户设置 - 转发和POP/IMAP - POP下载 - 对从现在起收到的邮件启用POP
		2、还须开启“安全性较低的应用的访问权限”
	qq：
		1、用户设置 - 帐户 - POP3/IMAP/SMTP/Exchange/CardDAV/CalDAV服务 - POP3/SMTP服务
		2、取得验证码，替代gomail中的kPasswd

* @ author zhoumf
* @ date 2018-11-15
***********************************************************************/
package email

import (
	"bytes"
	"conf"
	"gopkg.in/gomail"
	"reflect"
	"text/template"
)

var (
	g_dialer *gomail.Dialer
	_msg     = gomail.NewMessage()
)

func SendMail(subject, target, body, language string) error {
	if g_dialer == nil {
		g_dialer = gomail.NewDialer(
			conf.SvrCsv.EmailHost,
			conf.SvrCsv.EmailPort,
			conf.SvrCsv.EmailUser,
			conf.SvrCsv.EmailPasswd)
	}
	_msg.Reset()
	msg := _msg

	body = PackBody(subject, body, language) //嵌入模板，并本地化

	msg.SetAddressHeader("From", g_dialer.Username, "ChillyRoom")
	msg.SetHeader("To", target)
	//msg.SetHeader("Cc" /*抄送*/, "xxxx@foxmail.com")
	//msg.SetHeader("Bcc" /*暗送*/, "xxxx@gmail.com")
	msg.SetHeader("Subject", subject)

	msg.SetBody("text/html", body)

	//msg.Attach("我是附件")

	return g_dialer.DialAndSend(msg)
}

func PackBody(subject, body, language string) string {
	if csv, ok := G_Email[subject]; ok {
		ref := reflect.ValueOf(csv).Elem()
		if v := ref.FieldByName(language); v.IsValid() {
			if t, e := template.New(subject).Parse(v.String()); e == nil {
				var bf bytes.Buffer
				t.Execute(&bf, &body)
				return bf.String()
			}
		}
	}
	return body
}