package settingModel

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
	"github.com/wuchunfu/JobFlow/middleware/database"
)

type Setting struct {
	Id    int    `gorm:"type:int(20); primary_key; auto_increment; not null"`
	Code  string `gorm:"type:varchar(100); not null"`
	Key   string `gorm:"type:varchar(100); not null"`
	Value string `gorm:"type:varchar(4096); not null; default '' "`
}

const slackTemplate = `
任务ID:  {{.TaskId}}
任务名称: {{.TaskName}}
状态:    {{.Status}}
执行结果: {{.Result}}
`
const emailTemplate = `
任务ID:  {{.TaskId}}
任务名称: {{.TaskName}}
状态:    {{.Status}}
执行结果: {{.Result}}
`
const webhookTemplate = `
{
  "task_id": "{{.TaskId}}",
  "task_name": "{{.TaskName}}",
  "status": "{{.Status}}",
  "result": "{{.Result}}"
}
`

const (
	SlackCode        = "slack"
	SlackUrlKey      = "url"
	SlackTemplateKey = "template"
	SlackChannelKey  = "channel"
)

const (
	MailCode        = "mail"
	MailTemplateKey = "template"
	MailServerKey   = "server"
	MailUserKey     = "user"
)

const (
	WebhookCode        = "webhook"
	WebhookTemplateKey = "template"
	WebhookUrlKey      = "url"
)

// 初始化基本字段 邮件、slack等
func (setting *Setting) InitBasicField() {
	setting.Code = SlackCode
	setting.Key = SlackUrlKey
	setting.Value = ""
	db := database.GetDB()
	db.Model(&setting).Create(&setting)
	setting.Id = 0

	setting.Code = SlackCode
	setting.Key = SlackTemplateKey
	setting.Value = slackTemplate
	db.Model(&setting).Create(&setting)
	setting.Id = 0

	setting.Code = MailCode
	setting.Key = MailServerKey
	setting.Value = ""
	db.Model(&setting).Create(&setting)
	setting.Id = 0

	setting.Code = MailCode
	setting.Key = MailTemplateKey
	setting.Value = emailTemplate
	db.Model(&setting).Create(&setting)
	setting.Id = 0

	setting.Code = WebhookCode
	setting.Key = WebhookTemplateKey
	setting.Value = webhookTemplate
	db.Model(&setting).Create(&setting)
	setting.Id = 0

	setting.Code = WebhookCode
	setting.Key = WebhookUrlKey
	setting.Value = ""
	db.Model(&setting).Create(&setting)
}

// region slack配置

type Slack struct {
	Url      string    `json:"url"`
	Channels []Channel `json:"channels"`
	Template string    `json:"template"`
}

type Channel struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

func (setting *Setting) Slack() (Slack, error) {
	db := database.GetDB()
	list := make([]Setting, 0)
	err := db.Model(&setting).Where("code = ?", SlackCode).Find(&list)
	slack := Slack{}
	if err.Error != nil {
		return slack, err.Error
	}

	setting.formatSlack(list, &slack)

	return slack, err.Error
}

func (setting *Setting) formatSlack(list []Setting, slack *Slack) {
	for _, v := range list {
		switch v.Key {
		case SlackUrlKey:
			slack.Url = v.Value
		case SlackTemplateKey:
			slack.Template = v.Value
		default:
			slack.Channels = append(slack.Channels, Channel{
				v.Id, v.Value,
			})
		}
	}
}

func (setting *Setting) UpdateSlack(url, template string) error {
	db := database.GetDB()
	setting.Value = url

	db.Model(&setting).Select("value").Update(setting, Setting{Code: SlackCode, Key: SlackUrlKey})

	setting.Value = template
	db.Model(&setting).Select("value").Update(setting, Setting{Code: SlackCode, Key: SlackTemplateKey})

	return nil
}

// 创建slack渠道
func (setting *Setting) CreateChannel(channel string) (*Setting, error) {
	setting.Code = SlackCode
	setting.Key = SlackChannelKey
	setting.Value = channel
	db := database.GetDB()
	err := db.Model(&setting).Create(&setting)
	if err.Error != nil {
		return nil, err.Error
	}
	return setting, err.Error
}

