package mail

import (
	"gin-vue/middleware/notify/notify"
	"gin-vue/model/settingModel"
	"gin-vue/utils"
	"github.com/sirupsen/logrus"
	"gopkg.in/gomail.v2"
	"strconv"
	"strings"
	"time"
)

type Mail struct {
}

func (mail *Mail) Send(msg notify.Message) {
	model := new(settingModel.Setting)
	mailSetting, err := model.Mail()
	logrus.Debugf("%+v", mailSetting)
	if err != nil {
		logrus.Error("#mail#从数据库获取mail配置失败", err)
		return
	}
	if mailSetting.Host == "" {
		logrus.Error("#mail#Host为空")
		return
	}
	if mailSetting.Port == 0 {
		logrus.Error("#mail#Port为空")
		return
	}
	if mailSetting.User == "" {
		logrus.Error("#mail#User为空")
		return
	}
	if mailSetting.Password == "" {
		logrus.Error("#mail#Password为空")
		return
	}
	msg["content"] = notify.ParseNotifyTemplate(mailSetting.Template, msg)
	toUsers := mail.getActiveMailUsers(mailSetting, msg)
	mail.send(mailSetting, toUsers, msg)
}

func (mail *Mail) send(mailSetting settingModel.Mail, toUsers []string, msg notify.Message) {
	body := msg["content"].(string)
	body = strings.Replace(body, "\n", "<br>", -1)
	gomailMessage := gomail.NewMessage()
	gomailMessage.SetHeader("From", mailSetting.User)
	gomailMessage.SetHeader("To", toUsers...)
	gomailMessage.SetHeader("Subject", "gocron-定时任务通知")
	gomailMessage.SetBody("text/html", body)
	mailer := gomail.NewDialer(mailSetting.Host, mailSetting.Port,
		mailSetting.User, mailSetting.Password)
	maxTimes := 3
	i := 0
	for i < maxTimes {
		err := mailer.DialAndSend(gomailMessage)
		if err == nil {
			break
		}
		i += 1
		time.Sleep(2 * time.Second)
		if i < maxTimes {
			logrus.Errorf("mail#发送消息失败#%s#消息内容-%s", err.Error(), msg["content"])
		}
	}
}

func (mail *Mail) getActiveMailUsers(mailSetting settingModel.Mail, msg notify.Message) []string {
	taskReceiverIds := strings.Split(msg["task_receiver_id"].(string), ",")
	users := []string{}
	for _, v := range mailSetting.MailUsers {
		if utils.InStringSlice(taskReceiverIds, strconv.Itoa(v.Id)) {
			users = append(users, v.Email)
		}
	}
	return users
}
