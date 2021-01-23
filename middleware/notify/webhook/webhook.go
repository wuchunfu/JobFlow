package webhook

import (
	"gin-vue/middleware/httpclient"
	"gin-vue/middleware/notify/notify"
	"gin-vue/model/settingModel"
	"gin-vue/utils"
	"github.com/sirupsen/logrus"
	"html"
	"time"
)

type WebHook struct{}

func (webHook *WebHook) Send(msg notify.Message) {
	model := new(settingModel.Setting)
	webHookSetting, err := model.Webhook()
	if err != nil {
		logrus.Error("#webHook#从数据库获取webHook配置失败", err)
		return
	}
	if webHookSetting.Url == "" {
		logrus.Error("#webHook#webhook-url为空")
		return
	}
	logrus.Debugf("%+v", webHookSetting)
	msg["name"] = utils.EscapeJson(msg["name"].(string))
	msg["output"] = utils.EscapeJson(msg["output"].(string))
	msg["content"] = notify.ParseNotifyTemplate(webHookSetting.Template, msg)
	msg["content"] = html.UnescapeString(msg["content"].(string))
	webHook.send(msg, webHookSetting.Url)
}

func (webHook *WebHook) send(msg notify.Message, url string) {
	content := msg["content"].(string)
	timeout := 30
	maxTimes := 3
	i := 0
	for i < maxTimes {
		resp := httpclient.PostJson(url, content, timeout)
		if resp.StatusCode == 200 {
			break
		}
		i += 1
		time.Sleep(2 * time.Second)
		if i < maxTimes {
			logrus.Errorf("webHook#发送消息失败#%s#消息内容-%s", resp.Body, msg["content"])
		}
	}
}