func (setting *Setting) IsChannelExist(channel string) bool {
	setting.Code = SlackCode
	setting.Key = SlackChannelKey
	setting.Value = channel

	db := database.GetDB()
	var count int
	err := db.Model(&setting).Find(&setting).Count(&count)
	if err.Error != nil && count <= 0 {
		return false
	} else {
		return count > 0
	}
}

// 删除slack渠道
func (setting *Setting) RemoveChannel(id int) {
	setting.Code = SlackCode
	setting.Key = SlackChannelKey
	setting.Id = id
	db := database.GetDB()
	err := db.Model(&setting).Delete(&setting)
	if err.Error != nil {
		logrus.Error(err.Error)
	}
}

// endregion

type Mail struct {
	Host      string     `json:"host"`
	Port      int        `json:"port"`
	User      string     `json:"user"`
	Password  string     `json:"password"`
	MailUsers []MailUser `json:"mailUsers"`
	Template  string     `json:"template"`
}

type MailUser struct {
	Id       int    `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

// region 邮件配置
func (setting *Setting) Mail() (Mail, error) {
	db := database.GetDB()
	list := make([]Setting, 0)
	err := db.Model(&setting).Where("code = ?", MailCode).Find(&list)
	mail := Mail{MailUsers: make([]MailUser, 0)}
	if err.Error != nil {
		return mail, err.Error
	}

	setting.formatMail(list, &mail)
	return mail, err.Error
}

func (setting *Setting) formatMail(list []Setting, mail *Mail) {
	mailUser := MailUser{}
	for _, v := range list {
		switch v.Key {
		case MailServerKey:
			json.Unmarshal([]byte(v.Value), mail)
		case MailUserKey:
			json.Unmarshal([]byte(v.Value), &mailUser)
			mailUser.Id = v.Id
			mail.MailUsers = append(mail.MailUsers, mailUser)
		case MailTemplateKey:
			mail.Template = v.Value
		}
	}
}

func (setting *Setting) UpdateMail(config, template string) error {
	setting.Value = config
	db := database.GetDB()
	db.Select("value").Update(&setting, Setting{Code: MailCode, Key: MailServerKey})

	setting.Value = template
	db.Select("value").Update(&setting, Setting{Code: MailCode, Key: MailTemplateKey})
	return nil
}

func (setting *Setting) CreateMailUser(username, email string) {
	setting.Code = MailCode
	setting.Key = MailUserKey
	mailUser := MailUser{0, username, email}
	jsonByte, err := json.Marshal(mailUser)
	if err != nil {
		logrus.Error(err.Error())
	}
	setting.Value = string(jsonByte)

	db := database.GetDB()
	createErr := db.Model(&setting).Create(&setting)
	if createErr.Error != nil {
		logrus.Error(createErr.Error)
	}
}

func (setting *Setting) RemoveMailUser(id int) {
	setting.Code = MailCode
	setting.Key = MailUserKey
	setting.Id = id
	db := database.GetDB()
	err := db.Model(&setting).Delete(&setting)
	if err.Error != nil {
		logrus.Error(err.Error)
	}
}

type WebHook struct {
	Url      string `json:"url"`
	Template string `json:"template"`
}

func (setting *Setting) Webhook() (WebHook, error) {
	list := make([]Setting, 0)
	db := database.GetDB()
	err := db.Model(&setting).Where("code = ?", WebhookCode).Find(&list)
	webHook := WebHook{}
	if err.Error != nil {
		return webHook, err.Error
	}

	setting.formatWebhook(list, &webHook)

	return webHook, err.Error
}

func (setting *Setting) formatWebhook(list []Setting, webHook *WebHook) {
	for _, v := range list {
		switch v.Key {
		case WebhookUrlKey:
			webHook.Url = v.Value
		case WebhookTemplateKey:
			webHook.Template = v.Value
		}
	}
}

func (setting *Setting) UpdateWebHook(url, template string) error {
	setting.Value = url

	db := database.GetDB()
	db.Model(&setting).Select("value").Update(setting, Setting{Code: WebhookCode, Key: WebhookUrlKey})

	setting.Value = template
	db.Model(&setting).Select("value").Update(setting, Setting{Code: WebhookCode, Key: WebhookTemplateKey})
	return nil
}

// endregion
